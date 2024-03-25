package converter

import (
	"github.com/google/uuid"
)

func UUIDsToStrings(recs []uuid.UUID) []string {
	res := make([]string, 0, len(recs))

	for _, rec := range recs {
		res = append(res, rec.String())
	}

	return res
}
