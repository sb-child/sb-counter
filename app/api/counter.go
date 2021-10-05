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
	"image/draw"
	"image/jpeg"
	"sb-counter/app/service"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

var Counter = CounterApi{}

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

type CounterApi struct{}

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

func (api *CounterApi) drawCard() []byte {
	img := image.NewRGBA(image.Rect(0, 0, 350, 150))
	
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
