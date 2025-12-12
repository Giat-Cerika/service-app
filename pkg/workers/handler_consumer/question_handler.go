package handlerconsumer

import (
	"context"
	"giat-cerika-service/configs"
	"giat-cerika-service/internal/models"
	questionrepo "giat-cerika-service/internal/repositories/question_repo"
	"giat-cerika-service/pkg/workers/payload"

	"github.com/redis/go-redis/v9"
)

type QuestionHandler struct {
	repo questionrepo.IQuestionRepository
	rdb  *redis.Client
}

func NewQuestionHandler() *QuestionHandler {
	return &QuestionHandler{
		repo: questionrepo.NewQuestionRepositoryImpl(configs.DB),
		rdb:  configs.RDB,
	}
}

func (q *QuestionHandler) HandleSingle(ctx context.Context, photoUrl string, payloads any) error {
	p, ok := payloads.(*payload.ImageUploadPayload)
	if !ok {
		return nil
	}
	err := q.repo.UpdateImageQuestion(ctx, p.ID, photoUrl)
	if err != nil {
		return err
	}
	q.deleteCacheQuestion(ctx, p.ID.String())
	return nil
}

func (q *QuestionHandler) HandleMany(ctx context.Context, image *models.Image, payloads any) error {
	p, ok := payloads.(*payload.ImageUploadPayload)
	if !ok {
		return nil
	}
	// err := q.repo.UpdateImageQuestion(ctx, p.ID, imageUrl)
	// if err != nil {
	// 	return err
	// }

	q.deleteCacheQuestion(ctx, p.ID.String())
	return nil
}

func (h *QuestionHandler) deleteCacheQuestion(ctx context.Context, questionID string) {
	h.rdb.Del(ctx, "question:"+questionID)

	iter := h.rdb.Scan(ctx, 0, "questions:*", 0).Iterator()
	for iter.Next(ctx) {
		h.rdb.Del(ctx, iter.Val())
	}
}
