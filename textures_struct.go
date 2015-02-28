package minecraft

import (
	"crypto/md5"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"net/http"
	// If we work with PNGs we need this
	_ "image/png"
	"io"
)

// Texture is our structure for the Cape/Skin structs and the functions for dealing with it
type Texture struct {
	// texture image...
	Image image.Image
	// md5 hash of the texture image
	Hash string
	// Location we grabbed the texture from. Mojang/S3/Char
	Source string
	// 4-byte signature of the background matte for the texture
	AlphaSig [4]uint8
	// URL of the texture
	URL string
}

func (t *Texture) fetch() error {
	if t.URL == "" {
		return errors.New("fetch failed: No Texture URL")
	}

	resp, err := http.Get(t.URL)
	if err != nil {
		return errors.New("fetch failed: Unable to Get URL - (" + err.Error() + ")")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("fetch failed: Error retrieving texture - (HTTP " + resp.Status + ")")
	}

	err = t.decode(resp.Body)
	if err != nil {
		return errors.New("fetch failed: (" + err.Error() + ")")
	}
	return nil
}

func (t *Texture) fetchWithSessionProfile(sessionProfile SessionProfileResponse, textureType string) error {
	profileTextureProperty, err := decodeTextureProperty(sessionProfile)
	if err != nil {
		return errors.New("fetchWithSessionProfile failed: (" + err.Error() + ")")
	}

	t.Source = "SessionProfile"

	url, err := decodeTextureURL(profileTextureProperty, textureType)
	if err != nil {
		return errors.New("fetchWithSessionProfile failed: (" + err.Error() + ")")
	}
	t.URL = url

	err = t.fetch()
	if err != nil {
		return errors.New("fetchWithSessionProfile failed: (" + err.Error() + ")")
	}
	return nil
}

// Includes the Skin not found detection that Mojang uses
func (t *Texture) fetchWithUsernameMojang(username string, textureType string) error {
	t.URL = "http://skins.minecraft.net/MinecraftSkins/" + username + ".png"
	if textureType != "Skin" {
		t.URL = "http://skins.minecraft.net/MinecraftCloaks/" + username + ".png"
	}

	t.Source = "Mojang"

	err := t.fetch()
	if err != nil {
		if err.Error() == "fetch failed: Error retrieving texture - (HTTP 404 Not Found)" {
			return errors.New("fetchWithUsernameMojang failed:  Texture not found - (" + err.Error() + ")")
		}
		return errors.New("fetchWithUsernameMojang failed:  (" + err.Error() + ")")
	}

	return nil
}

// Includes the Skin not found detection that S3 uses
func (t *Texture) fetchWithUsernameS3(username string, textureType string) error {
	t.URL = "http://s3.amazonaws.com/MinecraftSkins/" + username + ".png"
	if textureType != "Skin" {
		t.URL = "http://s3.amazonaws.com/MinecraftCloaks/" + username + ".png"
	}

	t.Source = "S3"

	err := t.fetch()
	if err != nil {
		if err.Error() == "fetch failed: Error retrieving texture - (HTTP 403 Forbidden)" {
			return errors.New("fetchWithUsernameS3 failed: Texture not found - (" + err.Error() + ")")
		}
		return errors.New("fetchWithUsernameS3 failed:  (" + err.Error() + ")")
	}

	return nil
}

// decode takes the image bytes and turns it into our Texture struct
func (t *Texture) decode(r io.Reader) error {
	err := t.castToNRGBA(r)
	if err != nil {
		return errors.New("decode failed: Error casting to NRGBA - (" + err.Error() + ")")
	}

	// And md5 hash its pixels
	hasher := md5.New()
	hasher.Write(t.Image.(*image.NRGBA).Pix)
	t.Hash = fmt.Sprintf("%x", hasher.Sum(nil))

	// Create the alpha signature
	img := t.Image.(*image.NRGBA)
	t.AlphaSig = [...]uint8{
		img.Pix[0],
		img.Pix[1],
		img.Pix[2],
		img.Pix[3],
	}

	return nil
}

func (t *Texture) castToNRGBA(r io.Reader) error {
	// Decode the skin
	textureImg, format, err := image.Decode(r)
	if err != nil {
		return errors.New("castToNRGBA failed: (" + err.Error() + ")")
	}

	// Convert it to NRGBA if necessary
	textureFinal := textureImg
	if format != "NRGBA" {
		bounds := textureImg.Bounds()
		textureFinal = image.NewNRGBA(bounds)
		draw.Draw(textureFinal.(draw.Image), bounds, textureImg, image.Pt(0, 0), draw.Src)
	}

	t.Image = textureFinal
	return nil
}
