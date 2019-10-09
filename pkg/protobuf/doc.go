package protobuf

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	docreq "github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/jhump/protoreflect/desc"
	docgen "github.com/pseudomuto/protoc-gen-doc"
	"github.com/tinyzimmer/protobuf-registry/pkg/config"
)

func (p *Protobuf) DocJSON(filename string) ([]byte, error) {
	descriptors, err := p.GetDescriptors()
	if err != nil {
		return nil, err
	}
	//defer os.RemoveAll(tempPath)
	var desc *desc.FileDescriptor
	for _, x := range descriptors {
		if x.GetName() == filename || x.GetName() == strings.TrimPrefix(filename, "/") {
			desc = x
			break
		}
	}
	if desc == nil {
		return nil, fmt.Errorf("No file %s in this protobuf package", filename)
	}
	return generateDocs(desc, descriptors)
}

func generateDocs(desc *desc.FileDescriptor, descriptors map[string]*desc.FileDescriptor) ([]byte, error) {
	plugin := &docgen.Plugin{}
	req := &docreq.CodeGeneratorRequest{
		FileToGenerate:  []string{desc.GetName()},
		ProtoFile:       toFileDescriptorProtos(descriptors),
		CompilerVersion: parseProtocVersion(),
		Parameter:       strPtr("json,docs.json"),
	}
	res, err := plugin.Generate(req)
	if err != nil {
		return nil, err
	} else if len(res.File) == 0 {
		return nil, fmt.Errorf("No documentation returned from the plugin")
	}
	content := *res.File[0].Content
	out := []byte(content)
	return out, nil
}

func intPtr(s string) *int32 {
	i, _ := strconv.Atoi(s)
	i32 := int32(i)
	return &i32
}

func strPtr(s string) *string {
	return &s
}

func parseProtocVersion() *docreq.Version {
	spl := strings.Split(config.GlobalConfig.ProtobufVersion, " ")
	vers := spl[len(spl)-1]
	versSplit := strings.Split(vers, ".")
	return &docreq.Version{
		Major: intPtr(versSplit[0]),
		Minor: intPtr(versSplit[1]),
		Patch: intPtr(versSplit[2]),
	}
}

func toFileDescriptorProtos(descriptors map[string]*desc.FileDescriptor) []*descriptor.FileDescriptorProto {
	var rawDescriptors []*descriptor.FileDescriptorProto
	for _, x := range descriptors {
		rawDescriptors = append(rawDescriptors, x.AsFileDescriptorProto())
	}
	return rawDescriptors
}
