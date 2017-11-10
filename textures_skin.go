package minecraft

import _ "image/png" // If we work with PNGs we need this

type Skin struct {
	Texture
}

func FetchSkinUUID(uuid string) (Skin, error) {
	skin := &Skin{}

	// Must be careful to not request same profile from session server more than once per ~30 seconds
	sessionProfile, err := GetSessionProfile(uuid)
	if err != nil {
		return *skin, err
	}

	return *skin, skin.FetchWithSessionProfile(sessionProfile, "Skin")
}

func FetchSkinUsernameMojang(username string) (Skin, error) {
	skin := &Skin{}

	return *skin, skin.FetchWithUsernameMojang(username, "Skin")
}

func FetchSkinUsernameS3(username string) (Skin, error) {
	skin := &Skin{}

	return *skin, skin.FetchWithUsernameS3(username, "Skin")
}
