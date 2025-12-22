package quizsessionservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	quizrequest "giat-cerika-service/internal/dto/request/quiz_request"
	quizsessionresponse "giat-cerika-service/internal/dto/response/quiz_session_response"
	"giat-cerika-service/internal/models"
	quizrepo "giat-cerika-service/internal/repositories/quiz_repo"
	quizsessionrepo "giat-cerika-service/internal/repositories/quiz_session_repo"
	studentrepo "giat-cerika-service/internal/repositories/student_repo"
	errorresponse "giat-cerika-service/pkg/constant/error_response"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type QuizSessionServiceImpl struct {
	quizSessionRepo quizsessionrepo.IQuizSessionRepository
	quizRepo        quizrepo.IQuizRepository
	studentRepo     studentrepo.IStudentRepository
	rdb             *redis.Client
}

func NewQuizSessionServiceImpl(qsRepo quizsessionrepo.IQuizSessionRepository, quizRepo quizrepo.IQuizRepository, studentRepo studentrepo.IStudentRepository, rdb *redis.Client) IQuizSessionService {
	return &QuizSessionServiceImpl{quizSessionRepo: qsRepo, quizRepo: quizRepo, studentRepo: studentRepo, rdb: rdb}
}

func (q *QuizSessionServiceImpl) invalidateCacheQuiz(ctx context.Context) {
	iter := q.rdb.Scan(ctx, 0, "quizzes:*", 0).Iterator()
	for iter.Next(ctx) {
		q.rdb.Del(ctx, iter.Val())
	}

	iterID := q.rdb.Scan(ctx, 0, "quiz:*", 0).Iterator()
	for iterID.Next(ctx) {
		q.rdb.Del(ctx, iterID.Val())
	}
}

// AssignCodeQuiz implements [IQuizSessionService].
func (q *QuizSessionServiceImpl) AssignCodeQuiz(ctx context.Context, userId uuid.UUID, quizId uuid.UUID, code string) (*models.QuizSession, error) {
	quiz, err := q.quizRepo.FindById(ctx, quizId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorresponse.NewCustomError(errorresponse.ErrNotFound, "quiz not found", 404)
		}
		return nil, errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get quiz", 500)
	}

	student, err := q.studentRepo.FindByStudentID(ctx, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorresponse.NewCustomError(errorresponse.ErrNotFound, "student not found", 404)
		}
		return nil, errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get student", 500)
	}

	if strings.TrimSpace(code) == "" {
		return nil, errorresponse.NewCustomError(errorresponse.ErrBadRequest, "code is required", 400)
	}

	if quiz.Status == 0 || quiz.Status == 2 {
		return nil, errorresponse.NewCustomError(errorresponse.ErrBadRequest, "quiz is not available", 400)
	}

	if quizId == quiz.ID && code != quiz.Code {
		return nil, errorresponse.NewCustomError(errorresponse.ErrExists, "code access incorrect", 409)
	}
	isComplete, err := q.quizSessionRepo.FindCompleteStatusQuizSession(ctx, student.ID, quiz.ID)
	if err != nil {
		return nil, errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to find quiz session", 500)
	}
	if isComplete {
		return nil, errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Quiz Already Attempt", 500)
	}

	existingSession, err := q.quizSessionRepo.FindByUserAndQuiz(ctx, student.ID, quiz.ID)
	if err == nil && existingSession != nil {
		return existingSession, nil
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to check existing session", 500)
	}

	newQuizSession := &models.QuizSession{
		ID:          uuid.New(),
		UserID:      student.ID,
		QuizID:      quiz.ID,
		Score:       0,
		MaxScore:    0,
		Status:      models.SessionStatusStarted,
		StartedAt:   nil,
		CompletedAt: nil,
	}

	if err := q.quizSessionRepo.SaveQuizSession(ctx, newQuizSession); err != nil {
		return nil, errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to save quiz session", 500)
	}

	go func(quizId uuid.UUID) {
		ctxBg := context.Background()
		if err := q.quizRepo.IncreamentAmountAssigned(ctxBg, quizId); err != nil {
			fmt.Println("failed to increament amount assigned quiz", err)
		}
	}(newQuizSession.QuizID)

	q.invalidateCacheQuiz(ctx)

	return newQuizSession, nil
}

