package questionservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"giat-cerika-service/configs"
	datasources "giat-cerika-service/internal/dataSources"
	questionrequest "giat-cerika-service/internal/dto/request/question_request"
	"giat-cerika-service/internal/models"
	answerrepo "giat-cerika-service/internal/repositories/answer_repo"
	questionrepo "giat-cerika-service/internal/repositories/question_repo"
	quizrepo "giat-cerika-service/internal/repositories/quiz_repo"
	errorresponse "giat-cerika-service/pkg/constant/error_response"
	rabbitmq "giat-cerika-service/pkg/constant/rabbitMq"
	"giat-cerika-service/pkg/utils"
	"giat-cerika-service/pkg/workers/payload"
	"io"
	"mime/multipart"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type QuestionServiceImpl struct {
	questionRepo questionrepo.IQuestionRepository
	quizRepo     quizrepo.IQuizRepository
	answerRepo   answerrepo.IAnswerRepository
	rdb          *redis.Client
	cld          datasources.CloudinaryService
}

func NewQuestionServiceImpl(
	questionRepo questionrepo.IQuestionRepository,
	quizRepo quizrepo.IQuizRepository,
	answerRepo answerrepo.IAnswerRepository,
	rdb *redis.Client,
	cld datasources.CloudinaryService,
) IQuestionService {
	return &QuestionServiceImpl{
		questionRepo: questionRepo,
		answerRepo:   answerRepo,
		quizRepo:     quizRepo,
		rdb:          rdb,
		cld:          cld,
	}
}

func (q *QuestionServiceImpl) invalidateCacheQuestion(ctx context.Context) {
	iter := q.rdb.Scan(ctx, 0, "questions:*", 0).Iterator()
	for iter.Next(ctx) {
		q.rdb.Del(ctx, iter.Val()).Err()
	}

	IterID := q.rdb.Scan(ctx, 0, "question:*", 0).Iterator()
	for IterID.Next(ctx) {
		q.rdb.Del(ctx, IterID.Val()).Err()
	}
}

func (q *QuestionServiceImpl) invalidateCacheQuiz(ctx context.Context) {
	iter := q.rdb.Scan(ctx, 0, "quizzes:*", 0).Iterator()
	for iter.Next(ctx) {
		q.rdb.Del(ctx, iter.Val())
	}

	iterID := q.rdb.Scan(ctx, 0, "quiz:*", 0).Iterator()
	for iterID.Next(ctx) {
		q.rdb.Del(ctx, iterID.Val())
	}
}

func fileQuestionToBytes(fh *multipart.FileHeader) ([]byte, error) {
	file, err := fh.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return io.ReadAll(file)
}

var PublishImageQuestion = func(p payload.ImageUploadPayload) {
	go func() {
		_ = rabbitmq.PublishToQueue(
			"",
			rabbitmq.SendImageQuestionQueueName,
			p,
		)
	}()
}

// CreateQuestion implements IQuestionService.
func (q *QuestionServiceImpl) CreateQuestion(ctx context.Context, req questionrequest.CreateQuestionRequest) error {
	quiz, err := q.quizRepo.FindById(ctx, req.QuizId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get quiz", 500)
	}
	if quiz == nil {
		return errorresponse.NewCustomError(errorresponse.ErrNotFound, "quiz not found", 404)
	}
	if req.QuizId == uuid.Nil {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "quiz cannot be empty", 400)
	}
	if strings.TrimSpace(req.QuestionText) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "question text cannot be empty", 400)
	}
	if len(req.Answers) == 0 {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "answers cannot be empty", 400)
	}

	question := &models.Question{
		ID:           uuid.New(),
		QuizID:       quiz.ID,
		QuestionText: req.QuestionText,
	}

	if err := q.questionRepo.CreateQuestion(ctx, question); err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to create question", 500)
	}

	go func(quizId uuid.UUID) {
		ctxBg := context.Background()
		if err := q.quizRepo.IncreamentAmountQuestion(ctxBg, quizId); err != nil {
			fmt.Println("failed to increament amount question", err)
		}
	}(question.QuizID)

	if req.QuestionImage != nil {
		if bin, err := fileQuestionToBytes(req.QuestionImage); err == nil && len(bin) > 0 {
			go PublishImageQuestion(payload.ImageUploadPayload{
				ID:        question.ID,
				Type:      "single",
				FileBytes: bin,
				Folder:    "giat_ceria/questions",
				Filename:  fmt.Sprintf("question_%s_image", question.ID.String()),
			})
		}
	}

	for _, ansReq := range req.Answers {
		answer := &models.Answer{
			ID:         uuid.New(),
			QuestionID: question.ID,
			AnswerText: ansReq.AnswerText,
			ScoreValue: ansReq.ScoreValue,
		}
		if err := q.answerRepo.CreateAnswer(ctx, answer); err != nil {
			return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to create answer", 500)
		}
	}

	q.invalidateCacheQuestion(ctx)
	q.invalidateCacheQuiz(ctx)
	return nil
}

