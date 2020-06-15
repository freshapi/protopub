package protopubcli

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/containerd/containerd/remotes"
	"github.com/containerd/containerd/remotes/docker"
	auth "github.com/deislabs/oras/pkg/auth/docker"
	"github.com/spf13/pflag"
	"net/http"
	"os"
)

type resolverOpts struct {
	username  string
	password  string
	insecure  bool
	plainHTTP bool
	configs   []string
}

func resolverOptions(flags *pflag.FlagSet) resolverOpts {
	opts := resolverOpts{}
	flags.StringVarP(&opts.username, "username", "u", "", "registry username")
	flags.StringVarP(&opts.password, "password", "p", "", "registry password")
	flags.BoolVar(&opts.insecure, "insecure", false, "skip registry TLS verification")
	flags.BoolVar(&opts.plainHTTP, "plain", false, "plain http registry")
	flags.StringArrayVarP(&opts.configs, "config", "c", nil, "auth config path")
	return opts
}

// Resolver creates OIC resolver which could be used to authenticate registry requests
func Resolver(options resolverOpts) remotes.Resolver {
	opts := docker.ResolverOptions{
		PlainHTTP: options.plainHTTP,
	}

	client := http.DefaultClient
	if options.insecure {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}
	opts.Client = client

	if options.username != "" || options.password != "" {
		opts.Credentials = func(hostName string) (string, string, error) {
			return options.username, options.password, nil
		}
		return docker.NewResolver(opts)
	}
	cli, err := auth.NewClient(options.configs...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: Error loading auth file: %v\n", err)
	}
	resolver, err := cli.Resolver(context.Background(), client, options.plainHTTP)
	if err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: Error loading resolver: %v\n", err)
		resolver = docker.NewResolver(opts)
	}
	return resolver
}
