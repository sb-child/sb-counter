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
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

var Hello = helloApi{}

type helloApi struct{}

type User struct {
	Path string `json:"path"`
	View string `json:"view"`
	DB   string `json:"db"`
}

// Index is a demonstration route handler for output "Hello World!".
func (*helloApi) Index(r *ghttp.Request) {
	userPath := r.Get("user_path")
	r.Response.Writeln(r.Get("method"))
	users := []User{}
	g.Config().GetStructs("sbcounter.user", &users)
	r.Response.Writeln(users)
	for _, user := range users {
		if user.Path == userPath{
			r.Response.Writeln(user.DB)
		}
	}

}
