// minecraft_test.go
package minecraft

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestProfiles(t *testing.T) {

	Convey("clone1018 should match d9135e082f2244c89cb0bee234155292", t, func() {
		user, _ := GetUser("clone1018")

		So(user.Id, ShouldEqual, "d9135e082f2244c89cb0bee234155292")
	})

	Convey("CLone1018 should equal clone1018", t, func() {
		user, _ := GetUser("CLone1018")

		So(user.Name, ShouldEqual, "clone1018")
	})

	Convey("skmkj88200aklk should gracefully error", t, func() {
		_, err := GetUser("skmkj88200aklk")

		So(err.Error(), ShouldEqual, "User not found.")
	})

}

func TestAvatars(t *testing.T) {

	Convey("clone1018 should return valid image", t, func() {
		user := User{Name: "clone1018"}

		skin := GetSkin(user)

		So(skin, ShouldNotBeNil)
	})

}
