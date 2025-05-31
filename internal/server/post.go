package server

import (
	"context"
	"errors"
	"fmt"

	"github.com/comame/note.comame.xyz/internal/md"
)

type post struct {
	ID              uint64         `json:"id"`
	URLKey          string         `json:"url_key"`
	CreatedDatetime string         `json:"-"`
	UpdatedDatetime string         `json:"-"`
	Title           string         `json:"title"`
	Text            string         `json:"text"`
	Visibility      postVisibility `json:"visibility"`
	HTML            string         `json:"-"`
}

// FIXME: どう考えてもフロントエンドと統一したほうがいい
type postForFront struct {
	ID                  uint64     `json:"id"`
	URLKey              string     `json:"url_key"`
	CreatedDatetime     string     `json:"createdDatetime"`
	UpdatedDatetime     string     `json:"updatedDatetime"`
	Title               string     `json:"title"`
	Text                string     `json:"text"`
	HTML                string     `json:"html"`
	Permission          permission `json:"permission"`
	PermissionInherited bool       `json:"permissionInherited"`
}

type permission string

const (
	permissionPublic  permission = "public"
	permissionPrivate permission = "private"
	permissionURL     permission = "url"
)

func (p *post) toPostForFront() postForFront {
	permission := permissionPrivate
	switch p.Visibility {
	case postVisibilityPublic:
		permission = permissionPublic
	case postVisibilityPrivate:
		permission = permissionPrivate
	case postVisibilityUnlisted:
		permission = permissionURL
	}

	return postForFront{
		ID:                  p.ID,
		URLKey:              p.URLKey,
		CreatedDatetime:     p.CreatedDatetime,
		UpdatedDatetime:     p.UpdatedDatetime,
		Title:               p.Title,
		Text:                p.Text,
		HTML:                p.HTML,
		Permission:          permission,
		PermissionInherited: false,
	}
}

type postVisibility int

const (
	// 非公開
	postVisibilityPrivate postVisibility = 0
	// 限定公開
	postVisibilityUnlisted = 1
	// 全体公開
	postVisibilityPublic = 2
)

var (
	// post.ID = 0 のとき、意図せずゼロ値が入ってしまっている可能性が高いのでエラーとする
	errIDIsZero = errors.New("id is zero")
)

func (p *post) getURL() string {
	switch p.Visibility {
	case postVisibilityPublic:
		return fmt.Sprintf("/posts/public/%s", p.URLKey)
	case postVisibilityUnlisted:
		return fmt.Sprintf("/posts/unlisted/%s", p.URLKey)
	case postVisibilityPrivate:
		return fmt.Sprintf("/posts/private/%s", p.URLKey)
	}

	panic("unknown visibility")
}

func getPostByID(ctx context.Context, id uint64, visibility postVisibility) (*post, error) {
	c, err := GetConnection()
	if err != nil {
		return nil, err
	}

	p, err := c.findPostByID(ctx, id)
	if errors.Is(err, errNotFound) {
		return nil, errNotFound
	}
	if err != nil {
		return nil, err
	}

	if visibility != p.Visibility {
		return nil, errNotFound
	}

	p.HTML = md.ToHTML(p.Text)

	return p, nil
}

func getPostByURLKey(ctx context.Context, urlKey string, visibility postVisibility) (*post, error) {
	c, err := GetConnection()
	if err != nil {
		return nil, err
	}

	p, err := c.findPostByURLKey(ctx, urlKey)
	if errors.Is(err, errNotFound) {
		return nil, errNotFound
	}
	if err != nil {
		return nil, err
	}

	if visibility != p.Visibility {
		return nil, errNotFound
	}

	p.HTML = md.ToHTML(p.Text)

	return p, nil
}

func createPost(ctx context.Context, p post) (*post, error) {
	u, err := randomString(32)
	if err != nil {
		return nil, err
	}
	p.URLKey = u

	now := dateTimeNow()
	p.CreatedDatetime = now
	p.UpdatedDatetime = now

	con, err := GetConnection()
	if err != nil {
		return nil, err
	}

	if err := con.createPost(ctx, p); err != nil {
		return nil, err
	}

	return &p, nil
}

func updatePost(ctx context.Context, p post) error {
	if p.ID == 0 {
		return errIDIsZero
	}

	con, err := GetConnection()
	if err != nil {
		return err
	}

	if err := con.Begin(ctx); err != nil {
		return err
	}
	defer con.Rollback()

	if err := con.copyPostToPostLogInTransaction(ctx, p.ID); err != nil {
		return err
	}

	p.UpdatedDatetime = dateTimeNow()

	if err := con.updatePostInTransaction(ctx, p); err != nil {
		return err
	}

	if err := con.Commit(); err != nil {
		return err
	}

	return nil
}

func deletePost(ctx context.Context, postID uint64) error {
	if postID == 0 {
		return errIDIsZero
	}

	con, err := GetConnection()
	if err != nil {
		return err
	}

	if err := con.Begin(ctx); err != nil {
		return err
	}
	defer con.Rollback()

	if err := con.copyPostToPostLogInTransaction(ctx, postID); err != nil {
		return err
	}

	if err := con.deletePostInTransaction(ctx, postID); err != nil {
		return err
	}

	if err := con.Commit(); err != nil {
		return err
	}

	return nil
}
