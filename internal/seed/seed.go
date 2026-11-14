package seed

import (
	"context"
	"errors"
	"strings"

	"github.com/HemanthGowda/library-management-go/internal/book"
	"github.com/HemanthGowda/library-management-go/internal/config"
	"github.com/HemanthGowda/library-management-go/internal/password"
	"github.com/HemanthGowda/library-management-go/internal/user"
)

func Run(ctx context.Context, cfg config.Config, users user.Repository, books book.Repository, passwords password.Encoder) error {
	if err := seedAdmin(ctx, cfg, users, passwords); err != nil {
		return err
	}
	return seedBooks(ctx, books)
}

func seedAdmin(ctx context.Context, cfg config.Config, users user.Repository, passwords password.Encoder) error {
	email := strings.ToLower(strings.TrimSpace(cfg.SeedAdminEmail))
	exists, err := users.ExistsByEmail(ctx, email)
	if err != nil || exists {
		return err
	}

	hash, err := passwords.Encode(cfg.SeedAdminPassword)
	if err != nil {
		return err
	}

	_, err = users.Create(ctx, user.User{
		Email:        email,
		PasswordHash: hash,
		FullName:     strings.TrimSpace(cfg.SeedAdminFullName),
		Roles:        []user.Role{user.RoleAdmin},
		Enabled:      true,
	})
	if errors.Is(err, user.ErrDuplicate) {
		return nil
	}
	return err
}

func seedBooks(ctx context.Context, books book.Repository) error {
	samples := []book.Book{
		{Title: "Effective Java", AuthorName: "Joshua Bloch", ISBN: "9780134685991", Description: "Best practices for the Java platform.", PublicationYear: 2018, TotalCopies: 5, AvailableCopies: 5},
		{Title: "Head First Design Patterns", AuthorName: "Eric Freeman", ISBN: "9780596009205", Description: "A brain-friendly guide to design patterns.", PublicationYear: 2004, TotalCopies: 3, AvailableCopies: 3},
		{Title: "Designing Data-Intensive Applications", AuthorName: "Martin Kleppmann", ISBN: "9781492078005", Description: "The big ideas behind reliable, scalable data systems.", PublicationYear: 2017, TotalCopies: 4, AvailableCopies: 4},
	}

	for _, sample := range samples {
		exists, err := books.ExistsByISBN(ctx, sample.ISBN)
		if err != nil {
			return err
		}
		if exists {
			continue
		}
		if _, err := books.Insert(ctx, sample); err != nil && !errors.Is(err, book.ErrDuplicate) {
			return err
		}
	}
	return nil
}