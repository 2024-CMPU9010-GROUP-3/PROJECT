package dtos

import "io"

type ResponseContentDto struct {
	HttpStatus int `json:"-"`
	Content    any `json:"content"`
}

func (self ResponseContentDto) Decode(r io.Reader) error {
	panic("not implemented") // TODO: Implement
}

func (self ResponseContentDto) Encode() (string, error) {
	panic("not implemented") // TODO: Implement
}

func (self ResponseContentDto) Validate() error {
	panic("not implemented") // TODO: Implement
}

