package author

import (
	"context"
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/HemanthGowda/library-management-go/internal/domain"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Add(ctx context.Context, request CreateAuthorRequest) (AuthorResponse, error) {
	name := strings.TrimSpace(request.Name)
	biography := strings.TrimSpace(request.Biography)
	nationality := strings.TrimSpace(request.Nationality)

	if violations := validateAuthor(name, biography, request.Birth, nationality); len(violations) > 0 {
		return AuthorResponse{}, domain.ValidationFailed(violations)
	}

	exists, err := s.repository.ExistByName(ctx, name)
	if err != nil {
		return AuthorResponse{}, err
	}

	if exists {
		return AuthorResponse{}, domain.Conflict("Author already exists")
	}

	created, err := s.repository.Create(ctx, Author{
		Name:        name,
		Biography:   biography,
		Birth:       request.Birth,
		Nationality: nationality,
	})

	if errors.Is(err, ErrDuplicate) {
		return AuthorResponse{}, domain.Conflict("Author already exists")
	}

	if err != nil {
		return AuthorResponse{}, err
	}

	return toResponse(created), nil
}

func validateAuthor(name string, biography string, birth int, nationality string) []domain.Violation {
	var violations []domain.Violation
	switch {
	case name == "":
		violations = append(violations, domain.Violation{Field: "name", Message: "Author name is required"})
	case utf8.RuneCountInString(name) > 160:
		violations = append(violations, domain.Violation{Field: "name", Message: "Author name must be at most 160 characters"})
	}

	if utf8.RuneCountInString(biography) > 4000 {
		violations = append(violations, domain.Violation{Field: "biography", Message: "Biography must be at most 4000 characters"})
	}

	switch {
	case birth < 1:
		violations = append(violations, domain.Violation{Field: "birth", Message: "Birth year must be greater than 0"})
	case birth > time.Now().Year():
		violations = append(violations, domain.Violation{Field: "birth", Message: "Birth year cannot be in the future"})
	}

	if utf8.RuneCountInString(nationality) > 80 {
		violations = append(violations, domain.Violation{Field: "nationality", Message: "Nationality must be at most 80 characters"})
	}

	return violations
}

func toResponse(author Author) AuthorResponse {
	return AuthorResponse{
		Name:        author.Name,
		Biography:   author.Biography,
		Birth:       author.Birth,
		Nationality: author.Nationality,
	}
}
