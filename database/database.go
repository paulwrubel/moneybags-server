package database

//go:generate mockgen -source=$GOFILE -destination=../mocks/database/mock_$GOFILE -package=mockdatabase

import "github.com/jackc/pgtype/pgxtype"

type IHandler interface {
	pgxtype.Querier
}
