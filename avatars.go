package minecraft

import (
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
	if err != nil {
		return Skin{}, err
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
	return "http://s3.amazonaws.com/MinecraftSkins/" + username + ".png"
}
