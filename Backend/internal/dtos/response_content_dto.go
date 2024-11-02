package dtos

type ResponseContentDto struct {
	HttpStatus int `json:"-"`
	Content    any `json:"content"`
}