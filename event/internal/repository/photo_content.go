package repository

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
)

const (
	photoDirectoryName = "%v/photos"
)

type PhotoContentRepository struct {
	bucketName string
	s3         *minio.Client
}

func NewPhotoContentRepository(ctx context.Context, bucketName string, s3 *minio.Client) (*PhotoContentRepository, error) {
	ok, err := s3.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, fmt.Errorf("can't check if s3 bucket exists: %w", err)
	}

	if !ok {
		return nil, fmt.Errorf("bucket %s does not exist", bucketName)
	}

	return &PhotoContentRepository{
		bucketName: bucketName,
		s3:         s3,
	}, nil
}
