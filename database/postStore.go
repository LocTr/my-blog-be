package database

import (
	"context"

	"github.com/LocTr/my-blog-be/models"
	"github.com/uptrace/bun"
)

// PostStore implements database operations for posts.
type PostStore struct {
	db *bun.DB
}

// NewPostStore returns a PostStore.

func NewPostStore(db *bun.DB) *PostStore {
	return &PostStore{db: db}
}

// Get post by ID.
func (store *PostStore) GetPost(id int) (*models.Post, error) {
	post := &models.Post{ID: id}
	err := store.db.NewSelect().
		Model(post).
		WherePK().
		Scan(context.Background())

	return post, err
}

// GetPosts returns a list of posts with pagination.
func (store *PostStore) GetPosts(page, size int) ([]*models.Post, error) {
	var posts []*models.Post
	err := store.db.NewSelect().
		Model(&posts).
		Limit(size).
		Offset((page - 1) * size).
		Scan(context.Background())

	return posts, err
}
