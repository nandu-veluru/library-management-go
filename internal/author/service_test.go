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

func TestRejectsEmptyAuthorName(t *testing.T) {
	repository := newFakeAuthorRepo()
	service := NewService(repository)

	_, err := service.Add(context.Background(), CreateAuthorRequest{
		Name:        "",
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

func TestRejectsNullName(t *testing.T) {
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
		Name:        "Mitch Albom",
		Biography:   "I am author of book called Tuesdays with morrie",
		Birth:       1958,
	})
	require.NoError(t, err)
	require.Equal(t, "Mitch Albom", response.Name)
	require.Equal(t, "I am author of book called Tuesdays with morrie", response.Biography)
	require.Equal(t, 1958, response.Birth)
	require.Empty(t, response.Nationality)
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