// FindAllQuestions implements IQuestionService.
func (q *QuestionServiceImpl) FindAllQuestions(
	ctx context.Context,
	quizId uuid.UUID,
	page int,
	limit int,
	search string,
) ([]*models.Question, int, error) {

	cacheKey := fmt.Sprintf(
		"questions:quiz=%s:search=%s:page=%d:limit=%d",
		quizId.String(),
		search,
		page,
		limit,
	)

	if cached, err := configs.GetRedis(ctx, cacheKey); err == nil && len(cached) > 0 {
		var result struct {
			Data  []*models.Question `json:"data"`
			Total int                `json:"total"`
		}
		if json.Unmarshal([]byte(cached), &result) == nil {
			return result.Data, result.Total, nil
		}
	}

	offset := (page - 1) * limit

	quiz, err := q.quizRepo.FindById(ctx, quizId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, errorresponse.NewCustomError(
				errorresponse.ErrNotFound,
				"quiz not found",
				404,
			)
		}
		return nil, 0, errorresponse.NewCustomError(
			errorresponse.ErrInternal,
			"failed to get quiz",
			500,
		)
	}

	items, total, err := q.questionRepo.FindAllQuestions(
		ctx,
		quiz.ID,
		limit,
		offset,
		search,
	)
	if err != nil {
		return nil, 0, errorresponse.NewCustomError(
			errorresponse.ErrInternal,
			"failed to get questions",
			500,
		)
	}

	if items == nil {
		items = []*models.Question{}
	}

	buf, _ := json.Marshal(map[string]any{
		"data":  items,
		"total": total,
	})
	_ = configs.SetRedis(ctx, cacheKey, buf, time.Minute*30)

	return items, total, nil
}

// FindQuestionById implements IQuestionService.
func (q *QuestionServiceImpl) FindQuestionById(ctx context.Context, questionId uuid.UUID) (*models.Question, error) {
	cacheKey := fmt.Sprintf("question:id=%s", questionId.String())
	if cached, err := configs.GetRedis(ctx, cacheKey); err == nil && len(cached) > 0 {
		var question models.Question
		if json.Unmarshal([]byte(cached), &question) == nil {
			return &question, nil
		}
	}

	question, err := q.questionRepo.FindQuestionById(ctx, questionId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorresponse.NewCustomError(errorresponse.ErrNotFound, "question not found", 404)
		}
		return nil, errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get question", 500)
	}
	buf, _ := json.Marshal(question)
	_ = configs.SetRedis(ctx, cacheKey, buf, time.Minute*30)
	return question, nil
}

// UpdateQuestion implements IQuestionService.
func (q *QuestionServiceImpl) UpdateQuestion(ctx context.Context, questionId uuid.UUID, req questionrequest.UpdateQuestionRequest) error {
	question, err := q.questionRepo.FindQuestionById(ctx, questionId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorresponse.NewCustomError(errorresponse.ErrNotFound, "question not found", 404)
		}
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get question", 500)
	}

	quiz, err := q.quizRepo.FindById(ctx, req.QuizId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get quiz", 500)
	}

	if req.QuizId != uuid.Nil {
		question.QuizID = quiz.ID
	}
	if strings.TrimSpace(req.QuestionText) != "" {
		question.QuestionText = req.QuestionText
	}

	if req.QuestionImage != nil {
		if question.QuestionImage != "" {
			publicID := utils.ExtractPublicIDFromCloudinaryURL(question.QuestionImage)
			if publicID != "" {
				_ = q.cld.DestroyImage(ctx, publicID)
			}
		}
		if bin, err := fileQuestionToBytes(req.QuestionImage); err == nil && len(bin) > 0 {
			PublishImageQuestion(payload.ImageUploadPayload{
				ID:        question.ID,
				Type:      "single",
				FileBytes: bin,
				Folder:    "giat_ceria/questions",
				Filename:  fmt.Sprintf("question_%s_image", question.ID.String()),
			})
		}
	}

	if len(req.Answers) > 0 {
		// Delete old answers
		if err := q.answerRepo.DeleteByQuestionID(ctx, question.ID); err != nil {
			return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to delete old answers", 500)
		}

		// Insert new answers
		for _, ans := range req.Answers {
			newAnswer := &models.Answer{
				ID:         uuid.New(),
				QuestionID: question.ID,
				AnswerText: ans.AnswerText,
				ScoreValue: ans.ScoreValue,
			}

			if err := q.answerRepo.CreateAnswer(ctx, newAnswer); err != nil {
				return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to create answer", 500)
			}
		}
	}

	if err := q.questionRepo.UpdateQuestion(ctx, questionId, question); err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to update question", 500)
	}

	q.invalidateCacheQuestion(ctx)

	return nil
}

// DeleteQuestion implements IQuestionService.
func (q *QuestionServiceImpl) DeleteQuestion(ctx context.Context, questionId uuid.UUID) error {
	question, err := q.questionRepo.FindQuestionById(ctx, questionId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorresponse.NewCustomError(errorresponse.ErrNotFound, "question not found", 404)
		}
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get question", 500)
	}
	if question.Answers != nil {
		if err := q.answerRepo.DeleteByQuestionID(ctx, question.ID); err != nil {
			return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to delete answers", 500)
		}
	}
	if question.QuestionImage != "" {
		publicID := utils.ExtractPublicIDFromCloudinaryURL(question.QuestionImage)
		if publicID != "" {
			_ = q.cld.DestroyImage(ctx, publicID)
		}
	}

	if err := q.questionRepo.DeleteQuestion(ctx, questionId); err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to delete question", 500)
	}
	go func(quizId uuid.UUID) {
		ctxBg := context.Background()
		if err := q.quizRepo.DecreaseAmountQuestion(ctxBg, quizId); err != nil {
			fmt.Println("failed to decrease amount question", err)
		}
	}(question.QuizID)
	q.invalidateCacheQuestion(ctx)
	q.invalidateCacheQuiz(ctx)
	return nil
}
