package main

import (
	"github.com/jinzhu/gorm"
)

type StorageBlob struct {
	gorm.Model
	//Key string `gorm:"unique;not null"`
	Filename string `gorm:"not null"`
}
