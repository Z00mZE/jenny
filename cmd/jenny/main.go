package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/Z00mZE/jenny/internal/app/jenny"
)

func main() {
	ctx, ctxCloser := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer ctxCloser()

	fmt.Println(jenny.Run(ctx))
}
