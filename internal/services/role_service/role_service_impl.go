package roleservice

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"giat-cerika-service/configs"
	rolerequest "giat-cerika-service/internal/dto/request/role_request"
	"giat-cerika-service/internal/models"
	rolerepo "giat-cerika-service/internal/repositories/role_repo"
	errorresponse "giat-cerika-service/pkg/constant/error_response"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type RoleServiceImpl struct {
	roleRepo rolerepo.IRoleRepository
	rdb      *redis.Client
}

func NewRoleServiceImpl(roleRepo rolerepo.IRoleRepository, rdb *redis.Client) IRoleService {
	return &RoleServiceImpl{roleRepo: roleRepo, rdb: rdb}
}

func (r *RoleServiceImpl) invalidateCacheRole(ctx context.Context) {
	iter := r.rdb.Scan(ctx, 0, "roles:*", 0).Iterator()
	for iter.Next(ctx) {
		r.rdb.Del(ctx, iter.Val())
	}

	iterID := r.rdb.Scan(ctx, 0, "role:*", 0).Iterator()
	for iterID.Next(ctx) {
		r.rdb.Del(ctx, iterID.Val())
	}
}

// CreateRole implements IRoleService.
func (r *RoleServiceImpl) CreateRole(ctx context.Context, req rolerequest.CreateRoleRequest) error {
	existRole, err := r.roleRepo.FindByName(ctx, req.Name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "Failed to Get Role Name", 500)
	}

	if existRole != nil {
		return errorresponse.NewCustomError(errorresponse.ErrExists, "Role Name Already Exists", 409)
	}

	if strings.TrimSpace(req.Name) == "" {
		return errorresponse.NewCustomError(errorresponse.ErrBadRequest, "Name is required", 400)
	}

	newRole := &models.Role{
		ID:   uuid.New(),
		Name: req.Name,
	}

	err = r.roleRepo.Create(ctx, newRole)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "Failed to create role", 500)
	}

	r.invalidateCacheRole(ctx)

	return nil
}

// GetAllRole implements IRoleService.
func (r *RoleServiceImpl) GetAllRole(ctx context.Context, page int, limit int, search string) ([]*models.Role, int, error) {
	cacheKey := fmt.Sprintf("roles:search:%s:page:%d:limit:%d", search, page, limit)
	if cached, err := configs.GetRedis(ctx, cacheKey); err == nil && len(cached) > 0 {
		var result struct {
			Data  []*models.Role `json:"data"`
			Total int            `json:"total"`
		}
		if json.Unmarshal([]byte(cached), &result) == nil {
			return result.Data, result.Total, nil
		}
	}

	offset := (page - 1) * limit

	items, total, err := r.roleRepo.FindAll(ctx, limit, offset, search)
	if err != nil {
		return nil, 0, errorresponse.NewCustomError(errorresponse.ErrInternal, "Failed to get role", 500)
	}
	if len(items) == 0 {
		items = []*models.Role{}
	}

	buf, _ := json.Marshal(map[string]any{
		"data":  items,
		"total": total,
	})
	_ = configs.SetRedis(ctx, cacheKey, buf, time.Minute*30)

	return items, total, nil
}

// GetByIdRole implements IRoleService.
func (r *RoleServiceImpl) GetByIdRole(ctx context.Context, roleId uuid.UUID) (*models.Role, error) {
	cacheKey := fmt.Sprintf("role:%s", roleId)
	if cached, err := configs.GetRedis(ctx, cacheKey); err == nil && len(cached) > 0 {
		var role models.Role
		if json.Unmarshal([]byte(cached), &role) == nil {
			return &role, nil
		}
	}

	role, err := r.roleRepo.FindById(ctx, roleId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errorresponse.NewCustomError(errorresponse.ErrNotFound, "role not found", 404)
		}
		return nil, errorresponse.NewCustomError(errorresponse.ErrInternal, "role not found", 500)
	}

	buf, _ := json.Marshal(role)
	_ = configs.SetRedis(ctx, cacheKey, buf, time.Minute*30)
	return role, nil
}

// UpdateRole implements IRoleService.
func (r *RoleServiceImpl) UpdateRole(ctx context.Context, roleId uuid.UUID, req rolerequest.UpdateRoleRequest) error {
	role, err := r.roleRepo.FindById(ctx, roleId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorresponse.NewCustomError(errorresponse.ErrNotFound, "role not found", 404)
		}
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "role not found", 500)
	}

	existsRole, err := r.roleRepo.FindByName(ctx, req.Name)
	if err == nil && existsRole.ID != roleId {
		return errorresponse.NewCustomError(errorresponse.ErrExists, "Role Name Already Exists", 404)
	}

	if req.Name != "" {
		role.Name = req.Name
	}

	err = r.roleRepo.Update(ctx, roleId, role)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "Failed to update role", 500)
	}
	r.invalidateCacheRole(ctx)

	return nil
}

// DeleteRole implements IRoleService.
func (r *RoleServiceImpl) DeleteRole(ctx context.Context, roleId uuid.UUID) error {
	_, err := r.roleRepo.FindById(ctx, roleId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errorresponse.NewCustomError(errorresponse.ErrNotFound, "role not found", 404)
		}
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "role not found", 500)
	}

	err = r.roleRepo.Delete(ctx, roleId)
	if err != nil {
		return errorresponse.NewCustomError(errorresponse.ErrInternal, "failed to delete role", 500)
	}

	r.invalidateCacheRole(ctx)

	return nil
}
