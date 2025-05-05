package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/ritarock/pa55vault/app"
)

func Execute() {
	flag.Parse()
	handler := app.NewHandler(flag.Args())
	err := handler.Run()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
