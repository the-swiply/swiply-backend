package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"

	"github.com/the-swiply/swiply-backend/profile/internal/dbmodel"
)

const (
	photoDirectoryName = "%s/%s/photos"
)

type PhotoRepository struct {
	bucketName string
	s3         *minio.Client
}

func NewPhotoRepository(bucketName string, s3 *minio.Client) *PhotoRepository {
	return &PhotoRepository{
		bucketName: bucketName,
		s3:         s3,
	}
}

func (p *PhotoRepository) List(ctx context.Context, photoID uuid.UUID) ([]dbmodel.Photo, error) {
	p.s3.ListObjects(ctx, fmt.Sprintf(photoDirectoryName, p.bucketName, photoID.String()), minio.ListObjectsOptions{})
}

func (p *PhotoRepository) Get(ctx context.Context, userID, photoID uuid.UUID) (dbmodel.Photo, error) {

}

func (p *PhotoRepository) Delete(ctx context.Context, userID, photoID uuid.UUID) error {

}
