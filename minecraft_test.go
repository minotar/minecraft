// minecraft_test.go
package minecraft

import (
	"regexp"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

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

		invalidUUIDs := []string{"clone1018", "d9135e082f2244c8-9cb0-bee234155292"}
		validUUIDs := []string{"d9135e082f2244c89cb0bee234155292", "d9135e08-2f22-44c8-9cb0-bee234155292"}

		validUsernamesOrUUIDs := append(validUsernames, validUUIDs...)
		possiblyInvalidUsernamesOrUUIDs := append(invalidUsernames, invalidUUIDs...)

		usernameRegex := regexp.MustCompile("^" + ValidUsernameRegex + "$")
		uuidRegex := regexp.MustCompile("^" + ValidUuidRegex + "$")
		usernameOrUUIDRegex := regexp.MustCompile("^" + ValidUsernameOrUuidRegex + "$")

		Convey("Username regex works", func() {
			for _, validUsername := range validUsernames {
				So(usernameRegex.MatchString(validUsername), ShouldBeTrue)
			}

			for _, invalidUsername := range invalidUsernames {
				So(usernameRegex.MatchString(invalidUsername), ShouldBeFalse)
			}
		})

		Convey("UUID regex works", func() {
			for _, validUUID := range validUUIDs {
				So(uuidRegex.MatchString(validUUID), ShouldBeTrue)
			}

			for _, invalidUUID := range invalidUUIDs {
				So(uuidRegex.MatchString(invalidUUID), ShouldBeFalse)
			}
		})

		Convey("Username-or-UUID regex works", func() {
			for _, validThing := range validUsernamesOrUUIDs {
				So(usernameOrUUIDRegex.MatchString(validThing), ShouldBeTrue)
			}

			for _, possiblyInvalidThing := range possiblyInvalidUsernamesOrUUIDs {
				resultOne := usernameRegex.MatchString(possiblyInvalidThing)
				resultTwo := uuidRegex.MatchString(possiblyInvalidThing)
				expectedResult := resultOne || resultTwo

				So(usernameOrUUIDRegex.MatchString(possiblyInvalidThing), ShouldEqual, expectedResult)
			}
		})

	})

}

func TestExtra(t *testing.T) {

	Convey("Test apiRequest", t, func() {

		Convey("Not a URL", func() {
			_, err := apiRequest("//")

			So(err.Error(), ShouldContainSubstring, "Unable to Get URL")
		})

	})
}
