package gogoproto

const entityTpl = `
syntax = "{{.Syntax}}";

package {{.PackageName}};
{{range $i,$m :=.Messages}}
message {{$m.Name}} {
	{{- range $i2,$field := $m.Fields}}
	{{getFileProtoType $field}} {{$field.Name}} = {{plus $i2}};
	{{- end}}
}

{{- end}}
`
