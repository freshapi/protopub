package protopubcli

import (
	"context"
	"github.com/freshapi/protopub/pkg/protopub"
	"github.com/spf13/cobra"
	"os"
)

type pushOpts struct {
	resolverOpts
	image string
	file  string
	files []string
}

// NewPush creates `push` command
func NewPush() *cobra.Command {
	var opts pushOpts
	push := cobra.Command{
		Use:   "push <image> <descriptor-path>",
		Short: "Pushes descriptor from file to OCI image registry",
		Long: `Pushes descriptor from file to OCI image registry.

Example - push descriptor to the registry:
  protopub push docker.io/freshapi/example:latest ./my-descriptor-set.bin

Example - push descriptor to the private registry without logging in first:
  protopub push --username=freshapi --password=my_password docker.io/freshapi/example:latest ./my-descriptor-set.bin
`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.image = args[0]
			opts.file = args[1]
			return runPush(&opts)
		},
	}
	flags := push.Flags()
	opts.resolverOpts = resolverOptions(flags)
	flags.StringArrayVarP(&opts.files, "files", "f", []string{}, ".proto file names from FileDescriptorSet to push")

	return &push
}

func runPush(opts *pushOpts) error {
	resolver := Resolver(opts.resolverOpts)
	_, err := protopub.PushFile(context.Background(), os.Stdout, resolver, opts.image, opts.file, opts.files...)
	if err != nil {
		return err
	}
	return nil
}
