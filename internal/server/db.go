package server

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"sync"
)

type connection struct {
	db  *sql.DB
	tx  *sql.Tx
	txm sync.Mutex
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
	c.txm.Lock()
	defer c.txm.Unlock()

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
	c.txm.Lock()
	defer c.txm.Unlock()

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
	c.txm.Lock()
	defer c.txm.Unlock()

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
	c.txm.Lock()
	defer c.txm.Unlock()

	if c.tx == nil {
		return errNotInTransaction
	}

	return nil
}

func (c *connection) findPostByURLKey(ctx context.Context, urlKey string) (*post, error) {
	rows, err := c.db.QueryContext(ctx, `
		SELECT id, url_key, created_datetime, updated_datetime, title, text, permission
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
	if err := rows.Scan(&p.ID, &p.URLKey, &p.CreatedDatetime, &p.UpdatedDatetime, &p.Title, &p.Text, &p.Permission); err != nil {
		return nil, err
	}

	return p, nil
}

func (c *connection) findPostByID(ctx context.Context, id uint64) (*post, error) {
	rows, err := c.db.QueryContext(ctx, `
		SELECT id, url_key, created_datetime, updated_datetime, title, text, permission
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
	if err := rows.Scan(&p.ID, &p.URLKey, &p.CreatedDatetime, &p.UpdatedDatetime, &p.Title, &p.Text, &p.Permission); err != nil {
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
		(url_key, created_datetime, updated_datetime, title, text, permission)
		values
		(?, ?, ?, ?, ?, ?)
		`, post.URLKey, post.CreatedDatetime, post.UpdatedDatetime, post.Title, post.Text, post.Permission); err != nil {
		return err
	}

	if err := t.Commit(); err != nil {
		return err
	}

	return nil
}

func (c *connection) getAllPostsForAnonymous(ctx context.Context) ([]post, error) {
	rows, err := c.db.QueryContext(ctx, `
		SELECT
			nt_post.id,
			nt_post.url_key,
			nt_post.created_datetime,
			nt_post.updated_datetime,
			nt_post.title,
			nt_post.text,
			nt_post.permission
		FROM nt_post
		WHERE permission = 'public'
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
			&post.Permission,
		); err != nil {
			return nil, err
		}
		p = append(p, post)
	}

	return p, nil
}

func (c *connection) getAllPostsForAdmin(ctx context.Context) ([]post, error) {
	rows, err := c.db.QueryContext(ctx, `
		SELECT
			nt_post.id,
			nt_post.url_key,
			nt_post.created_datetime,
			nt_post.updated_datetime,
			nt_post.title,
			nt_post.text,
			nt_post.permission
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
			&post.Permission,
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
			permission = ?
		WHERE
			id = ?
	`, post.UpdatedDatetime, post.Title, post.Text, post.Permission, post.ID)
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
			permission
		)
		SELECT
			id,
			url_key,
			created_datetime,
			updated_datetime,
			text,
			permission
		FROM nt_post
		WHERE nt_post.id = ?
	`, postID); err != nil {
		return err
	}

	return nil
}
