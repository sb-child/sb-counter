package counter

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"os"
	"path/filepath"
	"sb-counter/internal/consts"
	"sb-counter/internal/service"
	"sb-counter/internal/utils"
	"time"

	"github.com/disintegration/imaging"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/grand"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type sCounter struct {
	ImgSize     image.Rectangle
	FontRegular *truetype.Font
	FontMono    *truetype.Font
	Users       []consts.CounterUser
}

func prepareFont(path string) *truetype.Font {
	tt, err := freetype.ParseFont(utils.GetResource("resource/public/resource/font/" + path))
	if err != nil {
		panic(err)
	}
	return tt
}

func init() {
	service.RegisterCounter(New())
}

func New() *sCounter {
	users := []consts.CounterUser{}
	g.Config().MustGet(context.Background(), "sbcounter.user").Structs(&users)
	g.Log().Debug(context.Background(), users)
	return &sCounter{
		ImgSize:     image.Rect(0, 0, 350, 150),
		FontRegular: prepareFont("OPPOSans-R.ttf"),
		FontMono:    prepareFont("JetBrainsMono-Regular.ttf"),
		Users:       users,
	}
}

func (s *sCounter) GetUser(ctx context.Context, userPath string) *consts.CounterUser {
	selectedUser := -1
	for i, user := range s.Users {
		if user.Path == userPath {
			selectedUser = i
			break
		}
	}
	if selectedUser == -1 {
		return nil
	}
	return &s.Users[selectedUser]
}

func (s *sCounter) DrawCard(ctx context.Context, today, all, yesterday, beforeYesterday int) []byte {
	img := image.NewRGBA(s.ImgSize)
	bgImg := s.GetCardBackground(ctx)
	draw.Draw(img, s.ImgSize, bgImg, bgImg.Bounds().Min, draw.Src)
	img = s.drawMainCounter(img, all)
	img = s.drawDailyCounter(img, today, yesterday, beforeYesterday)
	img = s.drawTime(img)
	// 使用jpg减少图片传输开销
	buff := new(bytes.Buffer)
	jpeg.Encode(buff, img, &jpeg.Options{Quality: 80})
	return buff.Bytes()
}

func (s *sCounter) drawText(src *image.RGBA, x, y int, size float64, text string, mono bool) *image.RGBA {
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
					return s.FontMono
				} else {
					return s.FontRegular
				}
			}(),
			&truetype.Options{Size: size}),
		Dot: point,
	}
	mainCounterFont.DrawString(text)
	return drawDst
}

func (s *sCounter) drawMainCounter(src *image.RGBA, v int) *image.RGBA {
	dst := s.drawText(src, 10, 43, 50, fmt.Sprintf("%11d", v), true)
	for i := 10; i < dst.Rect.Dx()-10; i++ {
		for j := -1; j <= 1; j++ {
			dst.SetRGBA(i, 80+j, color.RGBA{0, 0, 0, 0xff})
		}
	}
	return s.drawText(dst, 235, 70, 25, "总访问量", false)
}

func (s *sCounter) drawDailyCounter(src *image.RGBA, v, y, b int) *image.RGBA {
	dst := src
	t := fmt.Sprintf("日活%d", v)
	if y <= 0 {
		t = fmt.Sprintf("%s/昨无", t)
	} else {
		t = fmt.Sprintf("%s/%+d", t, v-y)
	}
	if b <= 0 {
		t = fmt.Sprintf("%s/前无", t)
	} else {
		t = fmt.Sprintf("%s/%+d", t, y-b)
	}
	dst = s.drawText(dst, 10, 110, 20, t, false)
	return dst
}

func (s *sCounter) drawTime(src *image.RGBA) *image.RGBA {
	now := time.Now()
	dst := s.drawText(src, 10, src.Rect.Dy()-10, 25,
		fmt.Sprintf("%04d.%02d.%02d %02d:%02d:%02d",
			now.Year(), now.Month(), now.Day(),
			now.Hour(), now.Minute(), now.Second()), false)
	return dst
}

func (s *sCounter) GetCardBackground(ctx context.Context) *image.RGBA {
retry:
	backgroundImgDir := g.Config().MustGet(ctx, "sbcounter.backgroundImageDir").String()
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
	g.Log().Debug(ctx, "[background]:", selectedBackgroundImg)
	img := image.NewRGBA(s.ImgSize)
	bgImg, _, err := image.Decode(f)
	if err != nil {
		goto retry
	}
	bgImg = imaging.Resize(bgImg, s.ImgSize.Dx(), 0, imaging.Linear)
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
	scrollMax := bgImg.Bounds().Dy() - s.ImgSize.Dy()
	scroll := grand.Intn(scrollMax)
	draw.Draw(img, s.ImgSize, bgImg, image.Pt(0, scroll), draw.Src)
	return img
}
