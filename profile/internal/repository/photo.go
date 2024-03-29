package repository

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"go.uber.org/multierr"

	"github.com/the-swiply/swiply-backend/profile/internal/dbmodel"
)

const (
	photoDirectoryName = "%v/photos"
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

func (p *PhotoRepository) List(ctx context.Context, userID uuid.UUID) ([]dbmodel.Photo, error) {
	var (
		keys    []string
		photos  []dbmodel.Photo
		listErr error
	)

	for obj := range p.s3.ListObjects(ctx, p.bucketName, minio.ListObjectsOptions{
		Prefix:    fmt.Sprintf(photoDirectoryName, userID),
		Recursive: true,
	}) {
		listErr = multierr.Combine(listErr, obj.Err)
		keys = append(keys, obj.Key)
	}

	if listErr != nil {
		return nil, listErr
	}

	for _, key := range keys {
		photoID, err := uuid.Parse(key)
		if err != nil {
			return nil, err
		}

		photo, err := p.Get(ctx, userID, photoID)
		photos = append(photos, photo)
	}

	return photos, nil
}

func (p *PhotoRepository) Get(ctx context.Context, userID, photoID uuid.UUID) (dbmodel.Photo, error) {
	obj, err := p.s3.GetObject(ctx, p.bucketName, fmt.Sprintf(photoDirectoryName+"/%v", userID, photoID), minio.GetObjectOptions{})
	if err != nil {
		return dbmodel.Photo{}, err
	}
	defer func() { _ = obj.Close() }()

	content, err := io.ReadAll(obj)
	if err != nil {
		return dbmodel.Photo{}, err
	}

	return dbmodel.Photo{
		ID:      photoID,
		Content: content,
	}, nil
}

func (p *PhotoRepository) Create(ctx context.Context, userID uuid.UUID, photo dbmodel.Photo) error {
	buf := bytes.NewReader(photo.Content)

	_, err := p.s3.PutObject(ctx, p.bucketName, fmt.Sprintf(photoDirectoryName+"/%v", userID, photo.ID), buf, -1, minio.PutObjectOptions{})

	return err
}

func (p *PhotoRepository) Delete(ctx context.Context, userID, photoID uuid.UUID) error {
	return p.s3.RemoveObject(ctx, p.bucketName, fmt.Sprintf(photoDirectoryName+"/%v", userID, photoID), minio.RemoveObjectOptions{})
}
