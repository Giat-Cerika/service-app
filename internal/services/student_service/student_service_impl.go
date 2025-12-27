package studentservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"giat-cerika-service/configs"
	datasources "giat-cerika-service/internal/dataSources"
	studentrequest "giat-cerika-service/internal/dto/request/student_request"
	"giat-cerika-service/internal/models"
	classrepo "giat-cerika-service/internal/repositories/class_repo"
	studentrepo "giat-cerika-service/internal/repositories/student_repo"
	errorresponse "giat-cerika-service/pkg/constant/error_response"
	rabbitmq "giat-cerika-service/pkg/constant/rabbitMq"
	"giat-cerika-service/pkg/utils"
	"giat-cerika-service/pkg/workers/payload"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type StudentServiceImpl struct {
	studenRepo studentrepo.IStudentRepository
	classRepo  classrepo.IClassRepository
	rdb        *redis.Client
	cld        datasources.CloudinaryService
}

func NewStudentServiceImpl(studentRepo studentrepo.IStudentRepository, classRepo classrepo.IClassRepository, rdb *redis.Client, cld datasources.CloudinaryService) IStudentService {
	return &StudentServiceImpl{studenRepo: studentRepo, classRepo: classRepo, rdb: rdb, cld: cld}
}

func fileStudentToBytes(fh *multipart.FileHeader) ([]byte, error) {
	f, err := fh.Open()
	if err != nil {
		return nil, err
	}

	defer f.Close()
	return io.ReadAll(f)
}

var PublishImageAsync = func(p payload.ImageUploadPayload) {
	go func() {
		_ = rabbitmq.PublishToQueue(
			"",
			rabbitmq.SendImageProfileStudentQueueName,
			p,
		)
	}()
}

func (c *StudentServiceImpl) invalidateCacheToothBrush(ctx context.Context) {
	iter := c.rdb.Scan(ctx, 0, "toothbrush:*", 0).Iterator()
	for iter.Next(ctx) {
		c.rdb.Del(ctx, iter.Val())
	}
}

// Register implements IStudentService.
func (s *StudentServiceImpl) Register(ctx context.Context, req studentrequest.RegisterStudentRequest) error {
	uniqueUsername, err := s.studenRepo.FindUsernameUnique(ctx, req.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get unique username", 500)
	}
	uniqueNisn, err := s.studenRepo.FindNisnUnique(ctx, req.Nisn)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get unique nisn", 500)
	}

	if strings.TrimSpace(req.Name) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Name is required", 400)
	}
	if strings.TrimSpace(req.Username) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Username is required", 400)
	}
	if strings.TrimSpace(req.Password) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Password is required", 400)
	}
	if strings.TrimSpace(req.ConfirmPassword) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Confirm Password is required", 400)
	}
	if req.Password != req.ConfirmPassword {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Password and Confirm Password doesn't match", 409)
	}
	if strings.TrimSpace(req.Nisn) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Nisn is required", 400)
	}
	if req.DateOfBirth.IsZero() {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Date Of Birth is required", 400)
	}
	if req.Age == nil {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Age is required", 400)
	}
	if req.Photo == nil {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Photo is required", 400)
	}
	if req.ClassID == uuid.Nil {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "class is required", 400)
	}
	if uniqueUsername != "" {
		return errorresponse.NewCustomError(errorresponse.ErrExists, "Username already exists", 409)
	}
	if uniqueNisn != "" {
		return errorresponse.NewCustomError(errorresponse.ErrExists, "Nisn already exists", 409)
	}

	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "failed to hashing password", 400)
	}

	role, err := s.studenRepo.FindRoleStudent(ctx)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get role", 500)
	}

	class, err := s.classRepo.FindById(ctx, req.ClassID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorresponse.NewCustomError(errorresponse.ErrNotFound, "class not found", 404)
		}
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get class", 500)
	}

	newStudent := &models.User{
		ID:          uuid.New(),
		Name:        &req.Name,
		Username:    req.Username,
		Password:    hashed,
		Nisn:        &req.Nisn,
		DateOfBirth: &req.DateOfBirth,
		Age:         *req.Age,
		RoleID:      role.ID,
		ClassID:     &class.ID,
		Status:      1,
	}

	if err := s.studenRepo.Create(ctx, newStudent); err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to create student", 500)
	}

	if req.Photo != nil {
		if binner, err := fileStudentToBytes(req.Photo); err == nil && len(binner) > 0 {
			pay := payload.ImageUploadPayload{
				ID:        newStudent.ID,
				Type:      "single",
				FileBytes: binner,
				Folder:    "giat_ceria/photo_student",
				Filename:  fmt.Sprintf("student_%s_photo", newStudent.ID.String()),
			}

			_ = rabbitmq.PublishToQueue("", rabbitmq.SendImageProfileStudentQueueName, pay)
		}
	}

	return nil
}

