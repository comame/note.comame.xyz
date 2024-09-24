package server

import (
	"context"
	"database/sql"
	"errors"
	"os"
)

type connection struct {
	db *sql.DB
	tx *sql.Tx
}

var (
	errNotFound             = errors.New("not found")
	errNotInTransaction     = errors.New("not in transaction")
	errAlreadyInTransaction = errors.New("already in transaction")
)

var dbInstance *sql.DB

func GetConnection() (*connection, error) {
	if dbInstance == nil {
		s := os.Getenv("MYSQL_CONNECT")
		db, err := sql.Open("mysql", s)
		if err != nil {
			return nil, err
		}
		dbInstance = db
	}

	return &connection{db: dbInstance, tx: nil}, nil
}

func (c *connection) Begin(ctx context.Context) error {
	if c.tx != nil {
		return errAlreadyInTransaction
	}

	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	c.tx = tx
	return nil
}

func (c *connection) Rollback() error {
	if c.tx == nil {
		return nil
	}

	if err := c.tx.Rollback(); err != nil {
		return err
	}
	c.tx = nil
	return nil
}

func (c *connection) Commit() error {
	if c.tx == nil {
		return nil
	}

	if err := c.tx.Commit(); err != nil {
		return err
	}
	c.tx = nil
	return nil
}

func (c *connection) transactionGuard() error {
	if c.tx == nil {
		return errNotInTransaction
	}

	return nil
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
	defer rows.Close()

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
	defer rows.Close()

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
	defer rows.Close()

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

func (c *connection) updatePostInTransaction(ctx context.Context, post post) error {
	if err := c.transactionGuard(); err != nil {
		return err
	}

	r, err := c.tx.ExecContext(ctx, `
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

	return nil
}

func (c *connection) deletePostInTransaction(ctx context.Context, postID uint64) error {
	if err := c.transactionGuard(); err != nil {
		return err
	}

	r, err := c.tx.ExecContext(ctx, `
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

	return nil
}

func (c *connection) copyPostToPostLogInTransaction(ctx context.Context, postID uint64) error {
	if err := c.transactionGuard(); err != nil {
		return err
	}

	if _, err := c.tx.ExecContext(ctx, `
		INSERT INTO nt_post_log (
			post_id,
			url_key,
			created_datetime,
			updated_datetime,
			text,
			visibility
		)
		SELECT
			id,
			url_key,
			created_datetime,
			updated_datetime,
			text,
			visibility
		FROM nt_post
		WHERE nt_post.id = ?
	`, postID); err != nil {
		return err
	}

	return nil
}
