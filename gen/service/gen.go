package service

import (
	"github.com/myadamtest/adam/gen/gogoproto"
	"github.com/myadamtest/adam/gen/proto"
	"os"
	"text/template"
)

func Gen(fileProto ...*gogoproto.ProtoFile) (err error) {
	for _, f := range fileProto {
		for _, m := range f.Messages {
			err = gen(&m)
			if err != nil {
				return
			}
		}
	}
	return
}

func gen(message *proto.MessageElement) error {
	t, err := template.New("text").Funcs(template.FuncMap{"firstUpName": formName, "firstLowerName": formName}).
		Parse(serviceTpl)
	if err != nil {
		return err
	}
	return t.Execute(os.Stdout, message)
}

func formName(name string) string {
	//todo camel
	return name
}
