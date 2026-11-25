package author

import (
	"context"
	"strings"
	"testing"

	"github.com/HemanthGowda/library-management-go/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestAutorNameUnique(t *testing.T) {
	repository := newFakeAuthorRepo()
	service := NewService(repository)

	response, err := service.Add(context.Background(), CreateAuthorRequest{
		Name:        "Mitch Albom",
		Biography:   "I am author of book called Tuesdays with morrie",
		Birth:       1958,
		Nationality: "American",
	})
	require.NoError(t, err)
	require.Equal(t, "Mitch Albom", response.Name)
	require.Equal(t, "I am author of book called Tuesdays with morrie", response.Biography)
	require.Equal(t, 1958, response.Birth)
	require.Equal(t, "American", response.Nationality)
}

func TestAuthorNameIsRequired(t *testing.T) {
	repository := newFakeAuthorRepo()
	service := NewService(repository)

	_, err := service.Add(context.Background(), CreateAuthorRequest{
		Biography:   "Author biography",
		Birth:       1958,
		Nationality: "American",
	})
	var domainErr *domain.Error
	require.ErrorAs(t, err, &domainErr)
	require.Equal(t, domain.KindValidation, domainErr.Kind)
	require.Contains(t, domainErr.Message, "Request validation failed")
}

func TestRejectsAuthorNameWithOnlySpaces(t *testing.T) {
	repository := newFakeAuthorRepo()
	service := NewService(repository)

	_, err := service.Add(context.Background(), CreateAuthorRequest{
		Name:        "     ",
		Biography:   "Author biography",
		Birth:       1958,
		Nationality: "American",
	})

	var domainErr *domain.Error
	require.ErrorAs(t, err, &domainErr)
	require.Equal(t, domain.KindValidation, domainErr.Kind)
	require.Contains(t, domainErr.Message, "Request validation failed")
}

func TestRejectsAuthorDuplicateName(t *testing.T) {
	repository := newFakeAuthorRepo()
	repository.byName["Mitch Albom"] = Author{Name: "Mitch Albom"}
	service := NewService(repository)

	_, err := service.Add(context.Background(), CreateAuthorRequest{
		Name:        "Mitch Albom",
		Biography:   "I am author of book called Tuesdays with morrie",
		Birth:       1958,
		Nationality: "American",
	})
	var domainErr *domain.Error
	require.ErrorAs(t, err, &domainErr)
	require.Equal(t, domain.KindConflict, domainErr.Kind)
	require.Contains(t, domainErr.Message, "Author already exists")
}

func TestRejectsEmptyName(t *testing.T) {
	repository := newFakeAuthorRepo()
	service := NewService(repository)

	_, err := service.Add(context.Background(), CreateAuthorRequest{
		Name:        "",
		Biography:   "I am author of book called Tuesdays with morrie",
		Birth:       1958,
		Nationality: "American",
	})
	var domainErr *domain.Error
	require.ErrorAs(t, err, &domainErr)
	require.Equal(t, len(domainErr.Violations), 1)
	require.Contains(t, domainErr.Message, "Request validation failed")
}

func TestRejectsNameOutOfRange(t *testing.T) {
	repository := newFakeAuthorRepo()
	service := NewService(repository)

	_, err := service.Add(context.Background(), CreateAuthorRequest{
		Name:        strings.Repeat("Mitch Albom", 50),
		Biography:   "I am author of book called Tuesdays with morrie",
		Birth:       1958,
		Nationality: "American",
	})
	var domainErr *domain.Error
	require.ErrorAs(t, err, &domainErr)
	require.Equal(t, len(domainErr.Violations), 1)
	require.Contains(t, domainErr.Message, "Request validation failed")
}

func TestAuthorNameTrimmedSpaces(t *testing.T) {
	repository := newFakeAuthorRepo()
	service := NewService(repository)

	response, err := service.Add(context.Background(), CreateAuthorRequest{
		Name:        "  Mitch Albom  ",
		Biography:   "I am author of book called Tuesdays with morrie",
		Birth:       1958,
		Nationality: "American",
	})
	require.NoError(t, err)
	require.Equal(t, "Mitch Albom", response.Name)
	require.Equal(t, "I am author of book called Tuesdays with morrie", response.Biography)
	require.Equal(t, 1958, response.Birth)
	require.Equal(t, "American", response.Nationality)
}

func TestAuthorBiographyIsOptional(t *testing.T) {
	repository := newFakeAuthorRepo()
	service := NewService(repository)

	response, err := service.Add(context.Background(), CreateAuthorRequest{
		Name:        "Mitch Albom",
		Birth:       1958,
		Nationality: "American",
	})
	require.NoError(t, err)
	require.Equal(t, "Mitch Albom", response.Name)
	require.Empty(t, response.Biography)
	require.Equal(t, 1958, response.Birth)
	require.Equal(t, "American", response.Nationality)
}

