package repository

import "github.com/Masterminds/squirrel"

const (
	updateInfoTable  = "update_info"
	statisticsTable  = "statistic"
	profileTable     = "profile"
	interactionTable = "interaction"
	recHistoryTable  = "recommendation_history"
)

var (
	sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
)