// StartQuizSession implements [IQuizSessionService].
func (q *QuizSessionServiceImpl) StartQuizSession(
	ctx context.Context,
	userId uuid.UUID,
	quizSessionId uuid.UUID,
) (*quizsessionresponse.QuizSessionStartResponse, error) {

	// =========================
	// GET STUDENT
	// =========================
	student, err := q.studentRepo.FindByStudentID(ctx, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorresponse.NewCustomError(
				errorresponse.ErrNotFound,
				"student not found",
				404,
			)
		}
		return nil, errorresponse.NewCustomError(
			errorresponse.ErrInternal,
			"failed to get student",
			500,
		)
	}

	// =========================
	// GET QUIZ SESSION
	// =========================
	quizSession, err := q.quizSessionRepo.FindById(ctx, student.ID, quizSessionId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorresponse.NewCustomError(
				errorresponse.ErrNotFound,
				"quiz session not found",
				404,
			)
		}
		return nil, errorresponse.NewCustomError(
			errorresponse.ErrInternal,
			"failed to get quiz session",
			500,
		)
	}

	if quizSession.Status == models.SessionStatusCompleted {
		return nil, errorresponse.NewCustomError(
			errorresponse.ErrBadRequest,
			"quiz session already completed",
			400,
		)
	}

	// =========================
	// GET QUIZ DETAIL
	// =========================
	quiz, err := q.quizRepo.FindById(ctx, quizSession.QuizID)
	if err != nil {
		return nil, errorresponse.NewCustomError(
			errorresponse.ErrInternal,
			"failed to get quiz",
			500,
		)
	}

	// =========================
	// TIMEZONE FIX (NO .In())
	// =========================
	locJakarta, _ := time.LoadLocation("Asia/Jakarta")

	now := time.Now().In(locJakarta)

	// 🔥 REBUILD TIME AS WIB (timestamp without tz fix)
	start := time.Date(
		quiz.StartDate.Year(),
		quiz.StartDate.Month(),
		quiz.StartDate.Day(),
		quiz.StartDate.Hour(),
		quiz.StartDate.Minute(),
		quiz.StartDate.Second(),
		0,
		locJakarta,
	)

	end := time.Date(
		quiz.EndDate.Year(),
		quiz.EndDate.Month(),
		quiz.EndDate.Day(),
		quiz.EndDate.Hour(),
		quiz.EndDate.Minute(),
		quiz.EndDate.Second(),
		0,
		locJakarta,
	)

	// =========================
	// VALIDASI START TIME
	// =========================
	hasSpecificTime := !(start.Hour() == 0 && start.Minute() == 0 && start.Second() == 0)

	if hasSpecificTime {
		if now.Before(start) {
			return nil, errorresponse.NewCustomError(
				errorresponse.ErrBadRequest,
				fmt.Sprintf(
					"quiz can be started at %s",
					start.Format("02 Jan 2006 15:04"),
				),
				400,
			)
		}
	} else {
		nowDate := time.Date(
			now.Year(), now.Month(), now.Day(),
			0, 0, 0, 0, locJakarta,
		)

		startDate := time.Date(
			start.Year(), start.Month(), start.Day(),
			0, 0, 0, 0, locJakarta,
		)

		if nowDate.Before(startDate) {
			return nil, errorresponse.NewCustomError(
				errorresponse.ErrBadRequest,
				fmt.Sprintf(
					"quiz can be started on %s",
					startDate.Format("02 Jan 2006"),
				),
				400,
			)
		}
	}

	// =========================
	// DURASI & END TIME
	// =========================
	// =========================
	// DURASI & END TIME (FIXED)
	// =========================
	var durationSeconds int64
	var isUnlimited bool
	var endTime *time.Time

	redisKey := fmt.Sprintf("quiz_session:%s:duration", quizSession.ID.String())

	if start.IsZero() || end.IsZero() || start.Equal(end) {
		// Unlimited quiz
		isUnlimited = true
		durationSeconds = 0
	} else {
		existingTTL := q.rdb.TTL(ctx, redisKey).Val()

		if existingTTL > 0 {
			// Resume quiz
			durationSeconds = int64(existingTTL.Seconds())

			calculatedEndTime := now.Add(existingTTL)
			endTime = &calculatedEndTime

			isUnlimited = false
		} else {
			// 🔥 FIRST START → HITUNG SISA WAKTU
			remaining := end.Sub(now)

			// ❌ Sudah lewat end time
			if remaining <= 0 {
				return nil, errorresponse.NewCustomError(
					errorresponse.ErrBadRequest,
					"quiz time has ended",
					400,
				)
			}

			durationSeconds = int64(remaining.Seconds())

			calculatedEndTime := now.Add(remaining)
			endTime = &calculatedEndTime

			isUnlimited = false

			// Simpan sisa waktu ke Redis
			_ = q.rdb.Set(
				ctx,
				redisKey,
				durationSeconds,
				time.Duration(durationSeconds)*time.Second,
			).Err()
		}
	}

	// =========================
	// UPDATE STARTED AT
	// =========================
	if quizSession.Status == models.SessionStatusStarted {
		err = q.quizSessionRepo.CreateStartedAt(ctx, quizSession.ID)
		if err != nil {
			return nil, errorresponse.NewCustomError(
				errorresponse.ErrInternal,
				"failed to start quiz",
				500,
			)
		}
	}

	// =========================
	// RESPONSE
	// =========================
	response := &quizsessionresponse.QuizSessionStartResponse{
		QuizSessionID:   quizSession.ID,
		QuizID:          quiz.ID,
		DurationSeconds: durationSeconds,
		IsUnlimited:     isUnlimited,
		EndTime:         endTime,
	}

	return response, nil
}

