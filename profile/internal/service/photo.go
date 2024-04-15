package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/the-swiply/swiply-backend/profile/internal/dbmodel"
)

type PhotoRepository interface {
	List(ctx context.Context, userID uuid.UUID) ([]dbmodel.Photo, error)
	Get(ctx context.Context, userID, photoID uuid.UUID) (dbmodel.Photo, error)
	Create(ctx context.Context, userID uuid.UUID, photo dbmodel.Photo) error
	Delete(ctx context.Context, userID, photoID uuid.UUID) error
}

type PhotoService struct {
	cfg  PhotoConfig
	repo PhotoRepository
}

func NewPhotoService(cfg PhotoConfig, repo PhotoRepository) *PhotoService {
	return &PhotoService{
		cfg:  cfg,
		repo: repo,
	}
}
