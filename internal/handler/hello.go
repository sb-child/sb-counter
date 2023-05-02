package handler

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"sbcounter/apiv1"
)

var (
	Hello = hHello{}
)

type hHello struct{}

func (h *hHello) Hello(ctx context.Context, req *apiv1.HelloReq) (res *apiv1.HelloRes, err error) {
	g.RequestFromCtx(ctx).Response.Writeln("Hello World!")
	return
}
