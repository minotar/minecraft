// minecraft_test.go
package minecraft

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGetUser(t *testing.T) {

	Convey("CLone1018 should equal clone1018", t, func() {
		user, _ := GetUser("CLone1018")

		So(user.Name, ShouldEqual, "clone1018")
	})

	Convey("skmkj88200aklk should gracefully error", t, func() {
		_, err := GetUser("skmkj88200aklk")

		So(err.Error(), ShouldEqual, "User not found.")
	})

}
