package minecraft

import (
	"errors"
	_ "image/png"
)

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

	err := skin.FetchWithUsernameMojang(username, "Skin")
	if err != nil {
		return *skin, err
	}

	// Some proper testing is really required to determine if this is truly needed.
	if skin.Hash == SteveHash {
		return Skin{Texture{Source: "Steve"}}, errors.New("FetchSkinUsernameMojang failed: Rate limited")
	}

	return *skin, nil
}

func FetchSkinUsernameS3(username string) (Skin, error) {
	skin := &Skin{}

	return *skin, skin.FetchWithUsernameS3(username, "Skin")
}
