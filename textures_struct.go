package minecraft

import (
	"crypto/md5"
	"fmt"
	"image"
	"image/draw"
	"io"
	// If we work with PNGs we need this
	_ "image/png"

	"github.com/pkg/errors"
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

// CastToNRGBA takes image bytes and converts to NRGBA format if needed
func (t *Texture) CastToNRGBA(r io.Reader) error {
	// Decode the skin
	textureImg, format, err := image.Decode(r)
	if err != nil {
		return errors.Wrap(err, "unable to CastToNRGBA")
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

// Decode takes the image bytes and turns it into our Texture struct
func (t *Texture) Decode(r io.Reader) error {
	err := t.CastToNRGBA(r)
	if err != nil {
		return errors.WithStack(err)
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

// Fetch performs the GET for the texture, doing any required conversion and saving our Image property
func (t *Texture) Fetch() error {
	apiBody, err := apiRequest(t.URL)
	if apiBody != nil {
		defer apiBody.Close()
	}
	if err != nil {
		return errors.Wrap(err, "unable to Fetch Texture")
	}

	err = t.Decode(apiBody)
	if err != nil {
		return errors.Wrap(err, "unable to Decode Texture")
	}
	return nil
}

// FetchWithTextureProperty takes a already decoded Texture Property and will request either Skin or Cape as instructed
func (t *Texture) FetchWithTextureProperty(profileTextureProperty SessionProfileTextureProperty, textureType string) error {
	if textureType == "Skin" {
		t.URL = profileTextureProperty.Textures.Skin.URL
	} else {
		t.URL = profileTextureProperty.Textures.Cape.URL
	}
	if t.URL == "" {
		return errors.Errorf("%s URL not present", textureType)
	}
	t.Source = "SessionProfile"

	return errors.Wrap(t.Fetch(), "FetchWithTextureProperty failed")
}

// FetchWithSessionProfile will decode the Texture Property for you and request the Skin or Cape as instructed
// If requesting both Skin and Cape, this would result in 2 x decoding - useFetchWithTextureProperty instead
func (t *Texture) FetchWithSessionProfile(sessionProfile SessionProfileResponse, textureType string) error {
	profileTextureProperty, err := DecodeTextureProperty(sessionProfile)
	if err != nil {
		return errors.WithStack(err)
	}

	return t.FetchWithTextureProperty(profileTextureProperty, textureType)
}

// FetchWithUsernameMojang takes a username and will then request from the deprecated (but still updated) Mojang source
func (t *Texture) FetchWithUsernameMojang(username string, textureType string) error {
	if textureType == "Skin" {
		t.URL = "http://skins.minecraft.net/MinecraftSkins/" + username + ".png"
	} else {
		t.URL = "http://skins.minecraft.net/MinecraftCloaks/" + username + ".png"
	}
	t.Source = "Mojang"

	return errors.Wrap(t.Fetch(), "FetchWithUsernameMojang failed")
}

// FetchWithUsernameS3 uses the deprecated and stale S3 store for the Skin or Cape (truly last resort)
func (t *Texture) FetchWithUsernameS3(username string, textureType string) error {
	if textureType == "Skin" {
		t.URL = "http://s3.amazonaws.com/MinecraftSkins/" + username + ".png"
	} else {
		t.URL = "http://s3.amazonaws.com/MinecraftCloaks/" + username + ".png"
	}
	t.Source = "S3"

	return errors.Wrap(t.Fetch(), "FetchWithUsernameS3 failed")
}
