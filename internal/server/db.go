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

func (c *connection) createPost(ctx context.Context, post post) error {
	t, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer t.Rollback()

	r, err := c.db.Exec(`
		INSERT INTO nt_post
		(url_key, created_datetime, updated_datetime, title, text)
		values
		(?, ?, ?, ?, ?)
		`, post.URLKey, post.CreatedDatetime, post.UpdatedDatetime, post.Title, post.Text)
	if err != nil {
		return err
	}

	id, err := r.LastInsertId()
	if err != nil {
		return err
	}

	if _, err := c.db.ExecContext(ctx, `
		INSERT INTO nt_post_visibility
		(post_id, visibility)
		values
		(?, ?)
	`, id, post.Visibility); err != nil {
		return err
	}

	if err := t.Commit(); err != nil {
		return err
	}

	return nil
}

func (c *connection) getPosts(ctx context.Context) ([]post, error) {
	rows, err := c.db.QueryContext(ctx, `
		SELECT
			nt_post.id,
			nt_post.url_key,
			nt_post.created_datetime,
			nt_post.updated_datetime,
			nt_post.title,
			nt_post.text,
			nt_post_visibility.visibility
		FROM nt_post
		INNER JOIN nt_post_visibility
		ON nt_post.id = nt_post_visibility.post_id
	`)
	if err != nil {
		return nil, err
	}

	var p []post
	for rows.Next() {
		var post post
		if err := rows.Scan(
			&post.ID,
			&post.URLKey,
			&post.CreatedDatetime,
			&post.UpdatedDatetime,
			&post.Title,
			&post.Text,
			&post.Visibility,
		); err != nil {
			return nil, err
		}
		p = append(p, post)
	}

	return p, nil
}
