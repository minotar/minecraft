// Minecraft Avatars
package minecraft

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"io"
	"net/http"
	"strings"
)

const CharHash = "98903c1609352e11552dca79eb1ce3d6"

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

type MojangProfileResponse struct {
	Uuid       string                  `json:"id"`
	Username   string                  `json:"name"`
	Properties []MojangProfileProperty `json:"properties"`
}

type MojangProfileProperty struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type MojangProfileTextureProperty struct {
	TimestampMs uint64 `json:"timestamp"`
	ProfileUuid string `json:"profileId"`
	ProfileName string `json:"profileName"`
	IsPublic    bool   `json:"isPublic"`
	Textures    struct {
		Skin struct {
			Url string `json:"url"`
		} `json:"SKIN"`
		Cape struct {
			Url string `json:"url"`
		} `json:"CAPE"`
	} `json:"textures"`
}

func GetSkin(u User) (Skin, error) {
	username := u.Name

	skin, err := FetchSkinFromMojang(username)

	return skin, err
}

func FetchSkinFromUrl(url, username string) (Skin, error) {
	resp, err := http.Get(url + username + ".png")
	if err != nil {
		return Skin{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Skin{}, errors.New("Skin not found. (" + fmt.Sprintf("%v", resp) + ")")
	}

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

func FetchSkinFromMojangByUuid(uuid string) (Skin, error) {
	uuid = strings.Replace(uuid, "-", "", 4)

	url := "https://sessionserver.mojang.com/session/minecraft/profile/"
	url += uuid

	resp, err := http.Get(url)
	if err != nil {
		return Skin{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return Skin{}, errors.New("Skin not found. (HTTP 204 No Content)")
	} else if resp.StatusCode == 429 { // StatusTooManyRequests
		return Skin{}, errors.New("Rate limited")
	} else if resp.StatusCode != http.StatusOK {
		return Skin{}, errors.New("Error retrieving profile. (HTTP " + resp.Status + ")")
	}

	mojangProfile := MojangProfileResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&mojangProfile); err != nil {
		return Skin{}, errors.New("Error decoding profile. (" + err.Error() + ")")
	}

	var texturesProperty *MojangProfileProperty
	for _, v := range mojangProfile.Properties {
		if v.Name == "textures" {
			texturesProperty = &v
			break
		}
	}
	if texturesProperty == nil {
		return Skin{}, errors.New("Profile " + uuid + " has no textures property?")
	}

	texturesSkin := MojangProfileTextureProperty{}
	if err := json.NewDecoder(base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(texturesProperty.Value))).Decode(&texturesSkin); err != nil {
		return Skin{}, errors.New("Error decoding texture property. (" + err.Error() + ")")
	}

	skinResp, err := http.Get(texturesSkin.Textures.Skin.Url)
	if err != nil {
		return Skin{}, errors.New("Error retrieving skin. " + err.Error())
	}
	defer skinResp.Body.Close()

	if skinResp.StatusCode != http.StatusOK {
		return Skin{}, errors.New("Error retrieving skin. (HTTP " + skinResp.Status + ")")
	}

	return DecodeSkin(skinResp.Body)
}

func DecodeSkin(r io.Reader) (Skin, error) {
	skinImg := castToNRGBA(r)

	// And md5 hash its pixels
	hasher := md5.New()
	hasher.Write(skinImg.(*image.NRGBA).Pix)
	md5Hash := fmt.Sprintf("%x", hasher.Sum(nil))

	// Create an md5 sum
	// Finally, establish the skin
	skin := Skin{
		Image: skinImg,
		Hash:  md5Hash,
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

func castToNRGBA(r io.Reader) image.Image {
	// Decode the skin
	var s image.Image
	skinImg, format, err := image.Decode(r)
	if err != nil {
		chr, _ := FetchImageForChar()
		s = chr
	} else {
		s = skinImg
		format = ""
	}
	// Convert it to NRGBA if necessary
	skinFinal := s
	if format != "NRGBA" {
		bounds := s.Bounds()
		skinFinal = image.NewNRGBA(bounds)
		draw.Draw(skinFinal.(draw.Image), bounds, s, image.Pt(0, 0), draw.Src)
	}

	return skinFinal
}
