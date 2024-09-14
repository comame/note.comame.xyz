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

func getConnection() (*connection, error) {
	// TODO: コネクションプールとかを使う
	db, err := sql.Open("mysql", "root:root@(mysql.comame.dev)/note")
	if err != nil {
		return nil, err
	}

	return &connection{db}, nil
}

func (c *connection) Close() {
	c.db.Close()
}

func (c *connection) findPostByURLKey(ctx context.Context, urlKey string) (*post, error) {
	rows, err := c.db.QueryContext(ctx, `
		SELECT id, url_key, created_datetime, updated_datetime, title, text, visibility
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
	if err := rows.Scan(&p.ID, &p.URLKey, &p.CreatedDatetime, &p.UpdatedDatetime, &p.Title, &p.Text, &p.Visibility); err != nil {
		return nil, err
	}

	return p, nil
}

func (c *connection) findPostByID(ctx context.Context, id uint64) (*post, error) {
	rows, err := c.db.QueryContext(ctx, `
		SELECT id, url_key, created_datetime, updated_datetime, title, text, visibility
		FROM nt_post
		WHERE id = ?
	`, id)
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, errNotFound
	}

	p := new(post)
	if err := rows.Scan(&p.ID, &p.URLKey, &p.CreatedDatetime, &p.UpdatedDatetime, &p.Title, &p.Text, &p.Visibility); err != nil {
		return nil, err
	}

	return p, nil
}

func (c *connection) createPost(ctx context.Context, post post) error {
	t, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer t.Rollback()

	if _, err := c.db.Exec(`
		INSERT INTO nt_post
		(url_key, created_datetime, updated_datetime, title, text, visibility)
		values
		(?, ?, ?, ?, ?, ?)
		`, post.URLKey, post.CreatedDatetime, post.UpdatedDatetime, post.Title, post.Text, post.Visibility); err != nil {
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
			nt_post.visibility
		FROM nt_post
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

func (c *connection) updatePost(ctx context.Context, post post) error {
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	r, err := c.db.ExecContext(ctx, `
		UPDATE nt_post
		SET
			updated_datetime = ?,
			title = ?,
			text = ?,
			visibility = ?
		WHERE
			id = ?
	`, post.UpdatedDatetime, post.Title, post.Text, post.Visibility, post.ID)
	if err != nil {
		return err
	}

	a, err := r.RowsAffected()
	if err != nil {
		return err
	}

	if a != 1 {
		return errNotFound
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (c *connection) deletePost(ctx context.Context, postID uint64) error {
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	r, err := tx.ExecContext(ctx, `
		DELETE FROM nt_post
		WHERE id=?
	`, postID)
	if err != nil {
		return err
	}

	a, err := r.RowsAffected()
	if err != nil {
		return err
	}

	if a != 1 {
		return errNotFound
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
