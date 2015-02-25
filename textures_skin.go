package minecraft

import (
	"errors"
	"fmt"
	"image"
	_ "image/png"
	"net/http"
)

type Skin struct {
	// Skin image...
	Image image.Image
	// md5 hash of the skin image
	Hash string
	// Location we grabbed the skin from. Mojang/S3/Char
	Source string
	// 4-byte signature of the background matte for the skin
	AlphaSig [4]uint8
}

func FetchSkinFromUrl(url, username string) (Skin, error) {
	resp, err := http.Get(url + username + ".png")
	if err != nil {
		return Skin{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Skin{}, errors.New("Skin not found. (" + fmt.Sprintf("%v", resp) + ")")
	}

	return DecodeSkin(resp.Body)
}

func FetchSkinFromMojang(username string) (Skin, error) {
	url := "http://skins.minecraft.net/MinecraftSkins/"

	skin, err := FetchSkinFromUrl(url, username)
	skin.Source = "Mojang"

	if skin.Hash == SteveHash {
		return Skin{}, errors.New("Rate limited")
	}

	return skin, err
}

func FetchSkinFromS3(username string) (Skin, error) {
	url := "http://s3.amazonaws.com/MinecraftSkins/"

	skin, err := FetchSkinFromUrl(url, username)
	skin.Source = "S3"

	return skin, err
}

func FetchSkinFromMojangByUUID(uuid string) (Skin, error) {
	skinTextureURL, err := decodeTextureURL(uuid, "Skin")
	if err != nil {
		return Skin{}, err
	}

	skinTexture, err := fetchTexture(skinTextureURL)
	defer skinTexture.Close()
	if err != nil {
		return Skin{}, err
	}

	return DecodeSkin(skinTexture)
}
