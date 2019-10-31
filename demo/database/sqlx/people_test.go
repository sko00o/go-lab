package database

import (
	"strconv"
	"testing"

	"github.com/jmoiron/sqlx"
	. "github.com/smartystreets/goconvey/convey"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"

	"github.com/sko00o/go-lab/demo/database/test"
)

func TestPeople(t *testing.T) {
	for driver, source := range test.DSNTable {
		Convey("test by "+driver, t, func() {
			db, err := sqlx.Open(driver, source)
			So(err, ShouldBeNil)

			Reset(func() {
				So(db.Close(), ShouldBeNil)
			})

			Convey("clean database", func() {
				db.MustExec(down(db.DriverName()))
				db.MustExec(up(db.DriverName()))

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
						tx, err := db.Beginx()
						So(err, ShouldBeNil)
						defer func() {
							So(tx.Rollback(), ShouldBeNil)
						}()
						getPs, err := QueryByAgeRangeLock(tx, 1, 3)
						So(err, ShouldBeNil)

						for k := range getPs {
							TimeTruncate(&getPs[k])
						}
						So(getPs, ShouldResemble, ps[1:4])

						tx2, err := db.Beginx()
						So(err, ShouldBeNil)
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
