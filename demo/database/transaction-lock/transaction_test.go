package transaction

import (
	"strconv"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPeople(t *testing.T) {
	Setup(
		"mysql",
		"root:toor@tcp(localhost:3308)/demo_test?charset=utf8&parseTime=True",
		false,
	)
	//DB.LogMode(true)

	Convey("clean database", t, func() {
		DB.DropTableIfExists(&People{})
		DB.CreateTable(&People{})

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
				So(Create(DB, &ps[i]), ShouldBeNil)
				ps[i].CreatedAt = ps[i].CreatedAt.UTC().Round(time.Second)
				ps[i].UpdatedAt = ps[i].UpdatedAt.UTC().Round(time.Second)
			}

			Convey("get peoples", func() {
				for i, p := range ps {
					getP, err := QueryByName(DB, p.Name)
					So(err, ShouldBeNil)
					So(getP, ShouldResemble, ps[i])
				}
			})

			Convey("get peoples by age range", func() {
				for i := range ps {
					for j := i; j < len(ps); j++ {
						getPs, err := QueryByAgeRange(DB, ps[i].Age, ps[j].Age)
						So(err, ShouldBeNil)
						So(getPs, ShouldResemble, ps[i:j+1])
					}
				}
			})

			/*
				Convey("get peoples by age range lock", func() {
						tx := DB.Begin()
						defer tx.Rollback()
						getPs, err := QueryByAgeRangeLock(tx, 1, 3)
						So(err, ShouldBeNil)
						So(getPs, ShouldResemble, ps[1:4])

						tx2 := DB.Begin()
						defer tx2.Rollback()
						getPs, err = QueryByAgeRangeLock(tx2, 6, 7)
						So(err, ShouldBeNil)
						So(getPs, ShouldResemble, ps[6:8])
					})
			*/

			Convey("get peoples by age range lock", func() {
				tx, err := XDB.Beginx()
				So(err, ShouldBeNil)
				defer tx.Rollback()
				getPs, err := QueryByAgeRangeLockX(tx, 1, 3)
				So(err, ShouldBeNil)
				So(getPs, ShouldResemble, ps[1:4])

				// Why not work?
				tx2, err := XDB.Beginx()
				So(err, ShouldBeNil)
				defer tx2.Rollback()
				getPs, err = QueryByAgeRangeLockX(tx2, 6, 7)
				So(err, ShouldBeNil)
				So(getPs, ShouldResemble, ps[6:8])
			})
		})
	})
}
