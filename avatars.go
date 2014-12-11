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
	// Skin image...
	Image image.Image
	// md5 hash of the skin image
	Hash string
	// Location we grabbed the skin from. Mojang/S3/Char
	Source string
	// 4-byte signature of the background matte for the skin
	AlphaSig [4]uint8
}

func GetSkin(u User) (Skin, error) {
	username := u.Name

	Skin, err := FetchSkinFromMojang(username)

	return Skin, err
}

func FetchSkinFromUrl(url, username string) (Skin, error) {
	resp, err := http.Get(url + username + ".png")
	if err != nil || resp.StatusCode != http.StatusOK {
		return Skin{}, errors.New("Skin not found. (" + fmt.Sprintf("%v", resp) + ")")
	}
	defer resp.Body.Close()

	return DecodeSkin(resp.Body)
}

func FetchSkinFromMojang(username string) (Skin, error) {
	url := "http://skins.minecraft.net/MinecraftSkins/"

	skin, err := FetchSkinFromUrl(url, username)
	skin.Source = "Mojang"

	return skin, err
}

func FetchSkinFromS3(username string) (Skin, error) {
	url := "http://s3.amazonaws.com/MinecraftSkins/"

	skin, err := FetchSkinFromUrl(url, username)
	skin.Source = "S3"

	return skin, err
}

func DecodeSkin(r io.Reader) (Skin, error) {
	// decode the image from the reader
	skinImg, _, err := image.Decode(r)
	if err != nil {
		return Skin{}, err
	}

	// Pull that into a nice png
	buf := new(bytes.Buffer)
	encErr := png.Encode(buf, skinImg)
	if encErr != nil {
		return Skin{}, encErr
	}

	// Create an md5 sum
	hasher := md5.New()
	hasher.Write(buf.Bytes())
	skinHash := fmt.Sprintf("%x", hasher.Sum(nil))

	// Finally, establish the skin
	skin := Skin{
		Image: skinImg,
		Hash:  skinHash,
	}
	// Create the alpha signature
	img := skin.Image.(*image.NRGBA)
	skin.AlphaSig = [...]uint8{
		img.Pix[0],
		img.Pix[1],
		img.Pix[2],
		img.Pix[3],
	}

	// And return the skin
	return skin, nil
}
