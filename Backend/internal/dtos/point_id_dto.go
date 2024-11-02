package dtos

import "io"

type PointIdDto struct {
	Id int64 `json:"id"`
}

func (self PointIdDto) Decode(r io.Reader) error {
	panic("not implemented") // TODO: Implement
}

func (self PointIdDto) Encode() (string, error) {
	panic("not implemented") // TODO: Implement
}

func (self PointIdDto) Validate() error {
	panic("not implemented") // TODO: Implement
}

