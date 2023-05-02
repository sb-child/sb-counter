package main

import (
	_ "sb-counter/internal/packed"
	_ "sb-counter/internal/logic"

	"github.com/gogf/gf/v2/os/gctx"

	"sb-counter/internal/cmd"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
)

func main() {
	cmd.Main.Run(gctx.New())
}
