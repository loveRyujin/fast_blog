package main

import (
	"os"

	"github.com/onexstack_practice/fast_blog/cmd/fb-apiserver/app"
	_ "go.uber.org/automaxprocs"
)

func main() {
	cmd := app.NewFastBlogCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
