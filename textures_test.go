// textures_test.go
package minecraft

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTextures(t *testing.T) {

	Convey("Test DecodeTextureProperty", t, func() {

		testServer := startTestServer(returnMux())
		defer closeTestServer(testServer)

		Convey("Should correctly decode Skin and Cape URL", func() {
			sessionProfile, _ := GetSessionProfile("48a0a7e4d5594873a617dc189f76a8a1")
			profileTextureProperty, err := DecodeTextureProperty(sessionProfile)

			So(err, ShouldBeNil)
			So(profileTextureProperty.Textures.Skin.URL, ShouldEqual, "http://textures.minecraft.net/texture/e1c6c9b6de88f4188f9732909c76dfcd7b16a40a031ce1b4868e4d1f8898e4f")
			So(profileTextureProperty.Textures.Cape.URL, ShouldEqual, "http://textures.minecraft.net/texture/c3af7fb821254664558f28361158ca73303c9a85e96e5251102958d7ed60c4a3")
		})

		Convey("Should only decode Skin URL", func() {
			sessionProfile, _ := GetSessionProfile("d9135e082f2244c89cb0bee234155292")
			profileTextureProperty, err := DecodeTextureProperty(sessionProfile)

			So(err, ShouldBeNil)
			So(profileTextureProperty.Textures.Skin.URL, ShouldEqual, "http://textures.minecraft.net/texture/cd9ca55e9862f003ebfa1872a9244ad5f721d6b9e6883dd1d42f87dae127649")
			So(profileTextureProperty.Textures.Cape.URL, ShouldBeBlank)
		})

		Convey("Should error about no textures", func() {
			sessionProfile, _ := GetSessionProfile("00000000000000000000000000000004")
			profileTextureProperty, err := DecodeTextureProperty(sessionProfile)

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "DecodeTextureProperty failed: No textures property.")
			So(profileTextureProperty, ShouldResemble, SessionProfileTextureProperty{})
		})

		Convey("Should error trying to decode", func() {
			sessionProfile, _ := GetSessionProfile("00000000000000000000000000000005")
			profileTextureProperty, err := DecodeTextureProperty(sessionProfile)

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "DecodeTextureProperty failed: Error decoding texture property - (unexpected EOF)")
			So(profileTextureProperty, ShouldResemble, SessionProfileTextureProperty{})
		})

	})

	Convey("Test Texture.fetch", t, func() {

		testServer := startTestServer(returnMux())
		defer closeTestServer(testServer)

		Convey("clone1018 texture should return the correct skin", func() {
			texture := &Texture{URL: "http://textures.minecraft.net/texture/cd9ca55e9862f003ebfa1872a9244ad5f721d6b9e6883dd1d42f87dae127649"}

			err := texture.Fetch()

			So(err, ShouldBeNil)
			So(texture.Hash, ShouldEqual, "a04a26d10218668a632e419ab073cf57")
		})

		Convey("Bad texture requests should gracefully fail", func() {

			Convey("No texture URL", func() {
				texture := &Texture{URL: ""}

				err := texture.Fetch()

				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Fetch failed: No Texture URL")
			})

			Convey("Bad texture URL (non-image)", func() {
				texture := &Texture{URL: "http://testServer/200"}

				err := texture.Fetch()

				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Fetch failed: (Decode failed: Error casting to NRGBA - (CastToNRGBA failed: (image: unknown format)))")
			})

			Convey("Bad texture URL (non-200)", func() {
				texture := &Texture{URL: "http://testServer/404"}

				err := texture.Fetch()

				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Fetch failed: Error retrieving texture - (HTTP 404 Not Found)")
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
			So(steveSkin, ShouldNotResemble, Skin{})
			So(steveSkin.Hash, ShouldEqual, "98903c1609352e11552dca79eb1ce3d6")
		})

	})

	Convey("Test Skins", t, func() {

		testServer := startTestServer(returnMux())
		defer closeTestServer(testServer)

		Convey("d9135e082f2244c89cb0bee234155292 should return valid image from Mojang", func() {
			// clone1018
			skin, err := FetchSkinUUID("d9135e082f2244c89cb0bee234155292")

			So(err, ShouldBeNil)
			So(skin, ShouldNotResemble, Skin{})
			So(skin.Hash, ShouldEqual, "a04a26d10218668a632e419ab073cf57")
		})

		Convey("00000000000000000000000000000000 err from Mojang", func() {
			skin, err := FetchSkinUUID("10000000000000000000000000000000")

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "GetSessionProfile failed: (apiRequest failed: User not found - (HTTP 204 No Content))")
			So(skin, ShouldResemble, Skin{})
		})

		Convey("clone1018 should return valid image from Mojang", func() {
			skin, err := FetchSkinUsernameMojang("clone1018")

			So(err, ShouldBeNil)
			So(skin, ShouldNotResemble, Skin{})
			So(skin.Hash, ShouldEqual, "a04a26d10218668a632e419ab073cf57")
		})

		Convey("Wooxye should err from Mojang", func() {
			skin, err := FetchSkinUsernameMojang("Wooxye")

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "FetchWithUsernameMojang failed:  Texture not found - (Fetch failed: Error retrieving texture - (HTTP 404 Not Found))")
			So(skin, ShouldResemble, Skin{Texture{Source: "Mojang", URL: "http://skins.minecraft.net/MinecraftSkins/Wooxye.png"}})
		})

		Convey("clone1018 should return valid image from S3", func() {
			skin, err := FetchSkinUsernameS3("clone1018")

			So(err, ShouldBeNil)
			So(skin, ShouldNotResemble, Skin{})
			So(skin.Hash, ShouldEqual, "a04a26d10218668a632e419ab073cf57")
		})

		Convey("Wooxye should err from S3", func() {
			skin, err := FetchSkinUsernameS3("Wooxye")

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "FetchWithUsernameS3 failed: Texture not found - (Fetch failed: Error retrieving texture - (HTTP 403 Forbidden))")
			So(skin, ShouldResemble, Skin{Texture{Source: "S3", URL: "http://s3.amazonaws.com/MinecraftSkins/Wooxye.png"}})
		})

	})

	Convey("Test Capes", t, func() {

		testServer := startTestServer(returnMux())
		defer closeTestServer(testServer)

		Convey("48a0a7e4d5594873a617dc189f76a8a1 should return a Cape from Mojang", func() {
			// citricsquid
			cape, err := FetchCapeUUID("48a0a7e4d5594873a617dc189f76a8a1")

			So(err, ShouldBeNil)
			So(cape, ShouldNotResemble, Cape{})
			So(cape.Hash, ShouldEqual, "8cbf8786caba2f05383cf887be592ee6")
		})

		Convey("2f3665cc5e29439bbd14cb6d3a6313a7 should err from Mojang (No Cape)", func() {
			// lukegb
			cape, err := FetchCapeUUID("2f3665cc5e29439bbd14cb6d3a6313a7")

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "FetchWithSessionProfile failed: (FetchWithTextureProperty failed: (DecodeTextureURL failed: Cape URL is not present.))")
			So(cape, ShouldResemble, Cape{})
			So(cape.Hash, ShouldBeBlank)
		})

		Convey("10000000000000000000000000000000 should err from Mojang (No User)", func() {
			cape, err := FetchCapeUUID("10000000000000000000000000000000")

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "GetSessionProfile failed: (apiRequest failed: User not found - (HTTP 204 No Content))")
			So(cape, ShouldResemble, Cape{})
		})

		Convey("citricsquid should return a Cape from Mojang", func() {
			cape, err := FetchCapeUsernameMojang("citricsquid")

			So(err, ShouldBeNil)
			So(cape, ShouldNotResemble, Cape{})
			So(cape.Hash, ShouldEqual, "8cbf8786caba2f05383cf887be592ee6")
		})

		Convey("Wooxye should err from Mojang", func() {
			cape, err := FetchCapeUsernameMojang("Wooxye")

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "FetchWithUsernameMojang failed:  Texture not found - (Fetch failed: Error retrieving texture - (HTTP 404 Not Found))")
			So(cape, ShouldResemble, Cape{Texture{Source: "Mojang", URL: "http://skins.minecraft.net/MinecraftCloaks/Wooxye.png"}})
		})

		Convey("citricsquid should return a Cape from S3", func() {
			cape, err := FetchCapeUsernameS3("citricsquid")

			So(err, ShouldBeNil)
			So(cape, ShouldNotResemble, Cape{})
			So(cape.Hash, ShouldEqual, "8cbf8786caba2f05383cf887be592ee6")
		})

		Convey("Wooxye should err from S3", func() {
			cape, err := FetchCapeUsernameS3("Wooxye")

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "FetchWithUsernameS3 failed: Texture not found - (Fetch failed: Error retrieving texture - (HTTP 403 Forbidden))")
			So(cape, ShouldResemble, Cape{Texture{Source: "S3", URL: "http://s3.amazonaws.com/MinecraftCloaks/Wooxye.png"}})
		})

	})

	// This could be a lot more DRY but shush
	Convey("Test FetchTextures", t, func() {

		testServer := startTestServer(returnMux())
		defer closeTestServer(testServer)

		Convey("clone1018", func() {
			user, skin, cape, err := FetchTextures("clone1018")

			So(err, ShouldBeNil)
			So(user.UUID, ShouldEqual, "d9135e082f2244c89cb0bee234155292")
			So(cape, ShouldResemble, Cape{})
			So(skin.Hash, ShouldEqual, "a04a26d10218668a632e419ab073cf57")
		})

		Convey("d9135e082f2244c89cb0bee234155292", func() {
			user, skin, cape, err := FetchTextures("d9135e082f2244c89cb0bee234155292")

			So(err, ShouldBeNil)
			So(user.Username, ShouldEqual, "clone1018")
			So(cape, ShouldResemble, Cape{})
			So(skin.Hash, ShouldEqual, "a04a26d10218668a632e419ab073cf57")
		})

		Convey("citricsquid", func() {
			user, skin, cape, err := FetchTextures("citricsquid")

			So(err, ShouldBeNil)
			So(user.UUID, ShouldEqual, "48a0a7e4d5594873a617dc189f76a8a1")
			So(cape.Hash, ShouldEqual, "8cbf8786caba2f05383cf887be592ee6")
			So(skin.Hash, ShouldEqual, "c05454f331fa93b3e38866a9ec52c467")
		})

		Convey("48a0a7e4d5594873a617dc189f76a8a1", func() {
			user, skin, cape, err := FetchTextures("48a0a7e4d5594873a617dc189f76a8a1")

			So(err, ShouldBeNil)
			So(user.Username, ShouldEqual, "citricsquid")
			So(cape.Hash, ShouldEqual, "8cbf8786caba2f05383cf887be592ee6")
			So(skin.Hash, ShouldEqual, "c05454f331fa93b3e38866a9ec52c467")
		})

		Convey("RateLimitAPI", func() {
			user, skin, cape, err := FetchTextures("RateLimitAPI")

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "FetchTextures fallback to Steve: (GetAPIProfile failed: (apiRequest failed: Rate limited))")
			So(user, ShouldResemble, User{})
			So(cape, ShouldResemble, Cape{})
			So(skin.Source, ShouldEqual, "Steve")
			So(skin.Hash, ShouldEqual, "98903c1609352e11552dca79eb1ce3d6")
		})

		Convey("RateLimitSession", func() {
			user, skin, cape, err := FetchTextures("RateLimitSession")

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "FetchTextures fallback to Steve: (GetSessionProfile failed: (apiRequest failed: Rate limited))")
			So(user, ShouldResemble, User{})
			So(cape, ShouldResemble, Cape{})
			So(skin.Source, ShouldEqual, "Steve")
			So(skin.Hash, ShouldEqual, "98903c1609352e11552dca79eb1ce3d6")
		})

		Convey("MalformedAPI", func() {
			user, skin, cape, err := FetchTextures("MalformedAPI")

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "FetchTextures fallback to Steve: (GetAPIProfile failed: Error decoding profile - (unexpected EOF))")
			So(user, ShouldResemble, User{})
			So(cape, ShouldResemble, Cape{})
			So(skin.Source, ShouldEqual, "Steve")
			So(skin.Hash, ShouldEqual, "98903c1609352e11552dca79eb1ce3d6")
		})

		Convey("MalformedSession", func() {
			user, skin, cape, err := FetchTextures("MalformedSession")

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "FetchTextures fallback to Steve: (GetSessionProfile failed: Error decoding profile - (unexpected EOF))")
			So(user, ShouldResemble, User{})
			So(cape, ShouldResemble, Cape{})
			So(skin.Source, ShouldEqual, "Steve")
			So(skin.Hash, ShouldEqual, "98903c1609352e11552dca79eb1ce3d6")
		})

		Convey("NoTexture", func() {
			user, skin, cape, err := FetchTextures("NoTexture")

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "FetchTextures fallback to Steve: (FetchTexturesWithSessionProfile failed: Unable to decode sessionProfile (DecodeTextureProperty failed: No textures property.))")
			So(user.UUID, ShouldEqual, "00000000000000000000000000000004")
			So(cape, ShouldResemble, Cape{})
			So(skin.Source, ShouldEqual, "Steve")
			So(skin.Hash, ShouldEqual, "98903c1609352e11552dca79eb1ce3d6")
		})

		Convey("MalformedTexProp", func() {
			user, skin, cape, err := FetchTextures("MalformedTexProp")

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "FetchTextures fallback to Steve: (FetchTexturesWithSessionProfile failed: Unable to decode sessionProfile (DecodeTextureProperty failed: Error decoding texture property - (unexpected EOF)))")
			So(user.UUID, ShouldEqual, "00000000000000000000000000000005")
			So(cape, ShouldResemble, Cape{})
			So(skin.Source, ShouldEqual, "Steve")
			So(skin.Hash, ShouldEqual, "98903c1609352e11552dca79eb1ce3d6")
		})

		Convey("500API", func() {
			user, skin, cape, err := FetchTextures("500API")

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "FetchTextures fallback to Steve: (GetAPIProfile failed: (apiRequest failed: Error retrieving profile - (HTTP 500 Internal Server Error)))")
			So(user, ShouldResemble, User{})
			So(cape, ShouldResemble, Cape{})
			So(skin.Source, ShouldEqual, "Steve")
			So(skin.Hash, ShouldEqual, "98903c1609352e11552dca79eb1ce3d6")
		})

		Convey("500Session", func() {
			user, skin, cape, err := FetchTextures("500Session")

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "FetchTextures fallback to Steve: (GetSessionProfile failed: (apiRequest failed: Error retrieving profile - (HTTP 500 Internal Server Error)))")
			So(user, ShouldResemble, User{})
			So(cape, ShouldResemble, Cape{})
			So(skin.Source, ShouldEqual, "Steve")
			So(skin.Hash, ShouldEqual, "98903c1609352e11552dca79eb1ce3d6")
		})

		Convey("MalformedSTex", func() {
			user, skin, cape, err := FetchTextures("MalformedSTex")

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "FetchTextures fallback to Steve: (FetchTexturesWithSessionProfile failed: Unable to retrieve skin - (FetchWithTextureProperty failed: (Fetch failed: (Decode failed: Error casting to NRGBA - (CastToNRGBA failed: (unexpected EOF))))))")
			So(user.UUID, ShouldEqual, "00000000000000000000000000000008")
			So(cape, ShouldResemble, Cape{})
			So(skin.Source, ShouldEqual, "Steve")
			So(skin.Hash, ShouldEqual, "98903c1609352e11552dca79eb1ce3d6")
		})

		Convey("MalformedCTex", func() {
			user, skin, cape, err := FetchTextures("MalformedCTex")

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "FetchTextures unable to get the cape: (FetchTexturesWithSessionProfile failed: Unable to retrieve cape - (FetchWithTextureProperty failed: (Fetch failed: (Decode failed: Error casting to NRGBA - (CastToNRGBA failed: (unexpected EOF))))))")
			So(user.UUID, ShouldEqual, "00000000000000000000000000000009")
			So(cape, ShouldResemble, Cape{Texture{Source: "SessionProfile", URL: "http://textures.minecraft.net/texture/MalformedTexture"}})
			So(skin.Source, ShouldEqual, "SessionProfile")
			So(skin.Hash, ShouldEqual, "a04a26d10218668a632e419ab073cf57")
		})

		Convey("404STexture", func() {
			user, skin, cape, err := FetchTextures("404STexture")

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "FetchTextures fallback to Steve: (FetchTexturesWithSessionProfile failed: Unable to retrieve skin - (FetchWithTextureProperty failed: (Fetch failed: Error retrieving texture - (HTTP 404 Not Found))))")
			So(user.UUID, ShouldEqual, "00000000000000000000000000000010")
			So(cape, ShouldResemble, Cape{})
			So(skin.Source, ShouldEqual, "Steve")
			So(skin.Hash, ShouldEqual, "98903c1609352e11552dca79eb1ce3d6")
		})

		Convey("404CTexture", func() {
			user, skin, cape, err := FetchTextures("404CTexture")

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "FetchTextures unable to get the cape: (FetchTexturesWithSessionProfile failed: Unable to retrieve cape - (FetchWithTextureProperty failed: (Fetch failed: Error retrieving texture - (HTTP 404 Not Found))))")
			So(user.UUID, ShouldEqual, "00000000000000000000000000000011")
			So(cape, ShouldResemble, Cape{Texture{Source: "SessionProfile", URL: "http://textures.minecraft.net/texture/404Texture"}})
			So(skin.Source, ShouldEqual, "SessionProfile")
			So(skin.Hash, ShouldEqual, "a04a26d10218668a632e419ab073cf57")
		})

		Convey("RLSessionMojang", func() {
			user, skin, cape, err := FetchTextures("RLSessionMojang")

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "FetchTextures fallback to UsernameMojang: (GetSessionProfile failed: (apiRequest failed: Rate limited))")
			So(user, ShouldResemble, User{UUID: "", Username: "RLSessionMojang"})
			So(cape, ShouldResemble, Cape{})
			So(skin.Source, ShouldEqual, "Mojang")
			So(skin.Hash, ShouldEqual, "a04a26d10218668a632e419ab073cf57")
		})

		Convey("RLSessionS3", func() {
			user, skin, cape, err := FetchTextures("RLSessionS3")

			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "FetchTextures fallback to UsernameS3: (GetSessionProfile failed: (apiRequest failed: Rate limited))")
			So(user, ShouldResemble, User{UUID: "", Username: "RLSessionS3"})
			So(cape, ShouldResemble, Cape{})
			So(skin.Source, ShouldEqual, "S3")
			So(skin.Hash, ShouldEqual, "c05454f331fa93b3e38866a9ec52c467")
		})

	})

}
