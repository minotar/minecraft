package minecraft

import (
	"bytes"
	"encoding/base64"
)

func FetchSkinForChar() (Skin, error) {
	steveImgBytes, err := base64.StdEncoding.DecodeString(CHAR_SKIN_BASE64)
	if err != nil {
		return Skin{}, err
	}

	steveImgBuf := bytes.NewBuffer(steveImgBytes)

	skin, err := DecodeSkin(steveImgBuf)
	skin.Source = "Char"

	return skin, err
}

// The constant below contains Mojang AB copyrighted content.
// The use of this imagery is subject to their copyright.

const CHAR_SKIN_BASE64 = `iVBORw0KGgoAAAANSUhEUgAAAEAAAAAgCAMAAACVQ462AAABJlBMVEUAAAAAr68AmZkAqK
gmGgoAaGh1Ry8qHQ0Af3+WX0FqQDAvHw8oGwoAW1s6MYkAnp4pHAwsHg6GUzQrHg2qfWaW
b1skGAiBUzkwKHI/Pz8AzMyHVTtGOqUmGAsnGwsoGg0tHQ4tIBCaY0QzJBGcZ0gmIVsfEA
tRMSUjFwkkGAomGgwoGwsoHAsrHg4sHhEvIhEyIxA6KBRiQy90SC96TjOEUjGDVTtSPYmH
WDqIWjmKWTucY0WcaUyiakedak+0hG27iXL///8AAAB3QjUAYGAoKChWScxra2svIA00JR
I/KhVCKhJSKCZtQypvRSyAUzSPXj6QXkOWX0CcY0aaZEqfaEmcclysdlqze2K1e2etgG23
gnK2iWy+iGy9i3K9jnK9jnTGloDZnpKOAAAAAXRSTlMAQObYZgAAAk1JREFUeF6lkoVu5D
AURWsKcwYZyszMy8zM+/8/sdd2rW2n2s10ehTJV5Hu0XPyJgy2XcFTmbEumLgpdsX3IahY
XuZl2TgC3+dacAd4NxdUfHkH3x5bgC7n/m0Eqm8Ez8YQcD/hM1KgGEOQcGDbN/6NQcA7yV
q//yj3lyo8t/v9taTDg8DsxSgCnib5ZJ7nk0+fpEuTeZJyCMxeFArSNk+WF+eT/NvP3z9+
fXyczC8uJ7ydmr0onqDbS4NOulo6LX35cHZ0tJp2grTXDcxeFAuCYGFu73D99Hz389nueW
n9cG9uAS/1Xowk6LZXtt39t5++et7345f77vZKuwuB7I8iSIPexskgs16/uQteVa1scLLR
C1KzF4WCLBtkU+7B5ovn1erWzs7W5sH7qXdTgyyzOLCyiVF4AO4DI21Yisa9C0yBNRrMm/
VmiwWe4pqgAW4rYOh7t7nCdUGtNl1rARzyjBzHCWlcPa5W49gBMWOMEMZ0Juok4K9g+pIA
OEI4jMXoE4Is0AaMITqEkRjvhgW1KwIBAQoxRVMAomFKjFAoCIUIUYgpWjpTwqgWxOR/Ey
BEQIRhSKmsURoC1YbNdd1mCTQR5HlVEEWYQAvQEZTqmkCGB4lpQblUcpta8M8ryJIcgCmB
EKoNSVMLykpQLl8SgOG/YATIAtEIMLoUQASuC4ARENUiVGcNelrgDgnqdSOo1x8CRzgOpe
gzLUBWn8GVAuACeZVLAl00J64y3ZJEEMgliyJGadRSAjQxSZEAE2kD1RkBPhegV24OTfAH
F7Rm3RfymosAAAAASUVORK5CYII=`
