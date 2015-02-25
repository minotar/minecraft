package minecraft

import (
	"crypto/md5"
	"errors"
	"fmt"
	"image"
	"image/draw"
	// If we work with PNGs we need this
	_ "image/png"
	"io"
)

// DecodeSkin takes the image bytes and turns it into our Skin struct
func DecodeSkin(r io.Reader) (Skin, error) {
	skinImg, err := castToNRGBA(r)
	if err != nil {
		return Skin{}, errors.New("Error casting to NRGBA (" + err.Error() + ")")
	}

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

func castToNRGBA(r io.Reader) (image.Image, error) {
	// Decode the skin
	var s image.Image
	skinImg, format, err := image.Decode(r)
	if err != nil {
		return nil, err
	}

	s = skinImg
	format = ""

	// Convert it to NRGBA if necessary
	skinFinal := s
	if format != "NRGBA" {
		bounds := s.Bounds()
		skinFinal = image.NewNRGBA(bounds)
		draw.Draw(skinFinal.(draw.Image), bounds, s, image.Pt(0, 0), draw.Src)
	}

	return skinFinal, nil
}