func TestAuthorBiographyTrimmedSpaces(t *testing.T) {
	repository := newFakeAuthorRepo()
	service := NewService(repository)

	response, err := service.Add(context.Background(), CreateAuthorRequest{
		Name:        "Mitch Albom",
		Biography:   "   I am author of book called Tuesdays with morrie     ",
		Birth:       1958,
		Nationality: "American",
	})
	require.NoError(t, err)
	require.Equal(t, "Mitch Albom", response.Name)
	require.Equal(t, "I am author of book called Tuesdays with morrie", response.Biography)
	require.Equal(t, 1958, response.Birth)
	require.Equal(t, "American", response.Nationality)
}

func TestRejectsBiographyOutOfRange(t *testing.T) {
	repository := newFakeAuthorRepo()
	service := NewService(repository)

	_, err := service.Add(context.Background(), CreateAuthorRequest{
		Name:        "Mitch Albom",
		Biography:   strings.Repeat("I am author of book called Tuesdays with morrie", 100),
		Birth:       1958,
		Nationality: "American",
	})
	var domainErr *domain.Error
	require.ErrorAs(t, err, &domainErr)
	require.Equal(t, len(domainErr.Violations), 1)
	require.Contains(t, domainErr.Message, "Request validation failed")
}

func TestRejectsInvalidYearOfBirthLesserThanZero(t *testing.T) {
	repository := newFakeAuthorRepo()
	service := NewService(repository)

	_, err := service.Add(context.Background(), CreateAuthorRequest{
		Name:        "Mitch Albom",
		Biography:   "I am author of book called Tuesdays with morrie",
		Birth:       -1,
		Nationality: "American",
	})
	var domainErr *domain.Error
	require.ErrorAs(t, err, &domainErr)
	require.Equal(t, len(domainErr.Violations), 1)
	require.Contains(t, domainErr.Message, "Request validation failed")
}

func TestRejectsInvalidYearOfBirthLIsZero(t *testing.T) {
	repository := newFakeAuthorRepo()
	service := NewService(repository)

	_, err := service.Add(context.Background(), CreateAuthorRequest{
		Name:        "Mitch Albom",
		Biography:   "I am author of book called Tuesdays with morrie",
		Birth:       0,
		Nationality: "American",
	})
	var domainErr *domain.Error
	require.ErrorAs(t, err, &domainErr)
	require.Equal(t, len(domainErr.Violations), 1)
	require.Contains(t, domainErr.Message, "Request validation failed")
}

func TestRejectsInvalidYearOfBirthGreaterCanCurrentYear(t *testing.T) {
	repository := newFakeAuthorRepo()
	service := NewService(repository)

	_, err := service.Add(context.Background(), CreateAuthorRequest{
		Name:        "Mitch Albom",
		Biography:   "I am author of book called Tuesdays with morrie",
		Birth:       2027,
		Nationality: "American",
	})
	var domainErr *domain.Error
	require.ErrorAs(t, err, &domainErr)
	require.Equal(t, len(domainErr.Violations), 1)
	require.Contains(t, domainErr.Message, "Request validation failed")
}

func TestAuthorNationalityTrimmedSpaces(t *testing.T) {
	repository := newFakeAuthorRepo()
	service := NewService(repository)

	response, err := service.Add(context.Background(), CreateAuthorRequest{
		Name:        "Mitch Albom",
		Biography:   "   I am author of book called Tuesdays with morrie     ",
		Birth:       1958,
		Nationality: "American",
	})
	require.NoError(t, err)
	require.Equal(t, "Mitch Albom", response.Name)
	require.Equal(t, "I am author of book called Tuesdays with morrie", response.Biography)
	require.Equal(t, 1958, response.Birth)
	require.Equal(t, "American", response.Nationality)
}

func TestRejectsCountryNameOutOfRange(t *testing.T) {
	repository := newFakeAuthorRepo()
	service := NewService(repository)

	_, err := service.Add(context.Background(), CreateAuthorRequest{
		Name:        "Mitch Albom",
		Biography:   "I am author of book called Tuesdays with morrie",
		Birth:       1987,
		Nationality: strings.Repeat("American", 30),
	})
	var domainErr *domain.Error
	require.ErrorAs(t, err, &domainErr)
	require.Equal(t, len(domainErr.Violations), 1)
	require.Contains(t, domainErr.Message, "Request validation failed")
}

