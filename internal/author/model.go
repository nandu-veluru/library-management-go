package author

type Author struct {
	Id          string
	Name        string
	Biography   string
	Birth       int
	Nationality string
}

type CreateAuthorRequest struct {
	Name        string `json:"name"`
	Biography   string `json:"biography"`
	Birth       int    `json:"birth"`
	Nationality string `json:"nationality"`
}

type AuthorResponse struct {
	Name        string `json:"name"`
	Biography   string `json:"biography"`
	Birth       int    `json:"birth"`
	Nationality string `json:"nationality"`
}
