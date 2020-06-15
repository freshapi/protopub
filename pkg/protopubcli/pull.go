package protopubcli

import (
	"context"
	"github.com/freshapi/protopub/pkg/protopub"
	"github.com/spf13/cobra"
	"os"
)

type pullOpts struct {
	resolverOpts
	image string
	file  string
}

// NewPull creates `pull` command
func NewPull() *cobra.Command {
	var opts pullOpts
	push := cobra.Command{
		Use:   "pull <image> <descriptor-path>",
		Short: "Pulls descriptor set file from OCI image registry",
		Long: `Pulls descriptor set file from OCI image registry.

Example - pull descriptor from the registry:
  protopub pull docker.io/freshapi/example:latest ./my-descriptor-set.bin

Example - pull descriptor from private registry without logging in first:
  protopub pull --username=freshapi --password=my_password docker.io/freshapi/example:latest ./my-descriptor-set.bin
`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.image = args[0]
			opts.file = args[1]
			return runPull(&opts)
		},
	}
	flags := push.Flags()
	opts.resolverOpts = resolverOptions(flags)

	return &push
}

func runPull(opts *pullOpts) error {
	resolver := Resolver(opts.resolverOpts)
	err := protopub.PullToFile(context.Background(), os.Stdout, resolver, opts.image, opts.file)
	if err != nil {
		return err
	}
	return nil
}