// Login implements IStudentService.
func (s *StudentServiceImpl) Login(ctx context.Context, req studentrequest.LoginStudentRequet) (string, error) {
	student, err := s.studenRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return "", errorresponse.NewCustomError(errorresponse.ErrBadRequest, "invalid credential", 400)
	}

	isPassword := utils.CheckPasswordHash(req.Password, student.Password)
	if !isPassword {
		return "", errorresponse.NewCustomError(errorresponse.ErrBadRequest, "password incorrect", 400)
	}

	token, err := utils.GenerateToken(student.ID.String(), student.Role.Name)
	if err != nil {
		return "", errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to generate token", 500)
	}

	expiry, err := utils.GetExpiryFromToken(token)
	if err != nil {
		return "", errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get expiry token", 500)
	}

	redisKey := fmt.Sprintf("student_token:%s", student.ID)
	if err := configs.SetRedis(ctx, redisKey, token, time.Until(expiry)); err != nil {
		return "", errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to store token in cache", 400)
	}

	return token, nil
}

// GetProfile implements IStudentService.
func (s *StudentServiceImpl) GetProfile(ctx context.Context, studentId uuid.UUID, token string) (*models.User, error) {
	cacheKey := fmt.Sprintf("student_token:%s", studentId)
	storedToken, err := configs.GetRedis(ctx, cacheKey)
	if err != nil || storedToken != token {
		return nil, errorresponse.NewCustomError(errorresponse.ErrUnauthorized, "unauthorized access", 401)
	}

	student, err := s.studenRepo.FindByStudentID(ctx, studentId)
	if err != nil {
		return nil, errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get student", 500)
	}

	return student, nil
}

// CheckTokenBlacklisted implements IStudentService.
func (s *StudentServiceImpl) CheckTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	blackListeed := fmt.Sprintf("blacklistToken_student:%s", token)
	val, err := configs.GetRedis(ctx, blackListeed)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		return false, errorresponse.NewCustomError(errorresponse.ErrNotFound, "blacklisted token not found", 404)
	}
	return val == "blacklister", nil
}

// Logoutctx implements IStudentService.
func (s *StudentServiceImpl) Logout(ctx context.Context, studentID uuid.UUID, token string) error {
	expiry, err := utils.GetExpiryFromToken(token)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get expiry token", 500)
	}
	blacklistedKey := fmt.Sprintf("blacklistToken_student:%s", token)
	err = configs.SetRedis(ctx, blacklistedKey, "blacklister", time.Until(expiry))
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to store blacklisted token in cache", 500)
	}
	return nil
}

// CheckNisnAndDateOfBirth implements IStudentService.
func (s *StudentServiceImpl) CheckNisnAndDateOfBirth(ctx context.Context, req studentrequest.CheckNisnAndDateOfBirth) (*models.User, error) {
	if strings.TrimSpace(req.Nisn) == "" {
		return nil, errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Nisn is required", 400)
	}
	if req.DateOfBirth.IsZero() {
		return nil, errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Date Of Birth is required", 400)
	}
	student, err := s.studenRepo.CheckNisnAndDateOfBirth(ctx, req.Nisn, req.DateOfBirth)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorresponse.NewCustomError(errorresponse.ErrNotFound, "student not found", 404)
		}
		return nil, errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get student", 500)
	}
	if student == nil {
		return nil, errorresponse.NewCustomError(errorresponse.ErrNotFound, "student is empty", 404)
	}

	redisKey := fmt.Sprintf("id_forgot_password:%s", student.ID)
	data, _ := json.Marshal(map[string]any{
		"student_id":    student.ID,
		"name":          student.Name,
		"nisn":          req.Nisn,
		"date_of_birth": req.DateOfBirth,
	})
	_ = configs.SetRedis(ctx, redisKey, data, time.Minute*5)

	return student, nil
}

