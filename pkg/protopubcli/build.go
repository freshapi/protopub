package protopubcli

import (
	"github.com/freshapi/protopub/pkg/protopub"
	"github.com/spf13/cobra"
)

type buildOpts struct {
	output string
	paths  []string
}

// NewBuild creates `build` command
func NewBuild() *cobra.Command {
	var opts buildOpts
	inspect := cobra.Command{
		Use:   "build <output> [path...]",
		Short: "Builds given paths into protobuf descriptor set",
		Long: `Build composes .proto files to the descriptor set file.
It searches for .proto files recursively in given paths and runs protobuf compiler.

NOTE: current working directory is used as base path for .proto files when building descriptors

Example - build all .proto files from current directory:
  protopub build my-descriptor-set.bin

Example - build all .proto files from particular directory:
  protopub build my-descriptor-set.bin ./my_app_1

Example - build multiple directories:
  protopub build my-descriptor-set.bin ./my_app_1 ./my_app_2
`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.output = args[0]
			opts.paths = args[1:]
			return runBuild(&opts)
		},
	}

	return &inspect
}

func runBuild(opts *buildOpts) error {
	if len(opts.paths) == 0 {
		opts.paths = []string{"."}
	}
	return protopub.BuildToFile(".", opts.paths, opts.output)
}
