package model

import (
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type Blob struct {
	Model

	UXID string `json:"uxid" gorm:"uniqueIndex;column:uxid"`
	Name string `json:"name"`
}

func (b *Blob) BeforeCreate(tx *gorm.DB) error {
	if b.UXID == "" {
		uxid, err := uuid.NewV4()
		if err != nil {
			return err
		}
		b.UXID = uxid.String()
	}
	return nil
}