func TestAuthorNationalityIsOptional(t *testing.T) {
	repository := newFakeAuthorRepo()
	service := NewService(repository)

	response, err := service.Add(context.Background(), CreateAuthorRequest{
		Name:      "Mitch Albom",
		Biography: "I am author of book called Tuesdays with morrie",
		Birth:     1958,
	})
	require.NoError(t, err)
	require.Equal(t, "Mitch Albom", response.Name)
	require.Equal(t, "I am author of book called Tuesdays with morrie", response.Biography)
	require.Equal(t, 1958, response.Birth)
	require.Empty(t, response.Nationality)
}
func TestViewAuthor(t *testing.T) {
	repository := newFakeAuthorRepo()
	authorID := "d2e6b13b-e1b9-48c9-a623-84c034a1912e"
	repository.byName["Palle Vasu"] = Author{
		Id:          authorID,
		Name:        "Palle Vasu",
		Biography:   "Author of WE WERE NEVER MEANT TO BE",
		Birth:       1958,
		Nationality: "American",
	}
	service := NewService(repository)
	response, err := service.View(context.Background(), authorID)
	require.NoError(t, err)
	require.Equal(t, "Palle Vasu", response.Name)
	require.Equal(t, "Author of WE WERE NEVER MEANT TO BE", response.Biography)
	require.Equal(t, 1958, response.Birth)
	require.Equal(t, "American", response.Nationality)
}
func TestViewAuthorNotFound(t *testing.T) {
	repository := newFakeAuthorRepo()
	service := NewService(repository)

	_, err := service.View(context.Background(), "Unknown Author")

	var domainErr *domain.Error
	require.ErrorAs(t, err, &domainErr)
	require.Equal(t, domain.KindNotFound, domainErr.Kind)
	require.Contains(t, domainErr.Message, "Author not found")
}
func TestSearchAuthors(t *testing.T) {
	repository := newFakeAuthorRepo()
	service := NewService(repository)

	_, err := service.Add(context.Background(), CreateAuthorRequest{
		Name:        "Mitch Albom",
		Biography:   "I am author of book called Tuesdays with morrie",
		Birth:       1958,
		Nationality: "American",
	})

	_, err1 := service.Add(context.Background(), CreateAuthorRequest{
		Name:        "Michel jackson",
		Biography:   "I am an author",
		Birth:       1989,
		Nationality: "Mexican",
	})

	require.NoError(t, err)
	require.NoError(t, err1)

	authors, err2 := service.Search(context.Background(), "Mi", 1, 2)

	require.NoError(t, err2)

	require.Equal(t, len(authors.Author), 2)
	require.Equal(t, authors.Total, 2)
}
func TestPageOneContainsTwoRecords(t *testing.T) {
	repository := newFakeAuthorRepo()
	service := NewService(repository)

	_, err := service.Add(context.Background(), CreateAuthorRequest{
		Name:        "Mitch Albom",
		Biography:   "I am author of book called Tuesdays with morrie",
		Birth:       1958,
		Nationality: "American",
	})

	_, err1 := service.Add(context.Background(), CreateAuthorRequest{
		Name:        "Michel jackson",
		Biography:   "I am an author",
		Birth:       1989,
		Nationality: "Mexican",
	})

	_, err2 := service.Add(context.Background(), CreateAuthorRequest{
		Name:        "orry",
		Biography:   "I am author of book called Tuesdays with morrie",
		Birth:       1958,
		Nationality: "American",
	})

	_, err3 := service.Add(context.Background(), CreateAuthorRequest{
		Name:        "yamini",
		Biography:   "I am an author",
		Birth:       1989,
		Nationality: "Mexican",
	})

	require.NoError(t, err)
	require.NoError(t, err1)
	require.NoError(t, err2)
	require.NoError(t, err3)

	authors, err4 := service.Search(context.Background(), "", 1, 2)

	require.NoError(t, err4)

	require.Equal(t, len(authors.Author), 2)
	require.Equal(t, authors.Total, 4)

}

