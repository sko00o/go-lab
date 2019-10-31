package database

import (
	"time"

	"github.com/jinzhu/gorm"
)

type People struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	Name string
	Age  int
}

func Create(db *gorm.DB, p *People) error {
	return db.Create(p).Error
}

func QueryByName(db *gorm.DB, name string) (p People, err error) {
	err = db.First(&p, People{Name: name}).Error
	return
}

func QueryByAgeRange(db *gorm.DB, min, max int) (ps []People, err error) {
	err = db.Find(&ps, "age between ? and ?", min, max).Error
	return
}

func QueryByAgeRangeLock(db *gorm.DB, min, max int) (ps []People, err error) {
	err = db.Set("gorm:query_option", "FOR UPDATE SKIP LOCKED").
		Find(&ps, "age between ? and ?", min, max).Error
	return
}