// UpdateNewPasswordStudent implements IStudentService.
func (s *StudentServiceImpl) UpdateNewPasswordStudent(ctx context.Context, studentID uuid.UUID, req studentrequest.UpdatePassword) error {
	if strings.TrimSpace(req.NewPassword) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "new password is required", 400)
	}

	if strings.TrimSpace(req.ConfirmPassword) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "confirm password is required", 400)
	}

	if req.StudentID == uuid.Nil {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "student id is required", 400)
	}

	if strings.TrimSpace(req.NewPassword) != strings.TrimSpace(req.ConfirmPassword) {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "new password and confirm password doesn't match", 400)
	}

	redisKey := fmt.Sprintf("id_forgot_password:%s", studentID)
	cacheData, err := configs.GetRedis(ctx, redisKey)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return errorresponse.NewCustomError(errorresponse.ErrUnauthorized, "reset password session expired", 401)
		}
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get cache", 500)
	}

	var stored map[string]any
	_ = json.Unmarshal([]byte(cacheData), &stored)

	hashed, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed hashing password", 500)
	}

	if err := s.studenRepo.UpdateNewPassword(ctx, studentID, hashed); err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to update password", 500)
	}

	_ = configs.DeleteRedis(ctx, redisKey)

	return nil
}

// UpdateProfileStudent implements IStudentService.
func (s *StudentServiceImpl) UpdateProfileStudent(ctx context.Context, studentId uuid.UUID, req studentrequest.UpdateProfileRequest) error {
	student, err := s.studenRepo.FindByStudentID(ctx, studentId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorresponse.NewCustomError(errorresponse.ErrNotFound, "student not found", 404)
		}
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get student", 500)
	}

	// 2. Cek unique Username, pastikan bukan milik dirinya sendiri
	existUsername, err := s.studenRepo.FindByUsername(ctx, req.Username)
	if err == nil && existUsername.ID != studentId {
		return errorresponse.NewCustomError(errorresponse.ErrExists, "Username already exists", 409)
	}

	// 3. Cek unique NISN, pastikan bukan miliknya sendiri
	existNisn, err := s.studenRepo.FindNisnUnique(ctx, req.Nisn)
	if err == nil && existNisn != *student.Nisn {
		return errorresponse.NewCustomError(errorresponse.ErrExists, "NISN already exists", 409)
	}

	class, err := s.classRepo.FindById(ctx, req.ClassID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorresponse.NewCustomError(errorresponse.ErrNotFound, "class not found", 404)
		}
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get class", 500)
	}

	if req.Username != "" {
		student.Username = req.Username
	}
	if req.Nisn != "" {
		student.Nisn = &req.Nisn
	}
	if req.Name != "" {
		student.Name = &req.Name
	}
	if req.Age != 0 {
		student.Age = req.Age
	}
	if !req.DateOfBirth.IsZero() {
		student.DateOfBirth = &req.DateOfBirth
	}
	if req.ClassID != uuid.Nil {
		student.ClassID = &class.ID
	}

	err = s.studenRepo.UpdateProfile(ctx, studentId, student)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to update student", 500)
	}

	return nil
}

// UpdatePhotoStudent implements IStudentService.
func (s *StudentServiceImpl) UpdatePhotoStudent(ctx context.Context, studentId uuid.UUID, photo *multipart.FileHeader) error {
	student, err := s.studenRepo.FindByStudentID(ctx, studentId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorresponse.NewCustomError(errorresponse.ErrNotFound, "student not found", 404)
		}
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get student", 500)
	}

	if photo != nil {
		if bin, err := fileStudentToBytes(photo); err == nil && len(bin) > 0 {
			task := payload.ImageUploadPayload{
				ID:        student.ID,
				Type:      "single",
				FileBytes: bin,
				Folder:    "giat_ceria/photo_student",
				Filename:  fmt.Sprintf("student_%s_photo", studentId.String()),
			}
			_ = rabbitmq.PublishToQueue("", rabbitmq.SendImageProfileStudentQueueName, task)
		}
	}

	return nil
}