func TestPageSecondContainsTwoRecords(t *testing.T) {
	repository := newFakeAuthorRepo()
	service := NewService(repository)

	_, err := service.Add(context.Background(), CreateAuthorRequest{
		Name:        "Mitch Albom",
		Biography:   "I am author of book called Tuesdays with morrie",
		Birth:       1958,
		Nationality: "American",
	})

	_, err1 := service.Add(context.Background(), CreateAuthorRequest{
		Name:        "Michel jackson",
		Biography:   "I am an author",
		Birth:       1989,
		Nationality: "Mexican",
	})

	_, err2 := service.Add(context.Background(), CreateAuthorRequest{
		Name:        "orry",
		Biography:   "I am author of book called Tuesdays with morrie",
		Birth:       1958,
		Nationality: "American",
	})

	_, err3 := service.Add(context.Background(), CreateAuthorRequest{
		Name:        "yamini",
		Biography:   "I am an author",
		Birth:       1989,
		Nationality: "Mexican",
	})

	require.NoError(t, err)
	require.NoError(t, err1)
	require.NoError(t, err2)
	require.NoError(t, err3)

	authors, err4 := service.Search(context.Background(), "", 1, 2)

	require.NoError(t, err4)

	require.Equal(t, 2, len(authors.Author))
	require.Equal(t, 4, authors.Total)

}
func TestAuthorExists(t *testing.T) {
	repository := newFakeAuthorRepo()
	service := NewService(repository)

	response, err := service.Add(context.Background(), CreateAuthorRequest{
		Name:        "Mitch Albom",
		Biography:   "I am author of book called Tuesdays with morrie",
		Birth:       1958,
		Nationality: "American",
	})
	require.NoError(t, err)
	require.Equal(t, "Mitch Albom", response.Name)
}
func TestAuthorRemovedSucessfully(t *testing.T) {
	repository := newFakeAuthorRepo()
	authorID := "d2e6b13b-e1b9-48c9-a623-84c034a1912e"

	repository.byName["Palle Vasu"] = Author{
		Id:          authorID,
		Name:        "Palle Vasu",
		Biography:   "Author of WE WERE NEVER MEANT TO BE",
		Birth:       1958,
		Nationality: "American",
	}

	service := NewService(repository)
	response, err := service.Remove(context.Background(), authorID)

	require.NoError(t, err)
	require.Equal(t, "Palle Vasu", response.Name)
	require.NotContains(t, repository.byName, "Palle Vasu")
}
func TestRejectsUnkownAuthors(t *testing.T) {
	repository := newFakeAuthorRepo()
	service := NewService(repository)

	_, err := service.Remove(context.Background(), "Unknown Author")

	var domainErr *domain.Error
	require.ErrorAs(t, err, &domainErr)
	require.Equal(t, domain.KindNotFound, domainErr.Kind)
	require.Contains(t, domainErr.Message, "Author not found")
}
func TestUpdateAuthorWithValidDetails(t *testing.T) {
	repository := newFakeAuthorRepo()
	service := NewService(repository)

	author, _ := service.Add(context.Background(), CreateAuthorRequest{
		Name:        "mitch albom",
		Biography:   "I am author of book called Tuesdays with morrie",
		Birth:       1958,
		Nationality: "American",
	})

	updatedAuthor := CreateAuthorRequest{
		Name:        "mitch albom",
		Biography:   "I am author of book called Tuesdays with morrie and for one more day",
		Birth:       1958,
		Nationality: "American",
	}

	response, err := service.Update(context.Background(), author.Id, updatedAuthor)

	require.NoError(t, err)
	require.Equal(t, "mitch albom", response.Name)
	require.Equal(t, "I am author of book called Tuesdays with morrie and for one more day", response.Biography)
	require.Equal(t, 1958, response.Birth)
	require.Equal(t, "American", response.Nationality)
}
func TestUpdateAuthorNotFound(t *testing.T) {
	repository := newFakeAuthorRepo()
	service := NewService(repository)

	updatedAuthor := CreateAuthorRequest{
		Name:        "chota bheem",
		Biography:   "indian author",
		Birth:       1990,
		Nationality: "Indian",
	}

	_, err := service.Update(context.Background(), "unknown-id", updatedAuthor)

	require.Error(t, err)
	require.Contains(t, err.Error(), "Author not found")
}
func TestAcceptAuthorDetailsTrimSpaces(t *testing.T) {
	repository := newFakeAuthorRepo()
	service := NewService(repository)

	author, _ := service.Add(context.Background(), CreateAuthorRequest{
		Name:        "mitch albom",
		Biography:   "I am author of book called Tuesdays with morrie",
		Birth:       1958,
		Nationality: "American",
	})

	updatedAuthor := CreateAuthorRequest{
		Name:        "    mitch albom",
		Biography:   "   I am author of book called Tuesdays with morrie and for one more day    ",
		Birth:       1958,
		Nationality: "American    ",
	}

	response, err := service.Update(context.Background(), author.Id, updatedAuthor)

	require.NoError(t, err)
	require.Equal(t, "mitch albom", response.Name)
	require.Equal(t, "I am author of book called Tuesdays with morrie and for one more day", response.Biography)
	require.Equal(t, 1958, response.Birth)
	require.Equal(t, "American", response.Nationality)
}
func TestAccepctsEmptyDetails(t *testing.T) {
	repository := newFakeAuthorRepo()
	service := NewService(repository)

	author, _ := service.Add(context.Background(), CreateAuthorRequest{
		Name:        "mitch albom",
		Biography:   "I am author of book called Tuesdays with morrie",
		Birth:       1958,
		Nationality: "American",
	})

	updatedAuthor := CreateAuthorRequest{
		Name:        "    mitch albom",
		Biography:   " ",
		Birth:       1958,
		Nationality: " ",
	}

	response, err := service.Update(context.Background(), author.Id, updatedAuthor)

	require.NoError(t, err)
	require.Equal(t, "mitch albom", response.Name)
	require.Equal(t, "", response.Biography)
	require.Equal(t, 1958, response.Birth)
	require.Equal(t, "", response.Nationality)

}
func TestRejectsUpdateEmptyName(t *testing.T) {
	repository := newFakeAuthorRepo()
	service := NewService(repository)

	author, _ := service.Add(context.Background(), CreateAuthorRequest{
		Name:        "mitch albom",
		Biography:   "I am author of book called Tuesdays with morrie",
		Birth:       1958,
		Nationality: "American",
	})

	updatedAuthor := CreateAuthorRequest{
		Name:        " ",
		Biography:   "I am author of book called Tuesdays with morrie",
		Birth:       1958,
		Nationality: "American",
	}

	_, err := service.Update(context.Background(), author.Id, updatedAuthor)

	var domainErr *domain.Error
	require.ErrorAs(t, err, &domainErr)
	require.Equal(t, len(domainErr.Violations), 1)
	require.Equal(t, domainErr.Message, "Request validation failed")
}

