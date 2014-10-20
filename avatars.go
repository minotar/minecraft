// Minecraft Avatars
package minecraft

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"net/http"
)

type Skin struct {
	Image image.Image
	Hash  string
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
		return Skin{}, errors.New("Skin not found. (" + fmt.Sprintf("%v", resp) + ")")
	}
	defer resp.Body.Close()

	return decodeSkin(resp.Body)
}

func decodeSkin(r io.Reader) (Skin, error) {
	skinImg, _, err := image.Decode(r)
	if err != nil {
		return Skin{}, err
	}

	buf := new(bytes.Buffer)
	encErr := png.Encode(buf, skinImg)
	if encErr != nil {
		return Skin{}, encErr
	}
	hasher := md5.New()
	hasher.Write(buf.Bytes())
	skinHash := fmt.Sprintf("%x", hasher.Sum(nil))

	return Skin{
		Image: skinImg,
		Hash:  skinHash,
	}, err
}
