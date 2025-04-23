package main

import (
	_ "keyboard-api-go/internal/packed"

	_ "keyboard-api-go/internal/logic"

	"github.com/gogf/gf/v2/os/gctx"

	"keyboard-api-go/internal/cmd"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"

	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
