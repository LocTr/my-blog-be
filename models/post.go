package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Post struct {
	ID        int       `bun:"id,pk,autoincrement" json:"id"`
	Title     string    `bun:"title,notnull" json:"title"`
	Content   string    `bun:"content,notnull" json:"content"`
	CreatedAt time.Time `bun:"created_at,notnull" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,notnull" json:"updated_at"`
}

// BeforeInsert hook
func (p *Post) BeforeInsert(db *bun.DB) error {
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
	return nil
}