// GetQuizSessionDuration implements [IQuizSessionService].
func (q *QuizSessionServiceImpl) GetQuizSessionDuration(ctx context.Context, userId uuid.UUID, quizSessionId uuid.UUID) (*quizsessionresponse.QuizSessionDurationResponse, error) {
	redisKey := fmt.Sprintf("quiz_session:%s:duration", quizSessionId.String())

	// Get TTL dari Redis
	ttl := q.rdb.TTL(ctx, redisKey).Val()

	student, err := q.studentRepo.FindByStudentID(ctx, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorresponse.NewCustomError(errorresponse.ErrNotFound, "student not found", 404)
		}
		return nil, errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get student", 500)
	}

	// Cek quiz session
	quizSession, err := q.quizSessionRepo.FindById(ctx, student.ID, quizSessionId)
	if err != nil {
		return nil, errorresponse.NewCustomError(errorresponse.ErrNotFound, "quiz session not found", 404)
	}

	// Jika status completed
	if quizSession.Status == models.SessionStatusCompleted {
		return &quizsessionresponse.QuizSessionDurationResponse{
			RemainingSeconds: 0,
			IsExpired:        true,
			IsUnlimited:      false,
		}, nil
	}

	// Get quiz detail untuk cek apakah unlimited
	quiz, err := q.quizRepo.FindById(ctx, quizSession.QuizID)
	if err != nil {
		return nil, errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get quiz", 500)
	}

	// Cek apakah quiz memang unlimited dari awal
	isQuizUnlimited := quiz.StartDate.IsZero() || quiz.EndDate.IsZero() || quiz.StartDate.Equal(quiz.EndDate)

	// Jika quiz memang unlimited
	if isQuizUnlimited {
		return &quizsessionresponse.QuizSessionDurationResponse{
			RemainingSeconds: 0,
			IsExpired:        false,
			IsUnlimited:      true,
		}, nil
	}

	// Jika TTL masih ada (quiz belum expired)
	if ttl > 0 {
		return &quizsessionresponse.QuizSessionDurationResponse{
			RemainingSeconds: int64(ttl.Seconds()),
			IsExpired:        false,
			IsUnlimited:      false,
		}, nil
	}

	// Jika TTL habis atau tidak ada, tapi quiz ada batas waktu
	// Cek apakah quiz sudah pernah di-start (ada StartedAt)
	if quizSession.StartedAt != nil {
		// Quiz sudah di-start tapi TTL habis = EXPIRED
		return &quizsessionresponse.QuizSessionDurationResponse{
			RemainingSeconds: 0,
			IsExpired:        true,
			IsUnlimited:      false,
		}, nil
	}

	// Jika belum pernah di-start, hitung durasi dari quiz
	duration := quiz.EndDate.Sub(quiz.StartDate)
	return &quizsessionresponse.QuizSessionDurationResponse{
		RemainingSeconds: int64(duration.Seconds()),
		IsExpired:        false,
		IsUnlimited:      false,
	}, nil
}

