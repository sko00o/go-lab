package database

import (
	"strconv"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"xorm.io/xorm"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/ziutek/mymysql/godrv"

	"github.com/sko00o/go-lab/demo/database/test"
)

func TestPeople(t *testing.T) {
	for driver, source := range test.DSNTable {
		Convey("test by "+driver, t, func() {
			db, err := xorm.NewEngine(driver, source)
			So(err, ShouldBeNil)

			//db.ShowSQL(true)

			Reset(func() {
				So(db.Close(), ShouldBeNil)
			})

			Convey("clean database", func() {
				table := new(People)
				exist, err := db.IsTableExist(table)
				So(err, ShouldBeNil)
				if exist {
					So(db.DropTables(table), ShouldBeNil)
				}
				So(db.CreateTables(table), ShouldBeNil)

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
						tx := db.NewSession()
						So(tx.Begin(), ShouldBeNil)
						defer func() {
							So(tx.Rollback(), ShouldBeNil)
						}()
						getPs, err := QueryByAgeRangeLock(tx, 1, 3)
						So(err, ShouldBeNil)

						for k := range getPs {
							TimeTruncate(&getPs[k])
						}
						So(getPs, ShouldResemble, ps[1:4])

						tx2 := db.NewSession()
						So(tx2.Begin(), ShouldBeNil)
						defer func() {
							So(tx2.Rollback(), ShouldBeNil)
						}()
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
