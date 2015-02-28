// profiles_test.go
package minecraft

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestProfiles(t *testing.T) {

	// This will also effectively test GetUUID
	Convey("Test GetAPIProfile", t, func() {

		Convey("clone1018 should match d9135e082f2244c89cb0bee234155292", func() {
			apiProfile, err := GetAPIProfile("clone1018")

			So(err, ShouldBeNil)
			So(apiProfile.UUID, ShouldEqual, "d9135e082f2244c89cb0bee234155292")
		})

		Convey("CLone1018 should equal clone1018", func() {
			apiProfile, err := GetAPIProfile("CLone1018")

			So(err, ShouldBeNil)
			So(apiProfile.Username, ShouldEqual, "clone1018")
		})

		Convey("skmkj88200aklk should gracefully error", func() {
			apiProfile, err := GetAPIProfile("skmkj88200aklk")

			So(err.Error(), ShouldContainSubstring, "User not found")
			So(apiProfile, ShouldResemble, APIProfileResponse{})
		})

		Convey("bad_string/ should cause an HTTP error", func() {
			apiProfile, err := GetAPIProfile("bad_string/")

			So(err.Error(), ShouldContainSubstring, "Error retrieving profile")
			So(apiProfile, ShouldResemble, APIProfileResponse{})
		})

	})

	// Must be careful to not request same profile from session server more than once per ~30 seconds
	Convey("Test GetSessionProfile", t, func() {

		Convey("5c115ca73efd41178213a0aff8ef11e0 should equal LukeHandle", func() {
			// LukeHandle
			sessionProfile, err := GetSessionProfile("5c115ca73efd41178213a0aff8ef11e0")

			So(err, ShouldBeNil)
			So(sessionProfile.Username, ShouldEqual, "LukeHandle")
		})

		Convey("bad_string/ should cause an HTTP error", func() {
			sessionProfile, err := GetSessionProfile("bad_string/")

			So(err.Error(), ShouldContainSubstring, "Error retrieving profile")
			So(sessionProfile, ShouldResemble, SessionProfileResponse{})
		})

	})

	// Test a lot of what we did above, but this is a wrapper function that includes
	// common logic for solving the issues of being supplied with UUID and
	// Usernames and returning a uniform response (UUID of certain format)
	Convey("Test NormalizePlayerForUUID", t, func() {

		Convey("clone1018 should match d9135e082f2244c89cb0bee234155292", func() {
			playerUUID, err := NormalizePlayerForUUID("clone1018")

			So(err, ShouldBeNil)
			So(playerUUID, ShouldEqual, "d9135e082f2244c89cb0bee234155292")
		})

		Convey("CLone1018 should match d9135e082f2244c89cb0bee234155292", func() {
			playerUUID, err := NormalizePlayerForUUID("clone1018")

			So(err, ShouldBeNil)
			So(playerUUID, ShouldEqual, "d9135e082f2244c89cb0bee234155292")
		})

		Convey("d9135e08-2f22-44c8-9cb0-bee234155292 should match d9135e082f2244c89cb0bee234155292", func() {
			playerUUID, err := NormalizePlayerForUUID("d9135e082f2244c89cb0bee234155292")

			So(err, ShouldBeNil)
			So(playerUUID, ShouldEqual, "d9135e082f2244c89cb0bee234155292")
		})

		Convey("d9135e082f2244c89cb0bee234155292 should match d9135e082f2244c89cb0bee234155292", func() {
			playerUUID, err := NormalizePlayerForUUID("d9135e082f2244c89cb0bee234155292")

			So(err, ShouldBeNil)
			So(playerUUID, ShouldEqual, "d9135e082f2244c89cb0bee234155292")
		})

		Convey("skmkj88200aklk should gracefully error", func() {
			playerUUID, err := NormalizePlayerForUUID("skmkj88200aklk")

			So(err.Error(), ShouldContainSubstring, "User not found")
			So(playerUUID, ShouldBeBlank)
		})

		Convey("TooLongForAUsername should gracefully error", func() {
			playerUUID, err := NormalizePlayerForUUID("TooLongForAUsername")

			So(err.Error(), ShouldContainSubstring, "Invalid Username or UUID.")
			So(playerUUID, ShouldBeBlank)
		})

	})

}
