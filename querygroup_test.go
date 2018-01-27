package godb

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestQueryGroupDo(t *testing.T) {
	Convey("Given a test database", t, func() {
		db := fixturesSetup(t)
		defer db.Close()

		Convey("Do execute the raw query and fills a given instance", func() {
			qg := db.NewQueryGroup()
			allDummy := []Dummy{}
			qr1 := qg.Add(db.RawSQL("select * from dummies").DoLater(&allDummy))
			So(qr1.Err(), ShouldBeNil)

			singleDummy := Dummy{}
			qr2 := qg.Add(db.RawSQL("select * from dummies where an_integer = ?", 12).DoLater(&singleDummy))
			So(qr2.Err(), ShouldBeNil)

			err := qg.Do()
			So(err, ShouldBeNil)
			fmt.Println(allDummy)
			fmt.Println(singleDummy)
		})
	})
}
