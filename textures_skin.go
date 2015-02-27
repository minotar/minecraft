package minecraft

import (
	"errors"
	_ "image/png"
)

type Skin struct {
	Texture
}

func FetchSkinFromMojang(username string) (*Skin, error) {
	skin := &Skin{Texture{Source: "Mojang", URL: "http://skins.minecraft.net/MinecraftSkins/" + username + ".png"}}

	err := skin.fetch()
	if err != nil {
		if err.Error() == "Error retrieving texture. (HTTP 404 Not Found)" {
			return skin, errors.New("Skin not found. " + err.Error())
		}
		return skin, err
	}

	if skin.Hash == SteveHash {
		return &Skin{}, errors.New("Rate limited")
	}

	return skin, nil
}

func FetchSkinFromS3(username string) (*Skin, error) {
	skin := &Skin{Texture{Source: "S3", URL: "http://s3.amazonaws.com/MinecraftSkins/" + username + ".png"}}

	err := skin.fetch()
	if err != nil {
		if err.Error() == "Error retrieving texture. (HTTP 403 Forbidden)" {
			return skin, errors.New("Skin not found. " + err.Error())
		}
		return skin, err
	}

	return skin, nil
}

func FetchSkinFromMojangByUUID(uuid string) (*Skin, error) {
	skin := &Skin{Texture{Source: "Mojang"}}

	var err error
	skin.URL, err = decodeTextureURLWrapper(uuid, "Skin")
	if err != nil {
		return skin, err
	}

	return skin, skin.fetch()
}
