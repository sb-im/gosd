package main

import (
	"encoding/json"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB

type Plan struct {
	gorm.Model
	Name           string
	Description    string
	File           string
	Node_id        string
	Cycle_types_id string
}

func (plan *Plan) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID             uint   `json:"id"`
		Name           string `json:"name"`
		Description    string `json:"description"`
		File           string `json:"file"`
		Node_id        string `json:"node_id"`
		Cycle_types_id string `json:"cycle_types_id"`
	}{
		ID:             plan.ID,
		Name:           plan.Name,
		Description:    plan.Description,
		File:           plan.File,
		Node_id:        plan.Node_id,
		Cycle_types_id: plan.Cycle_types_id,
	})
}

func DBlink() {
	db, _ = gorm.Open("sqlite3", "test.db")
	//db, err := gorm.Open("sqlite3", "test.db")
	//if err != nil {
	//  panic("failed to connect database")
	//}
	//defer db.Close()

	db.AutoMigrate(&Plan{})
}

func (this *Plan) Find(id int) {
	db.First(this, id)
}

func (this *Plan) Create() {
	db.Create(this)
	//db.Create(&Plan{Name: "L1212"})

	// Update - update product's price to 2000
	//db.Model(&product).Update("Price", 2000)

	// Delete - delete product
	//db.Delete(&product)
}
