// minecraft_test.go
package minecraft

import (
	. "github.com/smartystreets/goconvey/convey"
	"regexp"
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

		skin, _ := GetSkin(user)

		So(skin, ShouldNotBeNil)
	})

	Convey("d9135e082f2244c89cb0bee234155292 should return valid image", t, func() {
		skin, err := FetchSkinFromMojangByUuid("d9135e082f2244c89cb0bee234155292")

		So(skin, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})

	Convey("Wooxye should err", t, func() {
		user := User{Name: "Wooxye"}

		_, err := GetSkin(user)
		So(err.Error(), ShouldStartWith, "Skin not found.")
	})

	Convey("Char should return valid image", t, func() {
		charImg, err := FetchImageForChar()

		So(charImg, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})

}

func TestRegexs(t *testing.T) {
	Convey("Regexs compile", t, func() {
		var err error

		_, err = regexp.Compile(ValidUsernameRegex)
		So(err, ShouldBeNil)

		_, err = regexp.Compile(ValidUuidRegex)
		So(err, ShouldBeNil)

		_, err = regexp.Compile(ValidUsernameOrUuidRegex)
		So(err, ShouldBeNil)
	})

	Convey("Regexs work", t, func() {
		invalidUsernames := []string{"d9135e082f2244c89cb0bee234155292", "_-proscope-_", "PeriScopeButTooLong"}
		validUsernames := []string{"clone1018", "lukegb", "Wooxye"}

		invalidUuids := []string{"clone1018"}
		validUuids := []string{"d9135e082f2244c89cb0bee234155292"}

		validUsernamesOrUuids := append(validUsernames, validUuids...)
		possiblyInvalidUsernamesOrUuids := append(invalidUsernames, invalidUuids...)

		usernameRegex := regexp.MustCompile("^" + ValidUsernameRegex + "$")
		uuidRegex := regexp.MustCompile("^" + ValidUuidRegex + "$")
		usernameOrUuidRegex := regexp.MustCompile("^" + ValidUsernameOrUuidRegex + "$")

		Convey("Username regex works", func() {
			for _, validUsername := range validUsernames {
				So(usernameRegex.MatchString(validUsername), ShouldBeTrue)
			}

			for _, invalidUsername := range invalidUsernames {
				So(usernameRegex.MatchString(invalidUsername), ShouldBeFalse)
			}
		})

		Convey("UUID regex works", func() {

			for _, validUuid := range validUuids {
				So(uuidRegex.MatchString(validUuid), ShouldBeTrue)
			}

			for _, invalidUuid := range invalidUuids {
				So(uuidRegex.MatchString(invalidUuid), ShouldBeFalse)
			}
		})

		Convey("Username-or-UUID regex works", func() {

			for _, validThing := range validUsernamesOrUuids {
				So(usernameOrUuidRegex.MatchString(validThing), ShouldBeTrue)
			}

			for _, possiblyInvalidThing := range possiblyInvalidUsernamesOrUuids {
				resultOne := usernameRegex.MatchString(possiblyInvalidThing)
				resultTwo := uuidRegex.MatchString(possiblyInvalidThing)
				expectedResult := resultOne || resultTwo

				So(usernameOrUuidRegex.MatchString(possiblyInvalidThing), ShouldEqual, expectedResult)
			}
		})
	})
}
