package storage

import (
	"crypto/sha1"
	"fmt"
	"io"
	"remindbot/lib/e"
)

type Storage interface {
	Save(p *Page) error
	PickRandom(userName string) (*Page, error)
	PickAll(userName string) ([]Page, error)
	Remove(p *Page) error
	IfExists(p *Page) (bool, error)
}

type Page struct {
	URL       string
	UserName  string
	CreatedAt string
}

func (p *Page) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap("failed to hash", err)
	}

	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", e.Wrap("failed to hash", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
