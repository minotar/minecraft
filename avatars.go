// Minecraft Avatars
package minecraft

import (
	"bytes"
	"crypto/md5"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"io"
	"io/ioutil"
	"net/http"
)

const CharHash = "613af1b0b41e4deae34e272f3487fba6"

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

	if skin.Hash == CharHash {
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

func DecodeSkin(r io.Reader) (Skin, error) {
	// decode the image from the reader
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return Skin{}, err
	}

	hashBuf := make(chan string)
	skinBuf := make(chan image.Image)

	// Decode the skin we go
	go func() {
		var s image.Image

		skinImg, format, err := image.Decode(bytes.NewReader(buf))
		if err != nil {
			chr, _ := FetchImageForChar()
			s = chr
		} else {
			s = skinImg
			format = ""
		}
		// Convert it to NRGBA if necessary
		out := s
		if format != "NRGBA" {
			bounds := s.Bounds()
			out = image.NewNRGBA(bounds)
			draw.Draw(out.(draw.Image), bounds, s, image.Pt(0, 0), draw.Src)
		}

		// Send it back down the channel
		skinBuf <- out
	}()

	// And md5 hash it
	go func() {
		hasher := md5.New()
		hasher.Write(buf)
		hashBuf <- fmt.Sprintf("%x", hasher.Sum(nil))
	}()

	// Create an md5 sum
	// Finally, establish the skin
	skin := Skin{
		Image: <-skinBuf,
		Hash:  <-hashBuf,
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
