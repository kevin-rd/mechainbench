package main

import (
	"context"
	"mechainbench/app/core"
	"mechainbench/app/engine"
	"time"
)

func main() {
	benchContext := &core.Context{
		BatchNum: 10,
		Quantity: 10,
		Timeout:  30 * time.Second,
	}
	engine.NewDefaultEngine(benchContext).Run(context.Background())
}
