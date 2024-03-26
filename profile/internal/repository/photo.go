package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"

	"github.com/the-swiply/swiply-backend/profile/internal/dbmodel"
)

const (
	photoDirectoryName = "photo"
)

type PhotoRepository struct {
	s3 *minio.Client
}

func NewPhotoRepository(s3 *minio.Client) *PhotoRepository {
	return &PhotoRepository{s3: s3}
}

func (p *PhotoRepository) List(ctx context.Context, photoID uuid.UUID) ([]dbmodel.Photo, error) {

}

func (p *PhotoRepository) Get(ctx context.Context, userID, photoID uuid.UUID) (dbmodel.Photo, error) {

}

func (p *PhotoRepository) Delete(ctx context.Context, userID, photoID uuid.UUID) error {

}
