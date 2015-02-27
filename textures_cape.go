package minecraft

import (
	"errors"
	_ "image/png"
)

type Cape struct {
	Texture
}

// FetchCape is a wrapper for the intelligence required to do things with UUIDs
// or Usernames. We'll do our best and fallback where appropriate
func FetchCape(user User) (*Cape, error) {
	return &Cape{}, errors.New("")
}

func FetchCapeFromMojangByUUID(uuid string) (*Cape, error) {
	capeTextureURL, err := decodeTextureURLWrapper(uuid, "Cape")
	if err != nil {
		return &Cape{}, err
	}

	capeTexture, err := fetchTexture(capeTextureURL)
	defer capeTexture.Close()
	if err != nil {
		return &Cape{}, err
	}

	cape := &Cape{}
	err = cape.decode(capeTexture)
	return cape, err
}
