package transaction

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
)

type People struct {
	ID        uint       `gorm:"primary_key"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `sql:"index" db:"deleted_at"`

	Name string
	Age  int
}

var DB *gorm.DB
var XDB *sqlx.DB

func Setup(driver string, source string, enableMigrate bool) {
	xdb, err := sqlx.Open(driver, source)
	if err != nil {
		panic(err)
	}
	XDB = xdb
	db, err := gorm.Open(xdb.DriverName(), xdb.DB)
	if err != nil {
		panic(err)
	}
	DB = db

	if enableMigrate {
		if err := db.AutoMigrate(&People{}).Error; err != nil {
			panic(err)
		}
	}
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

func QueryByAgeRangeLockX(db sqlx.Ext, min, max int) (ps []People, err error) {
	err = sqlx.Select(
		db, &ps,
		`select * from peoples where age between ? and ? for update skip locked`,
		min, max)
	return
}
