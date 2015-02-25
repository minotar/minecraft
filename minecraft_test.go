// minecraft_test.go
package minecraft

import (
	"regexp"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestProfiles(t *testing.T) {

	Convey("clone1018 should match d9135e082f2244c89cb0bee234155292", t, func() {
		uuid, err := GetUUID("clone1018")

		So(err, ShouldBeNil)
		So(uuid, ShouldEqual, "d9135e082f2244c89cb0bee234155292")
	})

	Convey("CLone1018 should equal clone1018", t, func() {
		apiProfile, err := GetAPIProfile("CLone1018")

		So(err, ShouldBeNil)
		So(apiProfile.Username, ShouldEqual, "clone1018")
	})

	Convey("skmkj88200aklk should gracefully error", t, func() {
		apiProfile, err := GetAPIProfile("skmkj88200aklk")

		So(err.Error(), ShouldStartWith, "User not found.")
		So(apiProfile, ShouldResemble, APIProfileResponse{})
	})

	Convey("bad_string/ should cause an HTTP error", t, func() {
		sessionProfile, err := GetSessionProfile("bad_string/")

		So(err.Error(), ShouldStartWith, "Error retrieving profile.")
		So(sessionProfile, ShouldResemble, SessionProfileResponse{})
	})

}

func TestTextures(t *testing.T) {

	Convey("clone1018 texture should return the correct skin", t, func() {
		skinTexture, err := fetchTexture("http://textures.minecraft.net/texture/cd9ca55e9862f003ebfa1872a9244ad5f721d6b9e6883dd1d42f87dae127649")
		defer skinTexture.Close()

		So(err, ShouldBeNil)

		skin, err := DecodeSkin(skinTexture)

		So(err, ShouldBeNil)
		So(skin.Hash, ShouldEqual, "a04a26d10218668a632e419ab073cf57")
	})

	Convey("Bad texture request should gracefully fail", t, func() {
		skinTexture, err := fetchTexture("http://textures.minecraft.net/texture/")
		defer skinTexture.Close()

		So(err.Error(), ShouldStartWith, "Error retrieving texture")

		Convey("Bad texture decode should gracefully fail", func() {
			skin, err := DecodeSkin(skinTexture)

			So(err.Error(), ShouldContainSubstring, "image: unknown format")
			So(skin, ShouldResemble, Skin{})
		})

	})

}

func TestTexturesSteve(t *testing.T) {

	Convey("Steve should return valid image", t, func() {
		steveImg, err := FetchImageForSteve()

		So(err, ShouldBeNil)
		So(steveImg, ShouldNotBeNil)
	})

	Convey("Steve should return valid image", t, func() {
		steveSkin, err := FetchSkinForSteve()

		So(err, ShouldBeNil)
		So(steveSkin, ShouldNotResemble, Skin{})
		So(steveSkin.Hash, ShouldEqual, "98903c1609352e11552dca79eb1ce3d6")
	})

}

func TestTexturesSkin(t *testing.T) {

	Convey("d9135e082f2244c89cb0bee234155292 should return valid image from Mojang", t, func() {
		skin, err := FetchSkinFromMojangByUUID("d9135e082f2244c89cb0bee234155292")

		So(err, ShouldBeNil)
		So(skin, ShouldNotResemble, Skin{})
		So(skin.Hash, ShouldEqual, "a04a26d10218668a632e419ab073cf57")
	})

	Convey("clone1018 should return valid image from Mojang", t, func() {
		skin, err := FetchSkinFromMojang("clone1018")

		So(err, ShouldBeNil)
		So(skin, ShouldNotResemble, Skin{})
		So(skin.Hash, ShouldEqual, "a04a26d10218668a632e419ab073cf57")
	})

	Convey("clone1018 should return valid image from S3", t, func() {
		skin, err := FetchSkinFromS3("clone1018")

		So(err, ShouldBeNil)
		So(skin, ShouldNotResemble, Skin{})
		So(skin.Hash, ShouldEqual, "a04a26d10218668a632e419ab073cf57")
	})

	Convey("Wooxye should err from Mojang", t, func() {
		skin, err := FetchSkinFromMojang("Wooxye")

		So(err.Error(), ShouldStartWith, "Skin not found.")
		So(skin, ShouldResemble, Skin{Source: "Mojang"})
	})

	Convey("Wooxye should err from S3", t, func() {
		skin, err := FetchSkinFromS3("Wooxye")

		So(err.Error(), ShouldStartWith, "Skin not found.")
		So(skin, ShouldResemble, Skin{Source: "S3"})
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
