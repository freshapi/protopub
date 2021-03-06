package protopub

import (
	v1 "github.com/freshapi/protopub/schema/v1"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"io"
	"io/ioutil"
	"os"

	"github.com/freshapi/protopub/schema"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

type Image struct {
	Config schema.Config
	Files  []*descriptorpb.FileDescriptorProto
}

type ImageInfo struct {
	Config schema.Config `json:"config"`
	Files  []*ImageFile  `json:"files"`
}

type ImageFile struct {
	Name              string   `json:"name"`
	Syntax            string   `json:"syntax"`
	Package           string   `json:"package"`
	Imports           []string `json:"imports"`
	Messages          []string `json:"messages"`
	Enums             []string `json:"enums"`
	Services          []string `json:"services"`
	HasSourceCodeInfo bool     `json:"hasSourceCodeInfo"`
}

func Info(image *Image, layers ...*ocispec.Descriptor) *ImageInfo {
	info := ImageInfo{}
	info.Config = image.Config

	for _, f := range image.Files {
		var (
			services []string
			enums    []string
			messages []string
		)
		for _, s := range f.GetService() {
			services = append(services, s.GetName())
		}
		for _, e := range f.GetEnumType() {
			enums = append(enums, e.GetName())
		}
		for _, m := range f.GetMessageType() {
			messages = append(messages, m.GetName())
		}
		file := ImageFile{
			Name:              f.GetName(),
			Syntax:            f.GetSyntax(),
			Package:           f.GetPackage(),
			Imports:           f.GetDependency(),
			Services:          services,
			Enums:             enums,
			Messages:          messages,
			HasSourceCodeInfo: f.SourceCodeInfo != nil && f.SourceCodeInfo.Location != nil,
		}
		info.Files = append(info.Files, &file)
	}

	return &info
}

func ImageFromPath(path string) (*Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ImageFromReader(file)
}

func ImageFromReader(r io.Reader) (*Image, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var fd descriptorpb.FileDescriptorSet
	err = proto.Unmarshal(b, &fd)
	if err != nil {
		return nil, err
	}
	return ImageFromDescriptorSet(&fd)
}

func ImageFromDescriptorSet(fd *descriptorpb.FileDescriptorSet) (*Image, error) {
	image := Image{}
	image.Files = fd.File
	return &image, nil
}

func NewConfigFromFiles(files []*descriptorpb.FileDescriptorProto) (schema.Config, error) {
	config := v1.NewConfig()
	var fileNames []string
	for _, f := range files {
		fileNames = append(fileNames, f.GetName())
	}
	config.SetFiles(fileNames)
	return config, nil
}
