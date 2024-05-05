package converter

import (
	"github.com/google/uuid"
	"github.com/the-swiply/swiply-backend/pkg/houston/loggy"
	"github.com/the-swiply/swiply-backend/recommendation/internal/domain"
	"github.com/the-swiply/swiply-backend/recommendation/internal/pb/profile"
)

func ProfilesFromProtoToDomain(profiles []*profile.UserProfile) []domain.Profile {
	res := make([]domain.Profile, 0, len(profiles))
	for _, prof := range profiles {
		parsedID, parseErr := uuid.Parse(prof.GetId())
		if parseErr != nil {
			loggy.Errorf("failed to parse profile id = %s: %v", prof.GetId(), parseErr)
			continue
		}

		interests := make([]int64, 0, len(prof.GetInterests()))
		for _, interest := range prof.GetInterests() {
			interests = append(interests, interest.GetId())
		}

		res = append(res, domain.Profile{
			ID:               parsedID,
			Interests:        interests,
			BirthDay:         prof.GetBirthDay().AsTime(),
			Gender:           prof.GetGender().String(),
			Info:             prof.GetInfo(),
			SubscriptionType: prof.GetSubscriptionType().String(),
			LocationLat:      prof.GetLocation().GetLat(),
			LocationLon:      prof.GetLocation().GetLong(),
		})
	}

	return res
}

func InteractionsFromProtoToDomain(interactions []*profile.Interaction) []domain.Interaction {
	res := make([]domain.Interaction, 0, len(interactions))
	for _, interaction := range interactions {
		parsedFromID, parseErr := uuid.Parse(interaction.GetFrom())
		if parseErr != nil {
			loggy.Errorf("failed to parse profile id = %s: %v", interaction.GetFrom(), parseErr)
			continue
		}

		parsedToID, parseErr := uuid.Parse(interaction.GetTo())
		if parseErr != nil {
			loggy.Errorf("failed to parse profile id = %s: %v", interaction.GetTo(), parseErr)
			continue
		}

		var positive bool
		switch interaction.GetType() {
		case profile.InteractionType_LIKE:
			positive = true
		case profile.InteractionType_DISLIKE:
			positive = false
		default:
			loggy.Infoln("unknown interaction type: %s", interaction.GetType().String())
			continue
		}

		res = append(res, domain.Interaction{
			From:     parsedFromID,
			To:       parsedToID,
			Positive: positive,
		})
	}

	return res
}
