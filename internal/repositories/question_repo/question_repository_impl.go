package questionrepo

import (
	"context"
	"giat-cerika-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuestionRepositoryImpl struct {
	db *gorm.DB
}

func NewQuestionRepositoryImpl(db *gorm.DB) IQuestionRepository {
	return &QuestionRepositoryImpl{db: db}
}

func (q *QuestionRepositoryImpl) preloadRelations(db *gorm.DB) *gorm.DB {
	return db.
		Preload("Quiz").
		Preload("Quiz.QuizType").
		Preload("Answers")
}

// CreateQuestion implements IQuestionRepository.
func (q *QuestionRepositoryImpl) CreateQuestion(ctx context.Context, data *models.Question) error {
	return q.db.WithContext(ctx).Create(data).Error
}

// FindAllQuestions implements IQuestionRepository.
func (q *QuestionRepositoryImpl) FindAllQuestions(ctx context.Context, limit int, offset int, search string) ([]*models.Question, int, error) {
	var (
		questions []*models.Question
		count     int64
	)
	query := q.db.WithContext(ctx).Model(&models.Question{})
	if search != "" {
		query = query.Where("question_text ILIKE ?", "%"+search+"%")
	}
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	query = q.preloadRelations(query)
	if err := query.Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&questions).Error; err != nil {
		return nil, 0, err
	}
	return questions, int(count), nil
}

// FindQuestionById implements IQuestionRepository.
func (q *QuestionRepositoryImpl) FindQuestionById(ctx context.Context, questionId uuid.UUID) (*models.Question, error) {
	var question models.Question
	if err := q.preloadRelations(q.db.WithContext(ctx)).
		First(&question, "id = ?", questionId).Error; err != nil {
		return nil, err
	}
	return &question, nil
}

// UpdateQuestion implements IQuestionRepository.
func (q *QuestionRepositoryImpl) UpdateQuestion(ctx context.Context, questionId uuid.UUID, data *models.Question) error {
	return q.db.WithContext(ctx).Model(&models.Question{}).Where("id = ?", questionId).Updates(data).Error
}

// DeleteQuestion implements IQuestionRepository.
func (q *QuestionRepositoryImpl) DeleteQuestion(ctx context.Context, questionId uuid.UUID) error {
	return q.db.WithContext(ctx).Delete(&models.Question{}, "id = ?", questionId).Error
}

// UpdateImageQuestion implements IQuestionRepository.
func (q *QuestionRepositoryImpl) UpdateImageQuestion(ctx context.Context, questionId uuid.UUID, image string) error {
	return q.db.WithContext(ctx).Model(&models.Question{}).Where("id = ?", questionId).Update("question_image", image).Error
}
