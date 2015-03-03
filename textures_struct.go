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

func (t *Texture) Fetch() error {
	if t.URL == "" {
		return errors.New("Fetch failed: No Texture URL")
	}

	resp, err := http.Get(t.URL)
	if err != nil {
		return errors.New("Fetch failed: Unable to Get URL - (" + err.Error() + ")")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("Fetch failed: Error retrieving texture - (HTTP " + resp.Status + ")")
	}

	err = t.Decode(resp.Body)
	if err != nil {
		return errors.New("Fetch failed: (" + err.Error() + ")")
	}
	return nil
}

func (t *Texture) FetchWithSessionProfile(sessionProfile SessionProfileResponse, textureType string) error {
	profileTextureProperty, err := DecodeTextureProperty(sessionProfile)
	if err != nil {
		return errors.New("FetchWithSessionProfile failed: (" + err.Error() + ")")
	}

	err = t.FetchWithTextureProperty(profileTextureProperty, textureType)
	if err != nil {
		return errors.New("FetchWithSessionProfile failed: (" + err.Error() + ")")
	}
	return nil
}

func (t *Texture) FetchWithTextureProperty(profileTextureProperty SessionProfileTextureProperty, textureType string) error {
	url, err := DecodeTextureURL(profileTextureProperty, textureType)
	if err != nil {
		return errors.New("FetchWithTextureProperty failed: (" + err.Error() + ")")
	}
	t.Source = "SessionProfile"
	t.URL = url

	err = t.Fetch()
	if err != nil {
		return errors.New("FetchWithTextureProperty failed: (" + err.Error() + ")")
	}
	return nil
}

// Includes the Skin not found detection that Mojang uses
func (t *Texture) FetchWithUsernameMojang(username string, textureType string) error {
	t.URL = "http://skins.minecraft.net/MinecraftSkins/" + username + ".png"
	if textureType != "Skin" {
		t.URL = "http://skins.minecraft.net/MinecraftCloaks/" + username + ".png"
	}

	t.Source = "Mojang"

	err := t.Fetch()
	if err != nil {
		if err.Error() == "Fetch failed: Error retrieving texture - (HTTP 404 Not Found)" {
			return errors.New("FetchWithUsernameMojang failed:  Texture not found - (" + err.Error() + ")")
		}
		return errors.New("FetchWithUsernameMojang failed:  (" + err.Error() + ")")
	}

	return nil
}

// Includes the Skin not found detection that S3 uses
func (t *Texture) FetchWithUsernameS3(username string, textureType string) error {
	t.URL = "http://s3.amazonaws.com/MinecraftSkins/" + username + ".png"
	if textureType != "Skin" {
		t.URL = "http://s3.amazonaws.com/MinecraftCloaks/" + username + ".png"
	}

	t.Source = "S3"

	err := t.Fetch()
	if err != nil {
		if err.Error() == "Fetch failed: Error retrieving texture - (HTTP 403 Forbidden)" {
			return errors.New("FetchWithUsernameS3 failed: Texture not found - (" + err.Error() + ")")
		}
		return errors.New("FetchWithUsernameS3 failed:  (" + err.Error() + ")")
	}

	return nil
}

// decode takes the image bytes and turns it into our Texture struct
func (t *Texture) Decode(r io.Reader) error {
	err := t.CastToNRGBA(r)
	if err != nil {
		return errors.New("Decode failed: Error casting to NRGBA - (" + err.Error() + ")")
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

func (t *Texture) CastToNRGBA(r io.Reader) error {
	// Decode the skin
	textureImg, format, err := image.Decode(r)
	if err != nil {
		return errors.New("CastToNRGBA failed: (" + err.Error() + ")")
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
