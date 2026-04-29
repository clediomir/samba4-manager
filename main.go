package main

import (
	"log/slog"
	"os"

	"samba4-manager/cmd"
)

func main() {
	if err := cmd.Execute(TemplatesFS, StaticFS, LocalesFS); err != nil {
		slog.Error("Startup failure", "error", err)
		os.Exit(1)
	}
}
