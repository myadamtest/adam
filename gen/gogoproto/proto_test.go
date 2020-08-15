package gogoproto

import (
	"github.com/myadamtest/adam/gen/proto"
	"testing"
)

func TestGen(t *testing.T) {
	fs := make([]*ProtoFile, 1)

	syntax := "proto3"
	pack := "mmmmPackage"
	msgName := "AAA"
	fieldName := "name"

	f2 := "age"

	fs[0] = &ProtoFile{}
	fs[0].Syntax = syntax
	fs[0].PackageName = pack
	fs[0].Messages = make([]proto.MessageElement, 1)
	fs[0].Messages[0] = proto.MessageElement{}
	fs[0].Messages[0].Name = msgName
	fs[0].Messages[0].Fields = make([]proto.FieldElement, 2)
	fs[0].Messages[0].Fields[0] = proto.FieldElement{}
	fs[0].Messages[0].Fields[0].Name = fieldName
	fs[0].Messages[0].Fields[0].Type, _ = proto.NewScalarDataType("string")

	fs[0].Messages[0].Fields[1] = proto.FieldElement{
		Name: f2,
		Type: fs[0].Messages[0].Fields[0].Type,
	}

	Gen(fs...)
}
