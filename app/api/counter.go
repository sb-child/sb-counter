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
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"sb-counter/app/service"

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

type View struct {
	Path string `json:"path"`
	View string `json:"view"`
	DB   string `json:"db"`
}

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
	g.Log().Info(selectedBackgroundImg)
	img := image.NewRGBA(api.ImgSize)
	bgImg, _, err := image.Decode(f)
	if err != nil {
		goto retry
	}
	bgImg = imaging.Resize(bgImg, api.ImgSize.Dx(), 0, imaging.Linear)
	bgImg = imaging.Blur(bgImg, 2)
	bgImg = imaging.AdjustBrightness(bgImg, 10)
	scrollMax := bgImg.Bounds().Dy() - api.ImgSize.Dy()
	scroll := grand.Intn(scrollMax)
	draw.Draw(img, api.ImgSize, bgImg, image.Pt(0, scroll), draw.Src)
	return img
}

func (api *CounterApi) drawMainCounter(src *image.RGBA) *image.RGBA {
	size := 50.0
	posX := 0
	posY := 50
	bg := image.NewUniform(color.RGBA{0, 0, 0, 0xff})
	point := fixed.Point26_6{X: fixed.Int26_6(posX * 64), Y: fixed.Int26_6(posY * 64)}
	drawDst := image.NewRGBA(src.Bounds())
	draw.Draw(drawDst, drawDst.Bounds(), src, src.Bounds().Min, draw.Src)
	mainCounterFont := &font.Drawer{
		Dst:  drawDst,
		Src:  bg,
		Face: truetype.NewFace(api.FontRegular, &truetype.Options{Size: size}),
		Dot:  point,
	}
	mainCounterFont.DrawString("45678中文")
	return drawDst
}

func (api *CounterApi) drawCard() []byte {
	img := image.NewRGBA(api.ImgSize)
	bgImg := api.getCardBackground()
	draw.Draw(img, api.ImgSize, bgImg, bgImg.Bounds().Min, draw.Src)
	img = api.drawMainCounter(img)
	// 使用jpg减少图片传输开销
	buff := new(bytes.Buffer)
	jpeg.Encode(buff, img, &jpeg.Options{Quality: 80})
	return buff.Bytes()
}

func (api *CounterApi) handleOutput(output string, r *ghttp.Request) {
	r.Response.Header().Set("Cache-Control", "no-cache,max-age=0,no-store,s-maxage=0,proxy-revalidate")
	switch output {
	case "card":
		r.Response.Header().Set("Content-Type", "image/jpeg")
		r.Response.Write(api.drawCard())
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
	api.handleOutput(output, r)
}
