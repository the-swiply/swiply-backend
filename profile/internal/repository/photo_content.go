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

type PhotoContentRepository struct {
	bucketName string
	s3         *minio.Client
}

func NewPhotoContentRepository(bucketName string, s3 *minio.Client) *PhotoContentRepository {
	return &PhotoContentRepository{
		bucketName: bucketName,
		s3:         s3,
	}
}

func (p *PhotoContentRepository) List(ctx context.Context, userID uuid.UUID) ([]dbmodel.Photo, error) {
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

func (p *PhotoContentRepository) Get(ctx context.Context, userID, photoID uuid.UUID) (dbmodel.Photo, error) {
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

func (p *PhotoContentRepository) Create(ctx context.Context, userID uuid.UUID, photo dbmodel.Photo) error {
	buf := bytes.NewReader(photo.Content)

	_, err := p.s3.PutObject(ctx, p.bucketName, fmt.Sprintf(photoDirectoryName+"/%v", userID, photo.ID), buf, -1, minio.PutObjectOptions{})

	return err
}

func (p *PhotoContentRepository) Delete(ctx context.Context, userID, photoID uuid.UUID) error {
	return p.s3.RemoveObject(ctx, p.bucketName, fmt.Sprintf(photoDirectoryName+"/%v", userID, photoID), minio.RemoveObjectOptions{})
}
