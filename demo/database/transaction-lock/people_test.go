package transaction

import (
	"strconv"
	"testing"

	"github.com/jinzhu/gorm"
	. "github.com/smartystreets/goconvey/convey"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var dsnTable = map[string]string{
	"mysql":    "root:toor@tcp(localhost:13306)/demo_test?parseTime=True",
	"postgres": "dbname=demo_test user=root password=toor port=15432 sslmode=disable",
}

func TestPeople(t *testing.T) {
	for driver, source := range dsnTable {
		Convey("test by "+driver, t, func() {
			db, err := gorm.Open(driver, source)
			So(err, ShouldBeNil)

			Reset(func() {
				So(db.Close(), ShouldBeNil)
			})

			Convey("clean database", func() {
				db.DropTableIfExists(&People{})
				db.CreateTable(&People{})

				Convey("create peoples", func() {
					// make people
					var ps []People
					for i := 0; i < 10; i++ {
						ps = append(ps, People{
							Name: strconv.Itoa(i),
							Age:  i,
						})
					}

					for i := range ps {
						So(Create(db, &ps[i]), ShouldBeNil)
						TimeTruncate(&ps[i])
					}

					Convey("get peoples", func() {
						for i, p := range ps {
							getP, err := QueryByName(db, p.Name)
							TimeTruncate(&getP)

							So(err, ShouldBeNil)
							So(getP, ShouldResemble, ps[i])
						}
					})

					Convey("get peoples by age range", func() {
						for i := range ps {
							for j := i; j < len(ps); j++ {
								getPs, err := QueryByAgeRange(db, ps[i].Age, ps[j].Age)
								for k := range getPs {
									TimeTruncate(&getPs[k])
								}

								So(err, ShouldBeNil)
								So(getPs, ShouldResemble, ps[i:j+1])
							}
						}
					})

					Convey("get peoples by age range lock", func() {
						tx := db.Begin()
						defer tx.Rollback()
						getPs, err := QueryByAgeRangeLock(tx, 1, 3)
						So(err, ShouldBeNil)

						for k := range getPs {
							TimeTruncate(&getPs[k])
						}
						So(getPs, ShouldResemble, ps[1:4])

						tx2 := db.Begin()
						defer tx2.Rollback()
						getPs, err = QueryByAgeRangeLock(tx2, 1, 7)
						So(err, ShouldBeNil)

						for k := range getPs {
							TimeTruncate(&getPs[k])
						}
						So(getPs, ShouldResemble, ps[4:8])
					})
				})
			})
		})
	}
}
