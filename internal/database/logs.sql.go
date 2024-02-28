// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: logs.sql

package database

import (
	"context"
	"database/sql"
	"encoding/json"
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

const getLogs = `-- name: GetLogs :many
SELECT id, text, apptoken, level, createdat, updatedat, context, ip, tags FROM logs WHERE appToken = ?
`

func (q *Queries) GetLogs(ctx context.Context, apptoken string) ([]Log, error) {
	rows, err := q.db.QueryContext(ctx, getLogs, apptoken)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Log
	for rows.Next() {
		var i Log
		if err := rows.Scan(
			&i.ID,
			&i.Text,
			&i.Apptoken,
			&i.Level,
			&i.Createdat,
			&i.Updatedat,
			&i.Context,
			&i.Ip,
			&i.Tags,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const saveLogs = `-- name: SaveLogs :execresult
INSERT INTO logs (appToken, text, createdAt, updatedAt, level, saved, context,ip, tags) VALUES (?, ?, NOW(), NOW(), ?, 0, ?, ?, ?)
`

type SaveLogsParams struct {
	Apptoken string
	Text     string
	Level    string
	Context  json.RawMessage
	Ip       sql.NullString
	Tags     json.RawMessage
}

func (q *Queries) SaveLogs(ctx context.Context, arg SaveLogsParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, saveLogs,
		arg.Apptoken,
		arg.Text,
		arg.Level,
		arg.Context,
		arg.Ip,
		arg.Tags,
	)
}
