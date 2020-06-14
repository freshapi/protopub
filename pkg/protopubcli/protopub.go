package protopubcli

import "github.com/spf13/cobra"

func NewProtopub() *cobra.Command {
	protopub := cobra.Command{
		Use:          "protopub [command]",
		Short:        "Protopub is a tool which allows you to push and pull .proto descriptors from OCI registry",
		SilenceUsage: true,
	}
	protopub.AddCommand(NewPush(), NewPull(), NewInspect(), NewBuild(), NewLogin())
	return &protopub
}
