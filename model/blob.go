package model

import (
	"io"
)

type Blob struct {
	ID       int64     `json:"id"`
	FileName string    `json:"filename"`
	Reader   io.Reader `json:"-"`
}

func NewBlob(filename string, reader io.Reader) *Blob {
	return &Blob{
		FileName: filename,
		Reader:   reader,
	}
}
