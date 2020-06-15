package protopub

import (
	"github.com/freshapi/protopub/schema"
	"github.com/opencontainers/go-digest"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// FileFilter is used to filter proto files
type FileFilter func(file *descriptorpb.FileDescriptorProto) bool

// PrepareImage creates OCI image manifest which could be used to push image
func PrepareImage(store *DescriptorStore, image *Image, filter FileFilter) (ocispec.Manifest, error) {
	config, err := AddConfig(store, image)
	if err != nil {
		return ocispec.Manifest{}, err
	}
	var files []*descriptorpb.FileDescriptorProto
	if filter == nil {
		files = image.Files
	} else {
		for _, file := range image.Files {
			if filter(file) {
				files = append(files, file)
			}
		}
	}
	var layers []ocispec.Descriptor
	for _, file := range files {
		layer, err := AddFile(store, file)
		if err != nil {
			return ocispec.Manifest{}, err
		}
		layers = append(layers, layer)
	}
	manifest := ocispec.Manifest{
		Config:      config,
		Layers:      layers,
		Annotations: GetAnnotations(image),
	}
	return manifest, nil
}

// AddFile adds proto file to the store and returns its descriptor
func AddFile(store *DescriptorStore, file *descriptorpb.FileDescriptorProto) (ocispec.Descriptor, error) {
	bytes, err := proto.Marshal(file)
	if err != nil {
		return ocispec.Descriptor{}, err
	}
	fileDesc := ocispec.Descriptor{
		MediaType: ocispec.MediaTypeImageLayer,
		Digest:    digest.FromBytes(bytes),
		Size:      int64(len(bytes)),
		Annotations: map[string]string{
			ocispec.AnnotationTitle: file.GetName(),
		},
	}
	store.Set(fileDesc, bytes)
	return fileDesc, nil
}

// GetAnnotations returns image annotations
func GetAnnotations(image *Image) map[string]string {
	return make(map[string]string)
}

// AddConfig adds config to the store and returns its descriptor
func AddConfig(store *DescriptorStore, image *Image) (ocispec.Descriptor, error) {
	configBytes, err := schema.RenderConfig(image.Config)
	if err != nil {
		return ocispec.Descriptor{}, err
	}
	configDesc := ocispec.Descriptor{
		MediaType: ocispec.MediaTypeImageConfig,
		Digest:    digest.FromBytes(configBytes),
		Size:      int64(len(configBytes)),
		Annotations: map[string]string{
			ocispec.AnnotationTitle: "_config.json",
		},
	}
	store.Set(configDesc, configBytes)
	return configDesc, nil
}