// SubmtiQuizSession implements [IQuizSessionService].
func (q *QuizSessionServiceImpl) SubmtiQuizSession(ctx context.Context, userId uuid.UUID, quizSessionId uuid.UUID, req quizrequest.SubmitQuizRequest) error {
	student, err := q.studentRepo.FindByStudentID(ctx, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorresponse.NewCustomError(errorresponse.ErrNotFound, "student not found", 404)
		}
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get student", 500)
	}
	quizSession, err := q.quizSessionRepo.FindById(ctx, student.ID, quizSessionId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorresponse.NewCustomError(errorresponse.ErrNotFound, "quiz session not found", 404)
		}
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get quiz session", 500)
	}
	if quizSession.Status == models.SessionStatusCompleted {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "quiz already submitted", 400)
	}

	quiz, err := q.quizRepo.FindById(ctx, quizSession.QuizID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorresponse.NewCustomError(errorresponse.ErrNotFound, "quiz not found", 404)
		}
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get quiz", 500)
	}

	totalQuestions := len(quiz.Questions)
	if totalQuestions == 0 {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "quiz has no question", 500)
	}

	var responseToSave []*models.Response
	var totalScore int
	var maxScore int

	submittedAnswers := make(map[uuid.UUID]uuid.UUID)
	for _, sub := range req.Answers {
		submittedAnswers[sub.QuestionID] = sub.AnswerID
	}

	answerScoreMap := make(map[uuid.UUID]int)
	for _, question := range quiz.Questions {
		qID := question.ID
		currentQuestionMaxScore := 0

		for _, answer := range question.Answers {
			if answer.ScoreValue > currentQuestionMaxScore {
				currentQuestionMaxScore = answer.ScoreValue
			}
			answerScoreMap[answer.ID] = answer.ScoreValue
		}

		maxScore += currentQuestionMaxScore

		answerID, submitted := submittedAnswers[qID]
		scoreEarned := 0

		isAnswerSubmitted := submitted && answerID != uuid.Nil

		if isAnswerSubmitted {
			scoreEarned = answerScoreMap[answerID]
			totalScore += scoreEarned

			response := &models.Response{
				ID:            uuid.New(),
				QuizSessionID: quizSessionId,
				QuestionID:    qID,
				AnswerID:      &answerID, // AnswerID disimpan (bukan nil)
				ScoreEarned:   scoreEarned,
			}

			responseToSave = append(responseToSave, response)
		} else {
			response := &models.Response{
				ID:            uuid.New(),
				QuizSessionID: quizSessionId,
				QuestionID:    qID,
				AnswerID:      nil, // AnswerID tetap nil
				ScoreEarned:   0,
			}
			responseToSave = append(responseToSave, response)
		}
	}

	redisKey := fmt.Sprintf("quiz_session:%s:duration", quizSessionId.String())
	q.rdb.Del(ctx, redisKey)

	if err := q.quizSessionRepo.BulkSaveResponses(ctx, quizSessionId, responseToSave); err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to save quiz responses", 500)
	}

	percentage := 0.0
	if maxScore > 0 {
		percentage = float64(totalScore) / float64(maxScore) * 100
	}

	quizHistoryStatus := 1
	if percentage < 40.0 {
		quizHistoryStatus = 3
	} else if percentage > 40.0 && percentage <= 60.0 {
		quizHistoryStatus = 2
	} else if percentage > 60.0 {
		quizHistoryStatus = 1
	} else {
		quizHistoryStatus = 0
	}

	completedAt := time.Now()

	err = q.quizSessionRepo.CompleteQuizSession(ctx, quizSessionId, totalScore, maxScore, &completedAt)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to submit quiz", 500)
	}

	startedAtTime := time.Time{}
	if quizSession.StartedAt != nil {
		startedAtTime = *quizSession.StartedAt
	}

	quizHistory := models.QuizHistory{
		ID:              uuid.New(),
		QuizID:          quizSession.QuizID,
		QuizSessionID:   quizSession.ID,
		Code:            quiz.Code,
		Title:           quiz.Title,
		Description:     quiz.Description,
		StartDate:       &quiz.StartDate,
		EndDate:         &quiz.EndDate,
		AmountQuestions: quiz.AmountQuestions,
		AmountAssigned:  quiz.AmountAssigned,
		UserID:          quizSession.UserID,
		Score:           totalScore,
		MaxScore:        maxScore,
		Percentage:      percentage,
		StartedAt:       &startedAtTime,
		CompletedAt:     &completedAt,
		Status:          models.SessionStatusCompleted,
		StatusCategory:  quizHistoryStatus,
	}

	var questionHistories []models.QuestionHistory
	var answerHistories []models.AnswerHistory

	responseAnswerIDMap := make(map[uuid.UUID]uuid.UUID)
	responseScoreMap := make(map[uuid.UUID]int)

	for _, r := range responseToSave {
		if r.AnswerID != nil {
			responseAnswerIDMap[r.QuestionID] = *r.AnswerID
			responseScoreMap[r.QuestionID] = r.ScoreEarned
		} else {
			responseAnswerIDMap[r.QuestionID] = uuid.Nil
			responseScoreMap[r.QuestionID] = 0
		}
	}

	for _, question := range quiz.Questions {
		qHistID := uuid.New()

		qHistory := models.QuestionHistory{
			ID:            qHistID,
			QuizHistoryID: quizHistory.ID,
			QuestionID:    question.ID,
			QuestionText:  question.QuestionText,
			QuestionImage: question.QuestionImage,
		}

		questionHistories = append(questionHistories, qHistory)

		submittedAnswerID := responseAnswerIDMap[question.ID]

		for _, answer := range question.Answers {
			scoreEarned := 0

			if answer.ID == submittedAnswerID && submittedAnswerID != uuid.Nil {
				scoreEarned = responseScoreMap[question.ID]
			}

			aHistory := models.AnswerHistory{
				ID:                uuid.New(),
				QuestionHistoryID: qHistID,
				AnswerID:          answer.ID,
				AnswerText:        answer.AnswerText,
				ScoreValue:        answer.ScoreValue,
				ScoreEarned:       scoreEarned,
			}

			answerHistories = append(answerHistories, aHistory)
		}
	}

	if err := q.quizSessionRepo.SaveQuizHistoryTransaction(ctx, &quizHistory, questionHistories, answerHistories); err != nil {
		fmt.Println("failed to save quiz history:", err)
	}

	q.invalidateCacheQuiz(ctx)
	return nil
}

