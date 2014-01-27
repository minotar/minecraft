// Minecraft Avatars
package minecraft

import (
	"errors"
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

	Skin, err := fetchFromUrl(awsUrl(username))

	return Skin, err
}

func fetchFromUrl(url string) (Skin, error) {
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		return Skin{}, errors.New("Skin not found.")
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

func awsUrl(username string) string {
	return "http://skins.minotar.net/" + username + ".png"
}
