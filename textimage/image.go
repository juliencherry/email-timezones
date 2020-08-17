package textimage

import (
	"bytes"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"net/http"

	"github.com/golang/freetype"
	"golang.org/x/image/font"
)

func Write(text []string, imageBuffer *bytes.Buffer)  error {

	resp, err := http.Get("https://github.com/golang/freetype/blob/e2365dfdc4a05e4b8299a783240d4a7d5a65d4e4/testdata/luxisr.ttf?raw=true")
	if err != nil {
   	return err
	}

	defer resp.Body.Close()
	fontBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
   	return err
	}

	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return err
	}

	fg, bg := image.Black, image.White
	rgba := image.NewRGBA(image.Rect(0, 0, 155, 55))
	draw.Draw(rgba, rgba.Bounds(), bg, image.ZP, draw.Src)

	size := 12.0
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(f)
	c.SetFontSize(size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fg)
	c.SetHinting(font.HintingNone)

	pt := freetype.Pt(10, 10 + int(c.PointToFixed(size) >> 6))
		for _, s := range text {
			_, err = c.DrawString(s, pt)
			if err != nil {
				return err
			}
			pt.Y += c.PointToFixed(size * 1.5)
		}

	err = png.Encode(imageBuffer, rgba)
	if err != nil {
		return err
	}

	return nil
}
