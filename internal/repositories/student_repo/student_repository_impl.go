package studentrepo

import (
	"context"
	"giat-cerika-service/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StudentRepositoryImpl struct {
	db *gorm.DB
}

func NewStudentRepositoryImpl(db *gorm.DB) IStudentRepository {
	return &StudentRepositoryImpl{db: db}
}

// Create implements IStudentRepository.
func (s *StudentRepositoryImpl) Create(ctx context.Context, data *models.User) error {
	return s.db.WithContext(ctx).Create(data).Error
}

// FindUsernameUnique implements IStudentRepository.
func (s *StudentRepositoryImpl) FindUsernameUnique(ctx context.Context, username string) (string, error) {
	var student models.User
	if err := s.db.WithContext(ctx).Preload("Role").Preload("Class").First(&student, "username = ?", username).Error; err != nil {
		return "", err
	}

	return student.Username, nil
}

// FindNisnUnique implements IStudentRepository.
func (s *StudentRepositoryImpl) FindNisnUnique(ctx context.Context, nisn string) (string, error) {
	var student models.User
	if err := s.db.WithContext(ctx).Preload("Role").Preload("Class").First(&student, "nisn = ?", nisn).Error; err != nil {
		return "", err
	}

	return *student.Nisn, nil
}

// FindRoleStudent implements IStudentRepository.
func (s *StudentRepositoryImpl) FindRoleStudent(ctx context.Context) (*models.Role, error) {
	var role models.Role
	if err := s.db.WithContext(ctx).First(&role, "name = ?", "student").Error; err != nil {
		return nil, err
	}

	return &role, nil
}

// UpdatePhotoStudent implements IStudentRepository.
func (s *StudentRepositoryImpl) UpdatePhotoStudent(ctx context.Context, studentId uuid.UUID, photo string) error {
	subQuery := s.db.Select("id").Where("name = ?", "student").Table("roles")
	return s.db.WithContext(ctx).Model(&models.User{}).
		Where("id = ?", studentId).
		Where("role_id IN (?)", subQuery).
		Update("photo", photo).Error
}

// FindByUsername implements IStudentRepository.
func (s *StudentRepositoryImpl) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var student models.User
	if err := s.db.WithContext(ctx).Preload("Role").Preload("Class").First(&student, "username = ?", username).Error; err != nil {
		return nil, err
	}

	return &student, nil
}

// FindByStudentID implements IStudentRepository.
func (s *StudentRepositoryImpl) FindByStudentID(ctx context.Context, studentID uuid.UUID) (*models.User, error) {
	var student models.User
	if err := s.db.WithContext(ctx).Preload("Role").Preload("Class").First(&student, "id = ?", studentID).Error; err != nil {
		return nil, err
	}

	return &student, nil
}

// CheckNisnAndDateOfBirth implements IStudentRepository.
func (s *StudentRepositoryImpl) CheckNisnAndDateOfBirth(ctx context.Context, nisn string, dateOfBirth time.Time) (*models.User, error) {
	var student models.User
	if err := s.db.WithContext(ctx).Preload("Role").Preload("Class").First(&student, "nisn = ? AND date_of_birth = ?", nisn, dateOfBirth).Error; err != nil {
		return nil, err
	}

	return &student, nil
}

// UpdateNewPassword implements IStudentRepository.
func (s *StudentRepositoryImpl) UpdateNewPassword(ctx context.Context, studentID uuid.UUID, password string) error {
	if err := s.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", studentID).Update("password", password).Error; err != nil {
		return err
	}

	return nil
}

func (s *StudentRepositoryImpl) UpdateProfile(ctx context.Context, studentId uuid.UUID, data *models.User) error {
	updateData := map[string]interface{}{
		"name":          data.Name,
		"username":      data.Username,
		"nisn":          data.Nisn,
		"age":           data.Age,
		"class_id":      data.ClassID,
		"date_of_birth": data.DateOfBirth,
	}

	return s.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", studentId).
		Updates(updateData).Error
}

// CreateTootBrush implements IStudentRepository.
func (s *StudentRepositoryImpl) CreateTootBrush(ctx context.Context, studentId uuid.UUID, data *models.ToootBrushLog) error {
	return s.db.WithContext(ctx).Create(data).Error
}

// CheckTootBrushExists implements IStudentRepository.
func (r *StudentRepositoryImpl) CheckTootBrushExists(ctx context.Context, studentId uuid.UUID, timeType string, logDate time.Time) (bool, error) {
	var count int64

	// Pastikan logDate sudah dalam UTC
	// Query dengan DATE() function yang consistent
	formattedDate := logDate.Format("2006-01-02")

	err := r.db.WithContext(ctx).
		Model(&models.ToootBrushLog{}).
		Where("user_id = ?", studentId).
		Where("time_type = ?", timeType).
		Where("DATE(log_date AT TIME ZONE 'UTC') = DATE(? AT TIME ZONE 'UTC')", formattedDate).
		Count(&count).Error

	return count > 0, err
}

// GetHistoryTootBrush implements IStudentRepository.
func (s *StudentRepositoryImpl) GetHistoryTootBrush(ctx context.Context, studentId uuid.UUID, typeTime string, limit int, offset int) ([]*models.ToootBrushLog, int, error) {
	var (
		logs  []*models.ToootBrushLog
		count int64
	)

	query := s.db.WithContext(ctx).Model(&models.ToootBrushLog{}).
		Where("user_id = ?", studentId)

	if typeTime != "" {
		query = query.Where("time_type = ?", typeTime)
	}

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Preload("User").
		Preload("User.Role").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, int(count), nil
}

// FindAll implements IStudentRepository.
func (q *StudentRepositoryImpl) FindAllStudents(ctx context.Context, limit, offset int, search string) ([]*models.User, int, error) {
	var (
		students []*models.User
		count    int64
	)

	query := q.db.WithContext(ctx).Model(&models.User{})
	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Find(&students).Error; err != nil {
		return nil, 0, err
	}

	return students, int(count), nil
}