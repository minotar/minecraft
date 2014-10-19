// Minecraft Avatars
package minecraft

import (
	"errors"
	"fmt"
	"image"
	_ "image/png"
	"io"
	"net/http"
)

type Skin struct {
	Image image.Image
}

func GetSkin(u User) (Skin, error) {
	username := u.Name

	Skin, err := FetchSkinFromUrl(username)

	return Skin, err
}

func FetchSkinFromUrl(username string) (Skin, error) {
	url := "http://skins.minecraft.net/MinecraftSkins/"
	resp, err := http.Get(url + username + ".png")
	if err != nil || resp.StatusCode != http.StatusOK {
		return Skin{}, errors.New("Skin not found. (" + fmt.Sprintf("%v", resp.Header) + ")")
	}
	defer resp.Body.Close()

	return decodeSkin(resp.Body)
}

func decodeSkin(r io.Reader) (Skin, error) {
	skinImg, _, err := image.Decode(r)
	if err != nil {
		return Skin{}, err
	}
	return Skin{
		Image: skinImg,
	}, err
}
