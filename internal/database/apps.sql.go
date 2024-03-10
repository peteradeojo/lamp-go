// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: apps.sql

package database

import (
	"context"
)

const getApps = `-- name: GetApps :many
SELECT id, token, userid FROM apps LIMIT $1
`

func (q *Queries) GetApps(ctx context.Context, limit int32) ([]App, error) {
	rows, err := q.db.QueryContext(ctx, getApps, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []App
	for rows.Next() {
		var i App
		if err := rows.Scan(&i.ID, &i.Token, &i.Userid); err != nil {
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
