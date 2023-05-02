package counter

import (
	"sb-counter/internal/service"

	"github.com/gogf/gf/v2/net/ghttp"
	// "context"
	// "github.com/gogf/gf/v2/frame/g"
)

func Handler(r *ghttp.Request) {
	ctx := r.GetCtx()
	userPath := r.Get("user_path", "").String()
	mode := r.Get("mode", "").String()
	output := r.Get("output", "").String()
	user := service.Counter().GetUser(ctx, userPath)
	if user == nil {
		r.Response.Status = 404
		r.Response.Writeln("user not found")
		return
	}
	r.Response.Header().Set("Cache-Control", "no-cache,max-age=0,no-store,s-maxage=0,proxy-revalidate")
	switch mode {
	case "rw":
		service.Database().Add(ctx, user.DB, r.GetClientIp())
	case "ro":
		// readonly
	default:
		// other
	}
	userData := service.Database().FetchData(ctx, user.DB)
	if userData == nil {
		r.Response.Status = 502
		r.Response.Writeln("internal error")
		return
	}
	switch output {
	case "card":
		r.Response.Header().Set("Content-Type", "image/jpeg")
		r.Response.Write(service.Counter().DrawCard(
			ctx, userData.Today, userData.All, userData.Yesterday, userData.BeforeYesterday,
		))
	case "json":
		r.Response.Header().Set("Content-Type", "application/json")
	default:
	}
}

// type Controller struct{}

// func New() *Controller {
// 	return &Controller{}
// }

// func (c *Controller) Hello(ctx context.Context, req *v1.Req) (res *v1.Res, err error) {
// 	g.RequestFromCtx(ctx).Response.Writeln("Hello World!")
// 	return
// }