// GetOrderedQuizQuestions implements [IQuizSessionService].
func (q *QuizSessionServiceImpl) GetOrderedQuizQuestions(ctx context.Context, userId uuid.UUID, quizSessionId uuid.UUID) (*quizsessionresponse.OrderedQuizQuestionsResponse, error) {
	student, err := q.studentRepo.FindByStudentID(ctx, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorresponse.NewCustomError(errorresponse.ErrNotFound, "student not found", 404)
		}
		return nil, errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get student", 500)
	}
	quizSession, err := q.quizSessionRepo.FindById(ctx, student.ID, quizSessionId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorresponse.NewCustomError(errorresponse.ErrNotFound, "quiz session not found", 404)
		}
		return nil, errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get quiz session", 500)
	}

	tempQuiz, err := q.quizRepo.FindById(ctx, quizSession.QuizID)
	if err != nil {
		return nil, errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get quiz data", 500)
	}

	cacheKey := fmt.Sprintf(
		"quiz_session:%s:ordered_questions:%s:user_id:%s",
		quizSessionId.String(),
		tempQuiz.QuestionOrderMode,
		quizSession.UserID,
	)

	if cached, err := q.rdb.Get(ctx, cacheKey).Result(); err == nil && cached != "" {
		var cachedResponse quizsessionresponse.OrderedQuizQuestionsResponse
		if json.Unmarshal([]byte(cached), &cachedResponse) == nil {
			return &cachedResponse, nil
		}
	}

	quiz, err := q.quizSessionRepo.FindQuizWithOrderedQuestions(ctx, quizSession.QuizID, string(tempQuiz.QuestionOrderMode))
	if err != nil {
		return nil, errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get ordered quiz questions", 500)
	}

	var questionResponse []quizsessionresponse.QuestionDetailResponse
	for _, question := range quiz.Questions {
		var answerResponse []quizsessionresponse.AnswerResponse
		for _, answer := range question.Answers {
			answerResponse = append(answerResponse, quizsessionresponse.AnswerResponse{
				ID:         answer.ID,
				AnswerText: answer.AnswerText,
				ScoreValue: answer.ScoreValue,
			})
		}
		questionResponse = append(questionResponse, quizsessionresponse.QuestionDetailResponse{
			ID:            question.ID,
			QuestionText:  question.QuestionText,
			QuestionImage: &question.QuestionImage,
			Answers:       answerResponse,
		})
	}

	response := &quizsessionresponse.OrderedQuizQuestionsResponse{
		QuizID:    quiz.ID,
		Questions: questionResponse,
	}

	if buf, err := json.Marshal(response); err == nil {
		_ = q.rdb.Set(ctx, cacheKey, buf, time.Minute*30).Err()
	}

	return response, nil
}

// GetQuizSessionStudentByQuiz implements [IQuizSessionService].
func (q *QuizSessionServiceImpl) GetQuizSessionStudentByQuiz(ctx context.Context) ([]quizsessionresponse.ListQuestionSessionResponse, error) {
	sessions, err := q.quizSessionRepo.FindQuizSessionByQuiz(ctx)
	if err != nil {
		return nil, err
	}

	grouped := make(map[uuid.UUID]*quizsessionresponse.ListQuestionSessionResponse)

	for _, s := range sessions {
		quizID := s.Quiz.ID

		if existing, ok := grouped[quizID]; !ok {
			// Jika belum ada, gunakan fungsi helper untuk inisialisasi
			resp := quizsessionresponse.ToListQuestionSessionResponse(s)
			grouped[quizID] = &resp
		} else {
			// Jika sudah ada, cukup append detailnya saja menggunakan helper detail
			existing.DetailQuizSession = append(
				existing.DetailQuizSession,
				quizsessionresponse.ToDetailQuizSession(s),
			)
		}
	}

	// Convert map ke slice
	result := make([]quizsessionresponse.ListQuestionSessionResponse, 0, len(grouped))
	for _, v := range grouped {
		result = append(result, *v)
	}

	return result, nil
}
