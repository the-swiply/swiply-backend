package converter

import (
	"time"

	"github.com/google/uuid"

	"github.com/the-swiply/swiply-backend/profile/internal/dbmodel"
	"github.com/the-swiply/swiply-backend/profile/internal/domain"
	"github.com/the-swiply/swiply-backend/profile/pkg/api/profile"
)

func InteractionFromDBModelToDomain(interaction dbmodel.Interaction) (domain.Interaction, error) {
	inter := domain.Interaction{
		ID:   interaction.ID,
		From: interaction.From,
		To:   interaction.To,
	}

	if err := inter.Type.Set(interaction.Type); err != nil {
		return inter, err
	}

	return inter, nil
}

func InteractionFromDomainToDBModel(interaction domain.Interaction) dbmodel.Interaction {
	return dbmodel.Interaction{
		ID:        interaction.ID,
		From:      interaction.From,
		To:        interaction.To,
		Type:      string(interaction.Type),
		CreatedAt: time.Now().UTC(),
	}
}

func InteractionFromProtoToDomain(interaction *profile.Interaction) domain.Interaction {
	inter := domain.Interaction{
		From: uuid.MustParse(interaction.GetFrom()),
		To:   uuid.MustParse(interaction.GetTo()),
	}

	switch interaction.GetType() {
	case profile.InteractionType_INTERACTION_TYPE_UNSPECIFIED:
		inter.Type = domain.InteractionTypeUnspecified
	case profile.InteractionType_LIKE:
		inter.Type = domain.InteractionTypeLike
	case profile.InteractionType_DISLIKE:
		inter.Type = domain.InteractionTypeDislike
	}

	return inter
}

func InteractionFromDomainToProto(interaction domain.Interaction) *profile.Interaction {
	inter := &profile.Interaction{
		From: interaction.From.String(),
		To:   interaction.To.String(),
	}

	switch interaction.Type {
	case domain.InteractionTypeUnspecified:
		inter.Type = profile.InteractionType_INTERACTION_TYPE_UNSPECIFIED
	case domain.InteractionTypeLike:
		inter.Type = profile.InteractionType_LIKE
	case domain.InteractionTypeDislike:
		inter.Type = profile.InteractionType_DISLIKE
	}

	return inter
}
