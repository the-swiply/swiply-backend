package repository

import (
	"bytes"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/the-swiply/swiply-backend/event/internal/domain"
	"go.uber.org/multierr"
	"io"
	"strconv"
	"strings"
)

const (
	eventsPhotoDirectoryPattern = "event/%d/photos"
)

type PhotoContentRepository struct {
	bucketName string
	s3         *minio.Client
}

func NewPhotoContentRepository(ctx context.Context, bucketName string, s3 *minio.Client, createBucketIfNotExsits bool) (*PhotoContentRepository, error) {
	ok, err := s3.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, fmt.Errorf("can't check if s3 bucket exists: %w", err)
	}

	if !ok {
		if !createBucketIfNotExsits {
			return nil, fmt.Errorf("bucket %s does not exist", bucketName)
		}

		err = s3.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("can't create bucket %s", bucketName)
		}
	}

	return &PhotoContentRepository{
		bucketName: bucketName,
		s3:         s3,
	}, nil
}

func (p *PhotoContentRepository) ListPhotos(ctx context.Context, eventID int64) ([]domain.Photo, error) {
	var (
		keys    []string
		listErr error
	)

	for obj := range p.s3.ListObjects(ctx, p.bucketName, minio.ListObjectsOptions{
		Prefix:    fmt.Sprintf(eventsPhotoDirectoryPattern, eventID),
		Recursive: true,
	}) {
		listErr = multierr.Combine(listErr, obj.Err)
		keys = append(keys, obj.Key)
	}

	if listErr != nil {
		return nil, listErr
	}

	photos := make([]domain.Photo, 0, len(keys))
	for _, key := range keys {
		keySplitted := strings.Split(key, "/")
		if len(keySplitted) < 1 {
			return nil, fmt.Errorf("invalid key: %s", key)
		}

		photoID, err := strconv.Atoi(keySplitted[len(keySplitted)-1])
		if err != nil {
			return nil, err
		}

		photo, err := p.GetPhoto(ctx, eventID, int64(photoID))
		photos = append(photos, photo)
	}

	return photos, nil
}

func (p *PhotoContentRepository) GetPhoto(ctx context.Context, eventID int64, photoID int64) (domain.Photo, error) {
	obj, err := p.s3.GetObject(ctx, p.bucketName, fmt.Sprintf(eventsPhotoDirectoryPattern+"/%d", eventID, photoID), minio.GetObjectOptions{})
	if err != nil {
		return domain.Photo{}, err
	}
	defer obj.Close()

	content, err := io.ReadAll(obj)
	if err != nil {
		return domain.Photo{}, err
	}

	return domain.Photo{
		ID:      photoID,
		Content: content,
	}, nil
}

func (p *PhotoContentRepository) UploadPhoto(ctx context.Context, eventID int64, photo domain.Photo) error {
	buf := bytes.NewReader(photo.Content)

	_, err := p.s3.PutObject(ctx, p.bucketName, fmt.Sprintf(eventsPhotoDirectoryPattern+"/%d", eventID, photo.ID), buf, -1, minio.PutObjectOptions{})

	return err
}
