package converter

import (
	"time"

	"github.com/the-swiply/swiply-backend/profile/internal/dbmodel"
	"github.com/the-swiply/swiply-backend/profile/internal/domain"
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