// CreateTootBrushStudent implements IStudentService.
func (s *StudentServiceImpl) CreateTootBrushStudent(ctx context.Context, studentId uuid.UUID, req studentrequest.CreateTootBrushRequest) error {
	if studentId == uuid.Nil {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "student id is required", 400)
	}

	if strings.TrimSpace(req.TimeType) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "type time is required", 400)
	}

	if req.TimeType != strings.ToUpper("MORNING") && req.TimeType != strings.ToUpper("NIGHT") {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "type time only 'MORNING' or 'NIGHT'", 400)
	}

	locJakarta, _ := time.LoadLocation("Asia/Jakarta")
	nowJakarta := time.Now().In(locJakarta)

	hour := nowJakarta.Hour()
	minute := nowJakarta.Minute()
	timeType := strings.ToUpper(req.TimeType)

	if timeType == "MORNING" {
		if hour < 5 || (hour == 7 && minute > 0) || hour > 7 {
			return errorresponse.NewCustomError(
				errorresponse.ErrBadRequest,
				"absen pagi hanya bisa antara jam 05:00 sampai 10:00",
				400,
			)
		}
	}

	if timeType == "NIGHT" {
		if hour < 17 || (hour == 22 && minute > 0) || hour > 22 {
			return errorresponse.NewCustomError(
				errorresponse.ErrBadRequest,
				"absen malam hanya bisa antara jam 17:00 sampai 22:00",
				400,
			)
		}
	}

	logDate := nowJakarta.Format("2006-01-02")

	exists, err := s.studenRepo.CheckTootBrushExists(ctx, studentId, strings.ToUpper(req.TimeType), nowJakarta)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to check log", 500)
	}

	if exists {
		return errorresponse.NewCustomError(errorresponse.ErrExists, "Anda sudah absen untuk sesi ini hari ini", 409)
	}

	student, err := s.studenRepo.FindByStudentID(ctx, studentId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorresponse.NewCustomError(errorresponse.ErrNotFound, "student not found", 404)
		}
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get student", 500)
	}

	fmt.Println(logDate)

	newLog := &models.ToootBrushLog{
		ID:        uuid.New(),
		UserID:    student.ID,
		TimeType:  timeType,
		LogDate:   logDate,
		LogTime:   nowJakarta,
		CreatedAt: nowJakarta,
	}

	if err := s.studenRepo.CreateTootBrush(ctx, student.ID, newLog); err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to create log", 500)
	}

	s.invalidateCacheToothBrush(ctx)
	return nil
}

// GetHitoryToothBrush implements IStudentService.
func (s *StudentServiceImpl) GetHitoryToothBrush(ctx context.Context, studentId uuid.UUID, typeTime string, page int, limit int) ([]*models.ToootBrushLog, int, error) {
	cacheKey := fmt.Sprintf("toothbrush:%s:type:%s:page:%d:limit:%d", studentId, typeTime, page, limit)

	// GET FROM CACHE
	if cached, err := configs.GetRedis(ctx, cacheKey); err == nil && len(cached) > 0 {
		var result struct {
			Data  []*models.ToootBrushLog `json:"data"`
			Total int                     `json:"total"`
		}
		if json.Unmarshal([]byte(cached), &result) == nil {
			return result.Data, result.Total, nil
		}
	}

	offset := (page - 1) * limit
	items, total, err := s.studenRepo.GetHistoryTootBrush(ctx, studentId, typeTime, limit, offset)
	if err != nil {
		return nil, 0, errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to get history toothbrush", 500)
	}

	if len(items) == 0 {
		items = []*models.ToootBrushLog{}
	}

	buf, _ := json.Marshal(map[string]any{
		"data":  items,
		"total": total,
	})
	_ = configs.SetRedis(ctx, cacheKey, buf, time.Minute*30)

	return items, total, nil
}

func (s *StudentServiceImpl) GetAllStudents(ctx context.Context, search string) ([]*models.User, int, error) {
	students, total, err := s.studenRepo.GetAllStudents(ctx, search)
	if err != nil {
		return nil, 0, errorresponse.NewCustomError(
			errorresponse.ErrInternal,
			"failed to get students",
			http.StatusInternalServerError)
	}
	if students == nil {
		students = []*models.User{}
	}
	return students, total, nil
}
