package main

import (
	_ "sbcounter/internal/packed"

	"sbcounter/internal/cmd"

	"github.com/gogf/gf/v2/os/gctx"
)

func main() {
	cmd.Main.Run(gctx.New())
}
