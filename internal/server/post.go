package server

import (
	"context"
	"errors"
	"fmt"

	"github.com/comame/note.comame.xyz/internal/md"
)

type post struct {
	ID                  uint64 `json:"id"`
	URLKey              string `json:"url_key"`
	CreatedDatetime     string `json:"createdDatetime"`
	UpdatedDatetime     string `json:"updatedDatetime"`
	Title               string `json:"title"`
	Text                string `json:"text"`
	HTML                string `json:"html"`
	PermissionInherited bool   `json:"permissionInherited"`
	Parent              uint64 `json:"parent"` // 最上位記事なら0

	// フロントエンド用の値であり、getPermissionするとセットされる。サーバー側ではこのフィールドを参照せずに、post.getPermission() を呼び出すこと。
	ResolvedPermission permission `json:"resolvedPermission"`
	// この記事に設定された権限。権限を取得するには post.getPermission() を呼び出すこと。
	Permission permission `json:"permission"`

	// 階層構造の取得は重たいので、そのキャッシュ用の内部的なフィールド。
	// 詳細については fetchHierarchy を参照。
	hierarchy        []post
	hierarchyFetched bool
}

type permission string

const (
	permissionPublic  permission = "public"
	permissionPrivate permission = "private"
	permissionURL     permission = "url"
)

var (
	// post.ID = 0 のとき、意図せずゼロ値が入ってしまっている可能性が高いのでエラーとする
	errIDIsZero     = errors.New("id is zero")
	errNoPermission = errors.New("no permission")
)

var maxAllowedHierarchyDepth = 15

// 階層構造を取得し、p.hierarchy に格納する
// p.hierarchy は先頭 (index:0) に p 自身が入る。末尾に最上位記事が入る。
// p.hierarchy には階層構造に関する一部フィールドのみ入る (getPartialPostHierarchyRelatedInfoを参照)
func (p *post) fetchHierarchy(ctx context.Context) error {
	if p.hierarchyFetched {
		return nil
	}
	defer func() {
		p.hierarchyFetched = true
	}()

	con, err := GetConnection()
	if err != nil {
		return err
	}

	curr := *p
	for range maxAllowedHierarchyDepth {
		p.hierarchy = append(p.hierarchy, curr)

		if curr.Parent == 0 {
			return nil
		}

		parent, err := con.getPartialPostHierarchyRelatedInfo(ctx, curr.Parent)
		if err != nil {
			return err
		}

		curr = *parent
	}

	// 無限ループを防ぐため、maxAllowedHierarchyDepth を超えた階層の場合はエラーにする
	return errors.New("記事の階層構造が深すぎる")
}

// 親記事をたどって権限を取得
func (p *post) getPermission(ctx context.Context) (permission, error) {
	if err := p.fetchHierarchy(ctx); err != nil {
		return "", err
	}
	for _, h := range p.hierarchy {
		if !h.PermissionInherited || h.Parent == 0 {
			p.ResolvedPermission = h.Permission
			return h.Permission, nil
		}
	}
	return "", errors.New("階層構造がおかしい")
}

func (p *post) getURL() string {
	return fmt.Sprintf("/posts/%s", p.URLKey)
}

func (p *post) isAllowedToView(ctx context.Context, isLoggedIn bool) (bool, error) {
	permission, err := p.getPermission(ctx)
	if err != nil {
		return false, err
	}

	// ログインしていたら全部見れる
	if isLoggedIn {
		return true, nil
	}

	if permission == permissionPublic {
		return true, nil
	}
	if permission == permissionURL {
		return true, nil
	}

	return false, nil
}

// 閲覧権限のある記事を取得
func getPostByID(ctx context.Context, id uint64, isLoggedIn bool) (*post, error) {
	c, err := GetConnection()
	if err != nil {
		return nil, err
	}

	p, err := c.findPostByIDWithContent(ctx, id)
	if errors.Is(err, errNotFound) {
		return nil, errNotFound
	}
	if err != nil {
		return nil, err
	}

	ok, err := p.isAllowedToView(ctx, isLoggedIn)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errNoPermission
	}

	p.HTML = md.ToHTML(p.Text)

	return p, nil
}

// 閲覧権限のある記事を取得
func getAllowedPostByURLKey(ctx context.Context, urlKey string, isLoggedIn bool) (*post, error) {
	c, err := GetConnection()
	if err != nil {
		return nil, err
	}

	p, err := c.findPostByURLKeyWithContent(ctx, urlKey)
	if errors.Is(err, errNotFound) {
		return nil, errNotFound
	}
	if err != nil {
		return nil, err
	}

	ok, err := p.isAllowedToView(ctx, isLoggedIn)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errNoPermission
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

func listPost(ctx context.Context, isLoggedIn bool) ([]post, error) {
	con, err := GetConnection()
	if err != nil {
		return nil, err
	}

	var posts []post
	if isLoggedIn {
		posts, err = con.getAllPostsForAdmin(ctx)
		if err != nil {
			return nil, err
		}
	} else {
		posts, err = con.getAllPostsForAnonymous(ctx)
		if err != nil {
			return nil, err
		}
	}

	for i := range posts {
		// .ResolvedPermission フィールドをセットするために呼び出しておく
		// FIXME: この処理はまあまあ重たいので何とかしたほうがよさそう
		posts[i].getPermission(ctx)
	}

	return posts, nil
}
