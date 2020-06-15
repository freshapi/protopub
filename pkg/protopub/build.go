package protopub

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// BuildToFile asks `protoc` to build descriptor set into output.
// basePath usually is a working directory
// paths - which paths to use as '-I' option of compiler
// output - output file path
func BuildToFile(basePath string, paths []string, output string) error {
	args := []string{
		"--descriptor_set_out=" + output,
		"--include_source_info",
	}
	for _, path := range paths {
		args = append(args, "-I"+path)
	}
	err := filepath.Walk(basePath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.HasSuffix(path, ".proto") {
				args = append(args, path)
			}
			return nil
		})
	if err != nil {
		return err
	}
	return exec.Command("protoc", args...).Run()
}
