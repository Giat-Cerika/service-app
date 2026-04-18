package quizhistoryservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"giat-cerika-service/configs"
	quizhistoryresponse "giat-cerika-service/internal/dto/response/quiz_history_response"
	"giat-cerika-service/internal/models"
	quizhistoryrepo "giat-cerika-service/internal/repositories/quiz_history_repo"
	quizrepo "giat-cerika-service/internal/repositories/quiz_repo"
	studentrepo "giat-cerika-service/internal/repositories/student_repo"
	errorresponse "giat-cerika-service/pkg/constant/error_response"
	"giat-cerika-service/pkg/utils"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type QuizHistoryServiceImpl struct {
	quizHistoryRepo quizhistoryrepo.IQuizHistoryRepository
	studentRepo     studentrepo.IStudentRepository
	quizRepo        quizrepo.IQuizRepository
	rdb             *redis.Client
}

func NewQuizHistoryServiceImpl(quizHistoryRepo quizhistoryrepo.IQuizHistoryRepository, studentRepo studentrepo.IStudentRepository, quizRepo quizrepo.IQuizRepository, rdb *redis.Client) IQuizHistoryService {
	return &QuizHistoryServiceImpl{quizHistoryRepo: quizHistoryRepo, studentRepo: studentRepo, quizRepo: quizRepo, rdb: rdb}
}

// GetHistoryQuizStudent implements [IQuizHistoryService].
func (q QuizHistoryServiceImpl) GetHistoryQuizStudent(
	ctx context.Context,
	userId uuid.UUID,
	search string,
) ([]quizhistoryresponse.QuizHistoryResponse, error) {

	student, err := q.studentRepo.FindByStudentID(ctx, userId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorresponse.NewCustomError(errorresponse.ErrNotFound, "student not found", 404)
		}
		return nil, errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get student", 500)
	}

	cacheKey := fmt.Sprintf("quizHistory:%s:search:%s", userId, search)
	if cached, err := configs.GetRedis(ctx, cacheKey); err == nil && len(cached) > 0 {
		var data []quizhistoryresponse.QuizHistoryResponse
		if json.Unmarshal([]byte(cached), &data) == nil {
			return data, nil
		}
	}

	items, err := q.quizHistoryRepo.FindHistoryByUserID(ctx, student.ID, search)
	if err != nil {
		return nil, errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get history quiz student", 500)
	}

	// --- Kumpulkan semua quiz ID unik dari seluruh history ---
	// Lalu ambil sekaligus dalam 1 query (menghindari N+1).
	quizIDSet := make(map[uuid.UUID]struct{}, len(items))
	for _, h := range items {
		quizIDSet[h.QuizID] = struct{}{}
	}
	quizIDs := make([]uuid.UUID, 0, len(quizIDSet))
	for id := range quizIDSet {
		quizIDs = append(quizIDs, id)
	}
	quizList, _ := q.quizRepo.FindByIds(ctx, quizIDs)
	// Bangun map QuizID → AmountAssigned untuk lookup O(1)
	quizMap := make(map[uuid.UUID]int, len(quizList))
	for _, qz := range quizList {
		quizMap[qz.ID] = qz.AmountAssigned
	}

	res := make([]quizhistoryresponse.QuizHistoryResponse, 0, len(items))

	for _, h := range items {
		// snapshot sebagai default
		currentAssigned := h.AmountAssigned

		// kalau quiz masih ada di DB → pakai nilai global (dari map, O(1))
		if assigned, ok := quizMap[h.QuizID]; ok {
			currentAssigned = assigned
		}

		res = append(res, quizhistoryresponse.QuizHistoryResponse{
			ID:              h.ID,
			Title:           h.Title,
			Description:     h.Description,
			StartDate:       utils.FormatDateTime(h.StartDate),
			EndDate:         utils.FormatDateTime(h.EndDate),
			AmountQuestions: h.AmountQuestions,
			AmountAssigned:  currentAssigned,
			Score:           h.Score,
			MaxScore:        h.MaxScore,
			Percentage:      h.Percentage,
			Status:          string(h.Status),
			StartedAt:       utils.FormatDateTime(h.StartedAt),
			CompletedAt:     utils.FormatDateTime(h.CompletedAt),
			StatusCategory:  h.StatusCategory,
			CreatedAt:       utils.FormatDate(h.CreatedAt),
			UpdatedAt:       utils.FormatDate(h.UpdatedAt),
		})
	}

	buf, _ := json.Marshal(res)
	_ = configs.SetRedis(ctx, cacheKey, buf, time.Minute*30)

	return res, nil
}

