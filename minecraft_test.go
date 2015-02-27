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

func TestProfiles(t *testing.T) {

	Convey("Test GetUUID", t, func() {

		Convey("clone1018 should match d9135e082f2244c89cb0bee234155292", func() {
			uuid, err := GetUUID("clone1018")

			So(err, ShouldBeNil)
			So(uuid, ShouldEqual, "d9135e082f2244c89cb0bee234155292")
		})

		Convey("skmkj88200aklk should gracefully error", func() {
			uuid, err := GetUUID("skmkj88200aklk ")

			So(err.Error(), ShouldStartWith, "User not found.")
			So(uuid, ShouldBeBlank)
		})

	})

	Convey("Test GetAPIProfile", t, func() {

		Convey("CLone1018 should equal clone1018", func() {
			apiProfile, err := GetAPIProfile("CLone1018")

			So(err, ShouldBeNil)
			So(apiProfile.Username, ShouldEqual, "clone1018")
		})

		Convey("skmkj88200aklk should gracefully error", func() {
			apiProfile, err := GetAPIProfile("skmkj88200aklk")

			So(err.Error(), ShouldStartWith, "User not found.")
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

			So(err.Error(), ShouldStartWith, "Error retrieving profile.")
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

			So(err.Error(), ShouldStartWith, "User not found.")
			So(playerUUID, ShouldBeBlank)
		})

		Convey("TooLongForAUsername should gracefully error", func() {
			playerUUID, err := NormalizePlayerForUUID("TooLongForAUsername")

			So(err.Error(), ShouldEqual, "Invalid Username or UUID.")
			So(playerUUID, ShouldBeBlank)
		})

	})

}

