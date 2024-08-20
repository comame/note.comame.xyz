package server

import (
	"context"
	"database/sql"
	"errors"
)

type connection struct {
	db *sql.DB
}

var errNotFound = errors.New("not found")
var errNoRows = errors.New("no rows")

func getConnection() (*connection, error) {
	db, err := sql.Open("mysql", "root:root@(mysql.comame.dev)/note")
	if err != nil {
		return nil, err
	}

	return &connection{db}, nil
}

func (c *connection) Close() {
	c.db.Close()
}

func (c *connection) findPost(ctx context.Context, urlKey string) (*post, error) {
	rows, err := c.db.QueryContext(ctx, `
		SELECT id, url_key, created_datetime, updated_datetime, title, text
		FROM nt_post
		WHERE url_key = ?
	`, urlKey)

	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, errNotFound
	}

	p := new(post)
	if err := rows.Scan(&p.ID, &p.URLKey, &p.CreatedDatetime, &p.UpdatedDatetime, &p.Title, &p.Text); err != nil {
		return nil, err
	}

	return p, nil
}

func (c *connection) findVisibility(ctx context.Context, post post) (*post, error) {
	// 別の記事の公開状態と取り違えないように、ゼロ値の場合はエラーとする
	if post.ID == 0 {
		return nil, errors.New("id is zero")
	}

	rows, err := c.db.QueryContext(ctx, `
		SELECT visibility
		FROM nt_post_visibility
		WHERE post_id = ?
	`, post.ID)
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, errNoRows
	}

	if err := rows.Scan(&post.Visibility); err != nil {
		return nil, err
	}

	return &post, nil
}
