package main

import (
	"fmt"
	"os"

	"github.com/onexstack_practice/fast_blog/cmd/fb-apiserver/app"
	_ "go.uber.org/automaxprocs"
)

func main() {
	cmd := app.NewFastBlogCommand()
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
