package main

import (
	gotenbergcmd "github.com/gotenberg/gotenberg/v8/cmd"
	_ "github.com/klippa-app/go-pdfium/multi_threaded/worker"

	// Gotenberg modules.
	_ "github.com/gotenberg/gotenberg/v8/pkg/standard"
)

func main() {
	gotenbergcmd.Run()
}
