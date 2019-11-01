package database

import (
	"strconv"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/rubenv/sql-migrate"
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
				m := &migrate.MemoryMigrationSource{
					Migrations: []*migrate.Migration{
						{
							Id:   "001",
							Up:   up(db.DriverName()),
							Down: down(db.DriverName()),
						},
					},
				}
				_, err = migrate.Exec(db.DB, db.DriverName(), m, migrate.Down)
				So(err, ShouldBeNil)
				_, err = migrate.Exec(db.DB, db.DriverName(), m, migrate.Up)
				So(err, ShouldBeNil)

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

					Convey("get peoples by age range lock no cross", func() {
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
						getPs, err = QueryByAgeRangeLock(tx2, 4, 6)
						So(err, ShouldBeNil)

						for k := range getPs {
							TimeTruncate(&getPs[k])
						}
						So(getPs, ShouldResemble, ps[4:7])
					})

					Convey("get peoples by age range lock cross", func() {
						tx, err := db.Beginx()
						So(err, ShouldBeNil)
						defer func() {
							So(tx.Rollback(), ShouldBeNil)
						}()
						getPs, err := QueryByAgeRangeLock(tx, 1, 4)
						So(err, ShouldBeNil)

						for k := range getPs {
							TimeTruncate(&getPs[k])
						}
						So(getPs, ShouldResemble, ps[1:5])

						tx2, err := db.Beginx()
						So(err, ShouldBeNil)
						defer func() {
							So(tx2.Rollback(), ShouldBeNil)
						}()
						getPs, err = QueryByAgeRangeLock(tx2, 3, 6)
						So(err, ShouldBeNil)

						for k := range getPs {
							TimeTruncate(&getPs[k])
						}
						So(getPs, ShouldResemble, ps[5:7])
					})

					Convey("get peoples by age range lock has limit", func() {
						tx, err := db.Beginx()
						So(err, ShouldBeNil)
						defer func() {
							So(tx.Rollback(), ShouldBeNil)
						}()
						getPs, err := QueryByAgeRangeLimitLock(tx, 1, 6, 4)
						So(err, ShouldBeNil)

						for k := range getPs {
							TimeTruncate(&getPs[k])
						}
						So(getPs, ShouldResemble, ps[1:5])

						tx2, err := db.Beginx()
						So(err, ShouldBeNil)
						defer func() {
							So(tx2.Rollback(), ShouldBeNil)
						}()
						getPs, err = QueryByAgeRangeLimitLock(tx2, 1, 6, 3)
						So(err, ShouldBeNil)

						for k := range getPs {
							TimeTruncate(&getPs[k])
						}
						So(getPs, ShouldResemble, ps[5:7])
					})
				})
			})
		})
	}
}
