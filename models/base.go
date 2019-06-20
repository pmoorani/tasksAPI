package models

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
)

type BaseModel struct {
	ID    uuid.UUID `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func (base *BaseModel) BeforeCreate(scope *gorm.Scope) (err error) {
	id := uuid.NewV4()

	scope.SetColumn("ID", id)
	return nil
}