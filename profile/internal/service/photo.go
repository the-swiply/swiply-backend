package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/the-swiply/swiply-backend/profile/internal/converter"
	"github.com/the-swiply/swiply-backend/profile/internal/dbmodel"
	"github.com/the-swiply/swiply-backend/profile/internal/domain"
)

type PhotoContentRepository interface {
	Get(ctx context.Context, userID, photoID uuid.UUID) (dbmodel.Photo, error)
	Create(ctx context.Context, userID uuid.UUID, photo dbmodel.Photo) error
	Delete(ctx context.Context, userID, photoID uuid.UUID) error
}

type PhotoRepository interface {
	Reorder(ctx context.Context, userID uuid.UUID, photoIDs []uuid.UUID) error
	List(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
	Append(ctx context.Context, userID, photoID uuid.UUID) error
	Delete(ctx context.Context, userID, photoID uuid.UUID) error
	Create(ctx context.Context, userID uuid.UUID, photoIDs []uuid.UUID) error
}

type PhotoService struct {
	cfg         PhotoConfig
	contentRepo PhotoContentRepository
	repo        PhotoRepository
}

func NewPhotoService(cfg PhotoConfig, contentRepo PhotoContentRepository, repo PhotoRepository) *PhotoService {
	return &PhotoService{
		cfg:         cfg,
		contentRepo: contentRepo,
		repo:        repo,
	}
}

func (p *PhotoService) Create(ctx context.Context, userID uuid.UUID, photo domain.Photo) (uuid.UUID, error) {
	ph := converter.PhotoFromDomainToDBModel(photo)
	if err := p.contentRepo.Create(ctx, userID, ph); err != nil {
		return uuid.Nil, err
	}

	if _, err := p.repo.List(ctx, userID); err == domain.ErrEntityIsNotExists {
		return ph.ID, p.repo.Create(ctx, userID, []uuid.UUID{ph.ID})
	} else if err == nil {
		return ph.ID, p.repo.Append(ctx, userID, ph.ID)
	} else {
		return uuid.Nil, err
	}
}

func (p *PhotoService) Get(ctx context.Context, userID, photoID uuid.UUID) (domain.Photo, error) {
	photo, err := p.contentRepo.Get(ctx, userID, photoID)
	if err != nil {
		return domain.Photo{}, err
	}

	return converter.PhotoFromDBModelToDomain(photo), nil
}

func (p *PhotoService) List(ctx context.Context, userID uuid.UUID) ([]domain.Photo, error) {
	photoIDs, err := p.repo.List(ctx, userID)
	if err != nil {
		return nil, err
	}

	var photos []domain.Photo
	for _, id := range photoIDs {
		photo, err := p.contentRepo.Get(ctx, userID, id)
		if err != nil {
			return nil, err
		}
		photos = append(photos, converter.PhotoFromDBModelToDomain(photo))
	}

	return photos, nil
}

func (p *PhotoService) Delete(ctx context.Context, userID, photoID uuid.UUID) (uuid.UUID, error) {
	if err := p.repo.Delete(ctx, userID, photoID); err != nil {
		return uuid.Nil, err
	}

	return photoID, p.contentRepo.Delete(ctx, userID, photoID)
}

func (p *PhotoService) Reorder(ctx context.Context, userID uuid.UUID, photoIDs []uuid.UUID) error {
	return p.repo.Reorder(ctx, userID, photoIDs)
}
