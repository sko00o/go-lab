package database

import (
	"time"

	"xorm.io/xorm"
)

type People struct {
	ID        uint       `xorm:"pk serial 'id'"`
	CreatedAt time.Time  `xorm:"datetime created 'created_at'"`
	UpdatedAt time.Time  `xorm:"datetime updated 'updated_at'"`
	DeletedAt *time.Time `xorm:"datetime deleted 'deleted_at'"`

	Name string
	Age  int
}

func (People) TableName() string {
	return "peoples"
}

func Create(db *xorm.Engine, p *People) (err error) {
	_, err = db.Insert(p)
	return
}

func QueryByName(db *xorm.Engine, name string) (p People, err error) {
	_, err = db.Where("name = ?", name).Get(&p)
	return
}

func QueryByAgeRange(db *xorm.Engine, min, max int) (ps []People, err error) {
	err = db.Where("age between ? and ?", min, max).Find(&ps)
	return
}

func QueryByAgeRangeLock(db *xorm.Session, min, max int) (ps []People, err error) {
	err = db.SQL(`
		SELECT * 
		FROM peoples
		WHERE age BETWEEN ? AND ?
		FOR UPDATE SKIP LOCKED`,
		min, max,
	).Find(&ps)

	return
}