type fakeAuthorRepo struct {
	byName map[string]Author
}

func newFakeAuthorRepo() *fakeAuthorRepo {
	return &fakeAuthorRepo{byName: map[string]Author{}}
}

func (f *fakeAuthorRepo) ExistByName(_ context.Context, name string) (bool, error) {
	_, ok := f.byName[name]
	return ok, nil
}

func (f *fakeAuthorRepo) Create(_ context.Context, author Author) (Author, error) {
	if _, ok := f.byName[author.Name]; ok {
		return Author{}, ErrDuplicate
	}

	f.byName[author.Name] = author
	return author, nil
}

func (f *fakeAuthorRepo) ViewAuthor(ctx context.Context, id string) (Author, error) {
	for _, author := range f.byName {
		if author.Id == id {
			return author, nil
		}
	}
	return Author{}, ErrNotFound
}

func (f *fakeAuthorRepo) GetByName(_ context.Context, name string) (Author, error) {
	author, ok := f.byName[name]
	if !ok {
		return Author{}, ErrNotFound
	}
	return author, nil
}

func (f *fakeAuthorRepo) SearchAuthor(_ context.Context, searchName string, limit, offset int) ([]Author, error) {
	authors := make([]Author, 0)

	for _, author := range f.byName {
		if strings.HasPrefix(strings.ToLower(author.Name), strings.ToLower(searchName)) {
			authors = append(authors, author)
		}
	}

	if offset >= len(authors) {
		return []Author{}, nil
	}

	end := offset + limit

	if end > len(authors) {
		end = len(authors)
	}

	return authors[offset:end], nil
}

func (f *fakeAuthorRepo) CountAuthors(_ context.Context, name string) (int, error) {
	var count int

	for name := range f.byName {
		if strings.HasPrefix(f.byName[name].Name, name) {
			count += 1
		}
	}
	return count, nil
}

func (f *fakeAuthorRepo) HasBooksAuthors(_ context.Context, id string) (bool, error) {
	if id == "45e9bc8e-d9c0-4940-bc5d-e5c0b9f8473c" {
		return true, nil
	}
	return false, nil
}

func (f *fakeAuthorRepo) DeleteAuthor(_ context.Context, id string) (Author, error) {
	for name, author := range f.byName {
		if author.Id == id {
			delete(f.byName, name)
			return author, nil
		}
	}
	return Author{}, ErrNotFound
}

func (f *fakeAuthorRepo) HasName(_ context.Context, id string, name string) (bool, error) {
	for _, author := range f.byName {
		if author.Name == name && author.Id != id {
			return true, nil
		}
	}
	return false, nil
}

func (f *fakeAuthorRepo) UpdateAuthor(_ context.Context, author Author) (Author, error) {
	for name, a := range f.byName {
		if a.Id == author.Id {
			delete(f.byName, name)
			f.byName[author.Name] = author
			return author, nil
		}
	}
	return Author{}, ErrNotFound
}
