// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: bearertokens.sql

package dbLayer

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createBearerToken = `-- name: CreateBearerToken :exec
INSERT INTO BearerTokens (
    tokenString, validTill, userName
) VALUES (
    $1, $2, $3
)
`

type CreateBearerTokenParams struct {
	Tokenstring string
	Validtill   pgtype.Timestamp
	Username    string
}

func (q *Queries) CreateBearerToken(ctx context.Context, arg CreateBearerTokenParams) error {
	_, err := q.db.Exec(ctx, createBearerToken, arg.Tokenstring, arg.Validtill, arg.Username)
	return err
}

const deleteBearerToken = `-- name: DeleteBearerToken :exec
UPDATE BearerTokens
SET valid = False
WHERE userName = $1 AND valid = True
`

func (q *Queries) DeleteBearerToken(ctx context.Context, username string) error {
	_, err := q.db.Exec(ctx, deleteBearerToken, username)
	return err
}

const retrieveBearerToken = `-- name: RetrieveBearerToken :one
SELECT tokenstring, validtill, username, valid FROM BearerTokens
WHERE tokenString = $1 AND valid = True
`

func (q *Queries) RetrieveBearerToken(ctx context.Context, tokenstring string) (Bearertoken, error) {
	row := q.db.QueryRow(ctx, retrieveBearerToken, tokenstring)
	var i Bearertoken
	err := row.Scan(
		&i.Tokenstring,
		&i.Validtill,
		&i.Username,
		&i.Valid,
	)
	return i, err
}

const updateBearerTokenExpiration = `-- name: UpdateBearerTokenExpiration :exec
UPDATE BearerTokens
SET validTill = $2
WHERE tokenString = $1 AND valid = True
`

type UpdateBearerTokenExpirationParams struct {
	Tokenstring string
	Validtill   pgtype.Timestamp
}

func (q *Queries) UpdateBearerTokenExpiration(ctx context.Context, arg UpdateBearerTokenExpirationParams) error {
	_, err := q.db.Exec(ctx, updateBearerTokenExpiration, arg.Tokenstring, arg.Validtill)
	return err
}
