package protopub

import (
	"context"
	"fmt"
	"github.com/freshapi/protopub/schema"
	v1 "github.com/freshapi/protopub/schema/v1"
	"io"
	"os"
	"sync"

	"github.com/containerd/containerd/images"
	"github.com/containerd/containerd/remotes"
	"github.com/deislabs/oras/pkg/content"
	"github.com/deislabs/oras/pkg/oras"
	"github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func PullToFile(ctx context.Context, statusW io.Writer, resolver remotes.Resolver, imageRef string, path string, files ...string) error {
	desc, err := PullImage(ctx, statusW, resolver, imageRef, files...)
	if err != nil {
		return err
	}
	var pbFiles []*descriptorpb.FileDescriptorProto
	for _, file := range desc.Files {
		pbFiles = append(pbFiles, file)
	}
	set := descriptorpb.FileDescriptorSet{
		File: pbFiles,
	}
	bytes, err := proto.Marshal(&set)
	if err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

func PullImage(ctx context.Context, statusW io.Writer, resolver remotes.Resolver, imageRef string, files ...string) (*Image, error) {
	store := NewDescriptorStore()
	var opts []oras.PullOpt
	opts = append(opts, oras.WithPullCallbackHandler(pullStatusTrack(statusW)))
	_, layers, err := oras.Pull(ctx, resolver, imageRef, store, opts...)
	if err != nil {
		return nil, err
	}
	var pbs []*descriptorpb.FileDescriptorProto
	var config schema.Config
	for _, layer := range layers {
		_, bytes, ok := store.Get(layer)
		if layer.MediaType == ocispec.MediaTypeImageConfig {
			config, err = v1.ParseConfig(bytes)
			if err != nil {
				return nil, err
			}
			continue
		}
		if !ok {
			return nil, fmt.Errorf("error fetching layer data")
		}
		pb := descriptorpb.FileDescriptorProto{}
		err = proto.Unmarshal(bytes, &pb)
		if err != nil {
			return nil, err
		}
		pbs = append(pbs, &pb)
	}
	image, err := ImageFromDescriptorSet(&descriptorpb.FileDescriptorSet{
		File: pbs,
	})
	if err != nil {
		return nil, err
	}
	image.Config = config
	return image, nil
}

func pullStatusTrack(w io.Writer) images.Handler {
	var printLock sync.Mutex
	return images.HandlerFunc(func(ctx context.Context, desc ocispec.Descriptor) ([]ocispec.Descriptor, error) {
		if name, ok := content.ResolveName(desc); ok {
			digestString := desc.Digest.String()
			if err := desc.Digest.Validate(); err == nil {
				if algo := desc.Digest.Algorithm(); algo == digest.SHA256 {
					digestString = desc.Digest.Encoded()[:12]
				}
			}
			printLock.Lock()
			defer printLock.Unlock()
			_, _ = fmt.Fprintln(w, digestString, name)
		}
		return nil, nil
	})
}
