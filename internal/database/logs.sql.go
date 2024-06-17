// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: logs.sql

package database

import (
	"context"
	"database/sql"

	"github.com/sqlc-dev/pqtype"
)

const createSystemLog = `-- name: CreateSystemLog :exec
INSERT INTO system_logs (id, text, stack, context, level, from_system, createdat, updatedat, origin) VALUES (
  uuid_generate_v4(), 
  $1, 
  $2, 
  $3, 
  $4, 
  B'1', 
  now(), 
  now(), 
  'go-api'
)
`

type CreateSystemLogParams struct {
	Text    string
	Stack   sql.NullString
	Context pqtype.NullRawMessage
	Level   LogLevel
}

func (q *Queries) CreateSystemLog(ctx context.Context, arg CreateSystemLogParams) error {
	_, err := q.db.ExecContext(ctx, createSystemLog,
		arg.Text,
		arg.Stack,
		arg.Context,
		arg.Level,
	)
	return err
}

const exportLogs = `-- name: ExportLogs :many
SELECT id, text, apptoken, level, createdat, updatedat, context, ip, tags FROM logs WHERE apptoken = $1 ORDER BY createdat DESC LIMIT 100
`

func (q *Queries) ExportLogs(ctx context.Context, apptoken string) ([]Log, error) {
	rows, err := q.db.QueryContext(ctx, exportLogs, apptoken)
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

const getAppWithToken = `-- name: GetAppWithToken :one
SELECT id, token, userid FROM apps WHERE token = $1
`

func (q *Queries) GetAppWithToken(ctx context.Context, token sql.NullString) (App, error) {
	row := q.db.QueryRowContext(ctx, getAppWithToken, token)
	var i App
	err := row.Scan(&i.ID, &i.Token, &i.Userid)
	return i, err
}

const getLogs = `-- name: GetLogs :many
SELECT id, text, apptoken, level, createdat, updatedat, context, ip, tags FROM logs WHERE appToken = $1
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

const saveLogs = `-- name: SaveLogs :one
INSERT INTO logs (apptoken, text, createdat, updatedat, level, context,ip, tags) VALUES ($1, $2, now(), now(), $3, $4, $5, $6) RETURNING 1
`

type SaveLogsParams struct {
	Apptoken string
	Text     string
	Level    string
	Context  pqtype.NullRawMessage
	Ip       sql.NullString
	Tags     pqtype.NullRawMessage
}

func (q *Queries) SaveLogs(ctx context.Context, arg SaveLogsParams) (int32, error) {
	row := q.db.QueryRowContext(ctx, saveLogs,
		arg.Apptoken,
		arg.Text,
		arg.Level,
		arg.Context,
		arg.Ip,
		arg.Tags,
	)
	var column_1 int32
	err := row.Scan(&column_1)
	return column_1, err
}
