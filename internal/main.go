package main

import (
	_ "internal/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"
	"internal/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.New())
}
