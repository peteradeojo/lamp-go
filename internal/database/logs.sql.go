// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: logs.sql

package database

import (
	"context"
	"database/sql"
)

const getAppWithToken = `-- name: GetAppWithToken :one
SELECT id, token, userid FROM apps WHERE token = ?
`

func (q *Queries) GetAppWithToken(ctx context.Context, token sql.NullString) (App, error) {
	row := q.db.QueryRowContext(ctx, getAppWithToken, token)
	var i App
	err := row.Scan(&i.ID, &i.Token, &i.Userid)
	return i, err
}

const saveLogs = `-- name: SaveLogs :execresult
INSERT INTO logs (appId, text, createdAt, updatedAt, level, saved) VALUES (?, ?, NOW(), NOW(), ?, 0)
`

type SaveLogsParams struct {
	Appid int64
	Text  string
	Level string
}

func (q *Queries) SaveLogs(ctx context.Context, arg SaveLogsParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, saveLogs, arg.Appid, arg.Text, arg.Level)
}