func TestTextures(t *testing.T) {

	Convey("Test decodeTextureProperty", t, func() {

		Convey("Should correctly decode Skin and Cape URL", func() {
			// citricsquid
			sessionProfileProperty := SessionProfileProperty{Name: "textures", Value: "eyJ0aW1lc3RhbXAiOjE0MjQ5ODM2MTI1NzgsInByb2ZpbGVJZCI6IjQ4YTBhN2U0ZDU1OTQ4NzNhNjE3ZGMxODlmNzZhOGExIiwicHJvZmlsZU5hbWUiOiJjaXRyaWNzcXVpZCIsInRleHR1cmVzIjp7IlNLSU4iOnsidXJsIjoiaHR0cDovL3RleHR1cmVzLm1pbmVjcmFmdC5uZXQvdGV4dHVyZS9lMWM2YzliNmRlODhmNDE4OGY5NzMyOTA5Yzc2ZGZjZDdiMTZhNDBhMDMxY2UxYjQ4NjhlNGQxZjg4OThlNGYifSwiQ0FQRSI6eyJ1cmwiOiJodHRwOi8vdGV4dHVyZXMubWluZWNyYWZ0Lm5ldC90ZXh0dXJlL2MzYWY3ZmI4MjEyNTQ2NjQ1NThmMjgzNjExNThjYTczMzAzYzlhODVlOTZlNTI1MTEwMjk1OGQ3ZWQ2MGM0YTMifX19=="}
			sessionProfile := SessionProfileResponse{Properties: []SessionProfileProperty{sessionProfileProperty}}

			profileTextureProperty, err := decodeTextureProperty(sessionProfile)

			So(err, ShouldBeNil)
			So(profileTextureProperty.Textures.Skin.URL, ShouldEqual, "http://textures.minecraft.net/texture/e1c6c9b6de88f4188f9732909c76dfcd7b16a40a031ce1b4868e4d1f8898e4f")
			So(profileTextureProperty.Textures.Cape.URL, ShouldEqual, "http://textures.minecraft.net/texture/c3af7fb821254664558f28361158ca73303c9a85e96e5251102958d7ed60c4a3")
		})

		Convey("Should only decode Skin URL", func() {
			// citricsquid
			sessionProfileProperty := SessionProfileProperty{Name: "textures", Value: "eyJ0aW1lc3RhbXAiOjE0MjQ5ODM2MTI1NzgsInByb2ZpbGVJZCI6IjQ4YTBhN2U0ZDU1OTQ4NzNhNjE3ZGMxODlmNzZhOGExIiwicHJvZmlsZU5hbWUiOiJjaXRyaWNzcXVpZCIsInRleHR1cmVzIjp7IlNLSU4iOnsidXJsIjoiaHR0cDovL3RleHR1cmVzLm1pbmVjcmFmdC5uZXQvdGV4dHVyZS9lMWM2YzliNmRlODhmNDE4OGY5NzMyOTA5Yzc2ZGZjZDdiMTZhNDBhMDMxY2UxYjQ4NjhlNGQxZjg4OThlNGYifX19"}
			sessionProfile := SessionProfileResponse{Properties: []SessionProfileProperty{sessionProfileProperty}}

			profileTextureProperty, err := decodeTextureProperty(sessionProfile)

			So(err, ShouldBeNil)
			So(profileTextureProperty.Textures.Skin.URL, ShouldEqual, "http://textures.minecraft.net/texture/e1c6c9b6de88f4188f9732909c76dfcd7b16a40a031ce1b4868e4d1f8898e4f")
			So(profileTextureProperty.Textures.Cape.URL, ShouldBeBlank)
		})

		Convey("Should only decode Cape URL", func() {
			// citricsquid
			sessionProfileProperty := SessionProfileProperty{Name: "textures", Value: "eyJ0aW1lc3RhbXAiOjE0MjQ5ODM2MTI1NzgsInByb2ZpbGVJZCI6IjQ4YTBhN2U0ZDU1OTQ4NzNhNjE3ZGMxODlmNzZhOGExIiwicHJvZmlsZU5hbWUiOiJjaXRyaWNzcXVpZCIsInRleHR1cmVzIjp7IkNBUEUiOnsidXJsIjoiaHR0cDovL3RleHR1cmVzLm1pbmVjcmFmdC5uZXQvdGV4dHVyZS9jM2FmN2ZiODIxMjU0NjY0NTU4ZjI4MzYxMTU4Y2E3MzMwM2M5YTg1ZTk2ZTUyNTExMDI5NThkN2VkNjBjNGEzIn19fQ=="}
			sessionProfile := SessionProfileResponse{Properties: []SessionProfileProperty{sessionProfileProperty}}

			profileTextureProperty, err := decodeTextureProperty(sessionProfile)

			So(err, ShouldBeNil)
			So(profileTextureProperty.Textures.Skin.URL, ShouldBeBlank)
			So(profileTextureProperty.Textures.Cape.URL, ShouldEqual, "http://textures.minecraft.net/texture/c3af7fb821254664558f28361158ca73303c9a85e96e5251102958d7ed60c4a3")
		})

		Convey("Should error about no textures", func() {
			sessionProfile := SessionProfileResponse{}

			profileTextureProperty, err := decodeTextureProperty(sessionProfile)

			So(err.Error(), ShouldEqual, "No textures property.")
			So(profileTextureProperty, ShouldResemble, SessionProfileTextureProperty{})
		})

		Convey("Should error trying to decode", func() {
			sessionProfileProperty := SessionProfileProperty{Name: "textures", Value: ""}
			sessionProfile := SessionProfileResponse{Properties: []SessionProfileProperty{sessionProfileProperty}}

			profileTextureProperty, err := decodeTextureProperty(sessionProfile)

			So(err.Error(), ShouldStartWith, "Error decoding texture property.")
			So(profileTextureProperty, ShouldResemble, SessionProfileTextureProperty{})
		})

	})

	// Must be careful to not request same profile from session server more than once per ~30 seconds
	Convey("Test decodeTextureURLWrapper", t, func() {

		Convey("48a0a7e4d5594873a617dc189f76a8a1 should return a Skin texture URL", func() {
			// citricsquid
			capeTextureURL, err := decodeTextureURLWrapper("48a0a7e4d5594873a617dc189f76a8a1", "Skin")

			So(err, ShouldBeNil)
			So(capeTextureURL, ShouldEqual, "http://textures.minecraft.net/texture/e1c6c9b6de88f4188f9732909c76dfcd7b16a40a031ce1b4868e4d1f8898e4f")
		})

		Convey("069a79f444e94726a5befca90e38aaf5 should return a Cape texture URL", func() {
			// Notch
			capeTextureURL, err := decodeTextureURLWrapper("069a79f444e94726a5befca90e38aaf5", "Cape")

			So(err, ShouldBeNil)
			So(capeTextureURL, ShouldEqual, "http://textures.minecraft.net/texture/3f688e0e699b3d9fe448b5bb50a3a288f9c589762b3dae8308842122dcb81")
		})

		Convey("Cape request for 2f3665cc5e29439bbd14cb6d3a6313a7 should gracefully error", func() {
			// lukegb
			capeTextureURL, err := decodeTextureURLWrapper("2f3665cc5e29439bbd14cb6d3a6313a7", "Cape")

			So(err.Error(), ShouldEqual, "Cape URL is not present.")
			So(capeTextureURL, ShouldBeBlank)
		})

	})

	Convey("Test fetchTexture", t, func() {

		Convey("clone1018 texture should return the correct skin", func() {
			skinTexture, err := fetchTexture("http://textures.minecraft.net/texture/cd9ca55e9862f003ebfa1872a9244ad5f721d6b9e6883dd1d42f87dae127649")
			defer skinTexture.Close()

			So(err, ShouldBeNil)

			skin := &Skin{}
			err = skin.decode(skinTexture)

			So(err, ShouldBeNil)
			So(skin.Hash, ShouldEqual, "a04a26d10218668a632e419ab073cf57")
		})

		Convey("Bad texture request should gracefully fail", func() {
			skinTexture, err := fetchTexture("http://textures.minecraft.net/texture/")
			defer skinTexture.Close()

			So(err.Error(), ShouldStartWith, "Error retrieving texture.")

			Convey("Bad texture decode should gracefully fail", func() {
				skin := &Skin{}
				err = skin.decode(skinTexture)

				So(err.Error(), ShouldContainSubstring, "image: unknown format")
				So(skin, ShouldResemble, &Skin{})
			})

		})

	})

	Convey("Test Steve", t, func() {

		Convey("Steve should return valid image", func() {
			steveImg, err := FetchImageForSteve()

			So(err, ShouldBeNil)
			So(steveImg, ShouldNotBeNil)
		})

		Convey("Steve should return valid skin", func() {
			steveSkin, err := FetchSkinForSteve()

			So(err, ShouldBeNil)
			So(steveSkin, ShouldNotResemble, &Skin{})
			So(steveSkin.Hash, ShouldEqual, "98903c1609352e11552dca79eb1ce3d6")
		})

	})

	// Must be careful to not request same profile from session server more than once per ~30 seconds
	Convey("Test Capes", t, func() {

		Convey("61699b2ed3274a019f1e0ea8c3f06bc6 should return a Cape from Mojang", func() {
			// Dinnerbone
			skin, err := FetchCapeFromMojangByUUID("61699b2ed3274a019f1e0ea8c3f06bc6")

			So(err, ShouldBeNil)
			So(skin, ShouldNotResemble, &Cape{})
			So(skin.Hash, ShouldEqual, "e45b09f09f971dd92fa51c550ca21876")
		})

		Convey("2aa9ff75db7140caa23189e693ad7d79 should err from Mojang", func() {
			// samuel
			skin, err := FetchCapeFromMojangByUUID("2aa9ff75db7140caa23189e693ad7d79")

			So(err.Error(), ShouldEqual, "Cape URL is not present.")
			So(skin, ShouldResemble, &Cape{})
			So(skin.Hash, ShouldBeBlank)
		})

	})

	Convey("Test Skins", t, func() {

		// Must be careful to not request same profile from session server more than once per ~30 seconds
		Convey("d9135e082f2244c89cb0bee234155292 should return valid image from Mojang", func() {
			// clone1018
			skin, err := FetchSkinFromMojangByUUID("d9135e082f2244c89cb0bee234155292")

			So(err, ShouldBeNil)
			So(skin, ShouldNotResemble, &Skin{})
			So(skin.Hash, ShouldEqual, "a04a26d10218668a632e419ab073cf57")
		})

		Convey("clone1018 should return valid image from Mojang", func() {
			skin, err := FetchSkinFromMojang("clone1018")

			So(err, ShouldBeNil)
			So(skin, ShouldNotResemble, &Skin{})
			So(skin.Hash, ShouldEqual, "a04a26d10218668a632e419ab073cf57")
		})

		Convey("clone1018 should return valid image from S3", func() {
			skin, err := FetchSkinFromS3("clone1018")

			So(err, ShouldBeNil)
			So(skin, ShouldNotResemble, &Skin{})
			So(skin.Hash, ShouldEqual, "a04a26d10218668a632e419ab073cf57")
		})

		Convey("Wooxye should err from Mojang", func() {
			skin, err := FetchSkinFromMojang("Wooxye")

			So(err.Error(), ShouldStartWith, "Skin not found.")
			So(skin, ShouldResemble, &Skin{Texture{Source: "Mojang", URL: "http://skins.minecraft.net/MinecraftSkins/Wooxye.png"}})
		})

		Convey("Wooxye should err from S3", func() {
			skin, err := FetchSkinFromS3("Wooxye")

			So(err.Error(), ShouldStartWith, "Skin not found.")
			So(skin, ShouldResemble, &Skin{Texture{Source: "S3", URL: "http://s3.amazonaws.com/MinecraftSkins/Wooxye.png"}})
		})

	})

}
