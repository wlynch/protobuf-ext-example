package main

import (
	"log"

	"github.com/golang/protobuf/descriptor"
	"github.com/golang/protobuf/proto"
	protobuf "github.com/golang/protobuf/protoc-gen-go/descriptor"

	"github.com/wlynch/protobuf-ext-example/message"
)

func main() {
	m := &message.MyMessage{}
	fd, md := descriptor.ForMessage(m)

	getField(md)
	removeFields(fd)
}

func getField(md *protobuf.DescriptorProto) {
	// Find all fields that have the given option.
	for _, field := range md.GetField() {
		opts := field.GetOptions()
		if opts == nil {
			continue
		}
		v, err := proto.GetExtension(field.GetOptions(), message.E_MyFieldOption)
		if err != nil {
			log.Fatal(err)
		}
		value, ok := v.(*string)
		if !ok {
			log.Fatalf("unexpected type for %s", message.E_MyFieldOption.Name)
		}
		log.Println("Field value:", *value)
	}
}

func removeFields(fd *protobuf.FileDescriptorProto) *protobuf.FileDescriptorProto {
	for _, m := range fd.GetMessageType() {
		// Removing things from lists is harder than it should be in Go.
		// I'm being lazy here.
		var newFields []*protobuf.FieldDescriptorProto
		for _, f := range m.GetField() {
			opts := f.GetOptions()
			if opts == nil {
				continue
			}
			v, err := proto.GetExtension(f.GetOptions(), message.E_MyFieldOption)
			if err != nil {
				log.Fatal(err)
			}
			// We'll just assume that the existance of the field means that we should
			// prune it.
			if v == nil {
				newFields = append(newFields, f)
			}
		}
		// Replace existing fields with pruned set.
		m.Field = newFields
	}
	return fd
}
