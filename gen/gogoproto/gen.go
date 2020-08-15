package gogoproto

import (
	"github.com/myadamtest/adam/gen/proto"
	"os"
	"text/template"
)

// rollback
func Gen(fileDescriptorProtoList ...*ProtoFile) error {
	for _, f := range fileDescriptorProtoList {
		err := gen(f)
		if err != nil {
			return err
		}
	}
	return nil
}

func gen(fileDescriptorProto *ProtoFile) error {
	t, err := template.New("text").Funcs(template.FuncMap{"getFileProtoType": getFileProtoType, "plus": plus}).
		Parse(entityTpl)
	if err != nil {
		return err
	}
	return t.Execute(os.Stdout, fileDescriptorProto)
}

func getFileProtoType(field *proto.FieldElement) string {
	return field.Type.Name()
}

func plus(i int) int {
	return i + 1
}
