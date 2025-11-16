package questionnairerepo

import (
	"context"
	"giat-cerika-service/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type QuestionnaireRepositoryImpl struct {
	db *gorm.DB
}

func NewQuestionnaireRepositoryImpl(db *gorm.DB) IQuestionnaireRepository {
	return &QuestionnaireRepositoryImpl{db: db}
}

// Create implements IQuestionnaireRepository.
func (c *QuestionnaireRepositoryImpl) Create(ctx context.Context, data *models.Questionnaire) error {
	return c.db.WithContext(ctx).Create(data).Error
}

// FindByTitle implements IQuestionnaireRepository.
func (c *QuestionnaireRepositoryImpl) FindByTitle(ctx context.Context, Title string) (*models.Questionnaire, error) {
	var questionnaire models.Questionnaire
	if err := c.db.WithContext(ctx).First(&questionnaire, "title = ?", Title).Error; err != nil {
		return nil, err
	}

	return &questionnaire, nil
}

// FindAll implements IQuestionnaireRepository.
func (c *QuestionnaireRepositoryImpl) FindAll(ctx context.Context, limit int, offset int, search string) ([]*models.Questionnaire, int, error) {
	var (
		questionnairees []*models.Questionnaire
		count   int64
	)

	query := c.db.WithContext(ctx).Model(&models.Questionnaire{})
	if search != "" {
		query = query.Where("name_questionnaire ILIKE ?", "%"+search+"%")
	}
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&questionnairees).Error; err != nil {
		return nil, 0, err
	}

	return questionnairees, int(count), nil
}

// FindById implements IQuestionnaireRepository.
func (c *QuestionnaireRepositoryImpl) FindById(ctx context.Context, questionnaireId uuid.UUID) (*models.Questionnaire, error) {
	var questionnaire models.Questionnaire
	if err := c.db.WithContext(ctx).First(&questionnaire, "id = ?", questionnaireId).Error; err != nil {
		return nil, err
	}

	return &questionnaire, nil
}

// Update implements IQuestionnaireRepository.
func (c *QuestionnaireRepositoryImpl) Update(ctx context.Context, questionnaireId uuid.UUID, data *models.Questionnaire) error {
	return c.db.WithContext(ctx).Save(data).Error
}

// Delete implements IQuestionnaireRepository.
func (c *QuestionnaireRepositoryImpl) Delete(ctx context.Context, questionnaireId uuid.UUID) error {
	return c.db.WithContext(ctx).Delete(&models.Questionnaire{}, "id = ?", questionnaireId).Error
}
