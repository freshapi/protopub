package protopub

import (
	"context"
	"fmt"
	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/remotes"
	"github.com/deislabs/oras/pkg/content"
	"github.com/deislabs/oras/pkg/oras"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"google.golang.org/protobuf/types/descriptorpb"
	"io"
	"os"
	"sync"
)

func PushFile(ctx context.Context, statusW io.Writer, resolver remotes.Resolver, imageRef string, path string, files ...string) (*Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return PushDescriptor(ctx, statusW, resolver, imageRef, file, files...)
}

func PushDescriptor(ctx context.Context, statusW io.Writer, resolver remotes.Resolver, imageRef string, r io.Reader, files ...string) (*Image, error) {
	image, err := ImageFromReader(r)
	if err != nil {
		return nil, err
	}
	config, err := NewConfigFromFiles(image.Files)
	if err != nil {
		return nil, err
	}
	image.Config = config
	store := NewDescriptorStore()
	var filter FileFilter
	if len(files) > 0 {
		filter = func(pbFile *descriptorpb.FileDescriptorProto) bool {
			for _, file := range files {
				if file == pbFile.GetName() {
					return true
				}
			}
			return false
		}
	}
	manifest, err := PrepareImage(store, image, filter)
	if err != nil {
		return nil, err
	}
	var opts []oras.PushOpt
	opts = append(opts, oras.WithConfig(manifest.Config))
	opts = append(opts, oras.WithManifestAnnotations(manifest.Annotations))
	opts = append(opts, oras.WithPushBaseHandler(pushStatusTrack(statusW)))
	_, err = oras.Push(ctx, resolver, imageRef, store, manifest.Layers, opts...)
	return image, err
}

func pushStatusTrack(w io.Writer) images.Handler {
	var printLock sync.Mutex
	return images.HandlerFunc(func(ctx context.Context, desc ocispec.Descriptor) ([]ocispec.Descriptor, error) {
		if name, ok := content.ResolveName(desc); ok {
			printLock.Lock()
			defer printLock.Unlock()
			_, _ = fmt.Fprintln(w, desc.Digest.Encoded()[:12], name)
		}
		return nil, nil
	})
}
