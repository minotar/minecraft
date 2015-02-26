package minecraft

import _ "image/png"

type Cape struct {
	Texture
}

func FetchCapeFromMojangByUUID(uuid string) (*Cape, error) {
	capeTextureURL, err := decodeTextureURL(uuid, "Cape")
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
