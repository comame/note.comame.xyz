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

type postVisibility int

const (
	// 非公開
	postVisibilityPrivate postVisibility = 0
	// 限定公開
	postVisibilityUnlisted = 1
	// 全体公開
	postVisibilityPublic = 2
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

func (p *post) editURL() string {
	return fmt.Sprintf("/edit/post/%d", p.ID)
}

func (p *post) visibilityLabel() string {
	switch p.Visibility {
	case postVisibilityPublic:
		return "一般公開"
	case postVisibilityUnlisted:
		return "限定公開"
	case postVisibilityPrivate:
		return "非公開"
	}

	panic("unknown visibility")
}

func getPost(ctx context.Context, urlKey string, visibility postVisibility) (*post, error) {
	c, err := getConnection()
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

	con, err := getConnection()
	if err != nil {
		return nil, err
	}

	if err := con.createPost(ctx, p); err != nil {
		return nil, err
	}

	return &p, nil
}

func updatePost(ctx context.Context, p post) (*post, error) {
	con, err := getConnection()
	if err != nil {
		return nil, err
	}

	p.UpdatedDatetime = dateTimeNow()

	if err := con.updatePost(ctx, p); err != nil {
		return nil, err
	}

	return &p, nil
}