// GetAllHistoryQuestionByQuizHistory implements [IQuizHistoryService].
func (q *QuizHistoryServiceImpl) GetAllHistoryQuestionByQuizHistory(ctx context.Context, quizHistoryId uuid.UUID) ([]*models.QuestionHistory, error) {
	cacheKey := fmt.Sprintf("questions_history:quiz_history:%s", quizHistoryId)

	if cached, err := configs.GetRedis(ctx, cacheKey); err == nil && len(cached) > 0 {
		var data []*models.QuestionHistory
		if json.Unmarshal([]byte(cached), &data) == nil {
			return data, nil
		}
	}

	quizHistory, err := q.quizHistoryRepo.FindQuizHistoryById(ctx, quizHistoryId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorresponse.NewCustomError(errorresponse.ErrNotFound, "quiz history not found", 404)
		}
		return nil, errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get quiz history", 500)
	}

	items, err := q.quizHistoryRepo.FindAllQuestionHistory(ctx, quizHistory.ID)
	if err != nil {
		return nil, errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get question history", 500)
	}

	if items == nil {
		items = []*models.QuestionHistory{}
	}

	buf, _ := json.Marshal(map[string]any{
		"data": items,
	})

	_ = configs.SetRedis(ctx, cacheKey, buf, time.Minute*30)

	return items, nil
}

// GetHistoryQuizByQuizID implements [IQuizHistoryService].
func (q *QuizHistoryServiceImpl) GetHistoryQuizByQuizID(ctx context.Context) ([]quizhistoryresponse.QuizHistoryGroupAdminResponse, error) {
	cacheKey := fmt.Sprintln("quizHistory:all")
	if cached, err := configs.GetRedis(ctx, cacheKey); err == nil && len(cached) > 0 {
		var data []quizhistoryresponse.QuizHistoryGroupAdminResponse
		if json.Unmarshal([]byte(cached), &data) == nil {
			return data, nil
		}
	}
	items, err := q.quizHistoryRepo.FindHistoryByQuizID(ctx)
	if err != nil {
		return nil, errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get history quiz student", 500)
	}

	// --- Kumpulkan semua quiz ID & user ID unik ---
	// Ambil sekaligus dalam 2 query bulk (menggantikan N+1+N+1).
	quizIDSet := make(map[uuid.UUID]struct{}, len(items))
	userIDSet := make(map[uuid.UUID]struct{}, len(items))
	for _, h := range items {
		quizIDSet[h.QuizID] = struct{}{}
		userIDSet[h.UserID] = struct{}{}
	}

	quizIDs := make([]uuid.UUID, 0, len(quizIDSet))
	for id := range quizIDSet {
		quizIDs = append(quizIDs, id)
	}
	userIDs := make([]uuid.UUID, 0, len(userIDSet))
	for id := range userIDSet {
		userIDs = append(userIDs, id)
	}

	quizList, _ := q.quizRepo.FindByIds(ctx, quizIDs)
	quizMap := make(map[uuid.UUID]*models.Quiz, len(quizList))
	for _, qz := range quizList {
		quizMap[qz.ID] = qz
	}

	studentList, _ := q.studentRepo.FindByUserIDs(ctx, userIDs)
	studentMap := make(map[uuid.UUID]*models.User, len(studentList))
	for _, s := range studentList {
		studentMap[s.ID] = s
	}

	grouped := make(map[uuid.UUID]*quizhistoryresponse.QuizHistoryGroupAdminResponse)

	for _, h := range items {
		quizID := h.QuizID

		// ===== INIT GROUP (1 QUIZ) =====
		if _, ok := grouped[quizID]; !ok {
			quiz, exists := quizMap[quizID]
			if !exists {
				continue
			}

			grouped[quizID] = &quizhistoryresponse.QuizHistoryGroupAdminResponse{
				QuizID:          quiz.ID,
				Title:           quiz.Title,
				Description:     quiz.Description,
				StartDate:       utils.FormatDateTime(&quiz.StartDate),
				EndDate:         utils.FormatDateTime(&quiz.EndDate),
				DetailHistories: []quizhistoryresponse.QuizHistoryDetailAdminResponse{},
			}
		}

		// ===== DETAIL HISTORY =====
		student, exists := studentMap[h.UserID]
		if !exists {
			continue
		}

		grouped[quizID].DetailHistories = append(
			grouped[quizID].DetailHistories,
			quizhistoryresponse.QuizHistoryDetailAdminResponse{
				ID:             h.ID,
				StudentName:    *student.Name,
				Class:          student.Class.NameClass,
				Score:          h.Score,
				MaxScore:       h.MaxScore,
				Percentage:     h.Percentage,
				Status:         string(h.Status),
				StatusCategory: h.StatusCategory,
				StartedAt:      utils.FormatDateTime(h.StartedAt),
				CompletedAt:    utils.FormatDateTime(h.CompletedAt),
				CreatedAt:      utils.FormatDate(h.CreatedAt),
			},
		)
	}

	// ===== MAP → SLICE =====
	result := make([]quizhistoryresponse.QuizHistoryGroupAdminResponse, 0, len(grouped))
	for _, v := range grouped {
		result = append(result, *v)
	}

	buf, _ := json.Marshal(result)
	_ = configs.SetRedis(ctx, cacheKey, buf, time.Minute*30)

	return result, nil
}
