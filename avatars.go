// Minecraft Avatars
package minecraft

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/gographics/imagick/imagick"
	"io"
	"net/http"
)

type Skin struct {
	Image *imagick.MagickWand
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

	return DecodeSkin(resp.Body)
}

func DecodeSkin(r io.Reader) (Skin, error) {
	skinImg := imagick.NewMagickWand()
	skinImg.SetFormat("PNG")
	// Read all the bytes out of the reader item
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	blob := buf.Bytes()

	// And stick them in the skin
	err := skinImg.ReadImageBlob(blob)
	if err != nil {
		return Skin{}, err
	}

	hasher := md5.New()
	hasher.Write(blob)
	skinHash := fmt.Sprintf("%x", hasher.Sum(nil))

	return Skin{
		Image: skinImg,
		Hash:  skinHash,
	}, err
}
