package protopubcli

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/freshapi/protopub/pkg/protopub"
)

type inspectOpts struct {
	resolverOpts
	target string
}

// NewInspect creates `inspect` command
func NewInspect() *cobra.Command {
	var opts inspectOpts
	inspect := cobra.Command{
		Use:   "inspect [target]",
		Short: "Inspect image or descriptor set file.",
		Long: `Inspect image or descriptor set file.

Example - inspect local file
  protopub inspect ./my-descriptor-set.bin

Example - inspect remote image
  protopub inspect docker.io/freshapi/example:latest
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.target = args[0]
			return runInspect(&opts)
		},
	}
	opts.resolverOpts = resolverOptions(inspect.Flags())

	return &inspect
}

func runInspect(opts *inspectOpts) error {
	var image *protopub.Image
	_, err := os.Stat(opts.target)
	if os.IsNotExist(err) {
		resolver := Resolver(opts.resolverOpts)
		image, err = protopub.PullImage(context.Background(), os.Stdout, resolver, opts.target)
	} else {
		image, err = protopub.ImageFromPath(opts.target)
		if err == nil {
			config, err := protopub.NewConfigFromFiles(image.Files)
			if err != nil {
				return err
			}
			image.Config = config
		}
	}
	if err != nil {
		return err
	}
	info := protopub.Info(image)
	b, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}
