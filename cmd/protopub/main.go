package main

import (
	"github.com/freshapi/protopub/pkg/protopubcli"
	"os"
)

var version = "unknown"

func main() {
	protopub := protopubcli.NewProtopub()
	protopub.Version = version
	err := protopub.Execute()
	if err != nil {
		os.Exit(1)
	}
}
