// Copyright 2021 sbchild

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"sb-counter/app/service"
	"time"

	"github.com/disintegration/imaging"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/grand"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

var Counter = CounterApi{
	ImgSize:     image.Rect(0, 0, 350, 150),
	FontRegular: prepareFont("OPPOSans-R.ttf"),
	FontMono:    prepareFont("JetBrainsMono-Regular.ttf"),
}

type User struct {
	Path string `json:"path"`
	View string `json:"view"`
	DB   string `json:"db"`
}

// type View struct {
// 	Path string `json:"path"`
// 	View string `json:"view"`
// 	DB   string `json:"db"`
// }

type CounterApi struct {
	ImgSize     image.Rectangle
	FontRegular *truetype.Font
	FontMono    *truetype.Font
}

func prepareFont(path string) *truetype.Font {
	tt, err := freetype.ParseFont(g.Res().GetContent("public/resource/font/" + path))
	if err != nil {
		panic(err)
	}
	return tt
}

func (api *CounterApi) getCardBackground() *image.RGBA {
retry:
	backgroundImgDir := g.Config().GetString("sbcounter.backgroundImageDir")
	backgroundImgs := make([]string, 0)
	backgroundPNGImgs, _ := filepath.Glob(backgroundImgDir + "*.png")
	backgroundJPGImgs, _ := filepath.Glob(backgroundImgDir + "*.jpg")
	backgroundImgs = append(backgroundImgs, backgroundPNGImgs...)
	backgroundImgs = append(backgroundImgs, backgroundJPGImgs...)
	selectedBackgroundImg := backgroundImgs[grand.Intn(len(backgroundImgs))]
	f, err := os.Open(selectedBackgroundImg)
	if err != nil {
		goto retry
	}
	defer f.Close()
	g.Log().Debug("bg", selectedBackgroundImg)
	img := image.NewRGBA(api.ImgSize)
	bgImg, _, err := image.Decode(f)
	if err != nil {
		goto retry
	}
	bgImg = imaging.Resize(bgImg, api.ImgSize.Dx(), 0, imaging.Linear)
	bgImg = imaging.AdjustFunc(bgImg, func(c color.NRGBA) color.NRGBA {
		v := 400
		m := 90
		if (int(c.R) + int(c.G) + int(c.B)) < v {
			c.R = uint8(int(c.R) + m)
			c.G = uint8(int(c.G) + m)
			c.B = uint8(int(c.B) + m)
		}
		return c
	})
	bgImg = imaging.AdjustBrightness(bgImg, 10)
	bgImg = imaging.Blur(bgImg, 4)
	scrollMax := bgImg.Bounds().Dy() - api.ImgSize.Dy()
	scroll := grand.Intn(scrollMax)
	draw.Draw(img, api.ImgSize, bgImg, image.Pt(0, scroll), draw.Src)
	return img
}

func (api *CounterApi) drawText(src *image.RGBA, x, y int, size float64, text string, mono bool) *image.RGBA {
	posX := x
	posY := y
	bg := image.NewUniform(color.RGBA{0, 0, 0, 0xff})
	point := fixed.Point26_6{X: fixed.Int26_6(posX * 64), Y: fixed.Int26_6(posY * 64)}
	drawDst := image.NewRGBA(src.Bounds())
	draw.Draw(drawDst, drawDst.Bounds(), src, src.Bounds().Min, draw.Src)
	mainCounterFont := &font.Drawer{
		Dst: drawDst,
		Src: bg,
		Face: truetype.NewFace(
			func() *truetype.Font {
				if mono {
					return api.FontMono
				} else {
					return api.FontRegular
				}
			}(),
			&truetype.Options{Size: size}),
		Dot: point,
	}
	mainCounterFont.DrawString(text)
	return drawDst
}

func (api *CounterApi) drawMainCounter(src *image.RGBA, v int) *image.RGBA {
	dst := api.drawText(src, 10, 43, 50, fmt.Sprintf("%11d", v), true)
	for i := 10; i < dst.Rect.Dx()-10; i++ {
		for j := -1; j <= 1; j++ {
			dst.SetRGBA(i, 80+j, color.RGBA{0, 0, 0, 0xff})
		}
	}
	return api.drawText(dst, 235, 70, 25, "总访问量", false)
}

func (api *CounterApi) drawDailyCounter(src *image.RGBA, v, y int) *image.RGBA {
	dst := src
	var t string
	if y <= 0 {
		t = fmt.Sprintf("日活%d", v)
	} else {
		t = fmt.Sprintf("日活%d, %+d", v, v-y)
	}
	dst = api.drawText(dst, 10, 110, 25, t, false)
	return dst
}

func (api *CounterApi) drawTime(src *image.RGBA) *image.RGBA {
	now := time.Now()
	dst := api.drawText(src, 10, src.Rect.Dy()-10, 25,
		fmt.Sprintf("%04d.%02d.%02d %02d:%02d:%02d",
			now.Year(), now.Month(), now.Day(),
			now.Hour(), now.Minute(), now.Second()), false)
	return dst
}

func (api *CounterApi) drawCard(today, all, yesterday int) []byte {
	img := image.NewRGBA(api.ImgSize)
	bgImg := api.getCardBackground()
	draw.Draw(img, api.ImgSize, bgImg, bgImg.Bounds().Min, draw.Src)
	img = api.drawMainCounter(img, all)
	img = api.drawDailyCounter(img, today, yesterday)
	img = api.drawTime(img)
	// 使用jpg减少图片传输开销
	buff := new(bytes.Buffer)
	jpeg.Encode(buff, img, &jpeg.Options{Quality: 80})
	return buff.Bytes()
}

func (api *CounterApi) handleMode(mode string, r *ghttp.Request, db string) {
	switch mode {
	case "rw":
		service.Counter().Add(db, r.GetClientIp())
	case "ro":
		// readonly
		return
	default:
		// other
		return
	}
}

func (api *CounterApi) handleOutput(output string, r *ghttp.Request, today, all, yesterday int) {
	r.Response.Header().Set("Cache-Control", "no-cache,max-age=0,no-store,s-maxage=0,proxy-revalidate")
	switch output {
	case "card":
		r.Response.Header().Set("Content-Type", "image/jpeg")
		r.Response.Write(api.drawCard(today, all, yesterday))
	case "json":
		r.Response.Header().Set("Content-Type", "application/json")
	default:
		return
	}
}

func (api *CounterApi) Index(r *ghttp.Request) {
	userPath := r.GetString("user_path")
	mode := r.GetString("mode")
	output := r.GetString("output")
	users := []User{}
	selectedUser := -1
	g.Config().GetStructs("sbcounter.user", &users)
	for i, user := range users {
		if user.Path == userPath {
			selectedUser = i
			break
		}
	}
	if selectedUser == -1 {
		return
	}
	api.handleMode(mode, r, users[selectedUser].DB)
	all := service.Counter().GetAll(users[selectedUser].DB)
	today := service.Counter().GetDay(users[selectedUser].DB, 1)
	yesterday := service.Counter().GetDay(users[selectedUser].DB, 2)
	api.handleOutput(output, r, today, all, yesterday)
}
