package main

import (
	"bytes"
	"html/template"

	"github.com/jreamy/go-opt/json/parse"
)

var intTemplate = template.Must(template.New("int_template").Parse(`
	// Write {{ .Str -}} . {{- .Fld }}
	{{ if .Omitempty }} if {{ .Str -}} . {{- .Fld }} != 0 { {{end}}
	{{ .Buf }} = append({{ .Buf }}, goOpt {{- .Lnm -}} {{- .Fld }}...)
	{{ .Buf }} = append({{ .Buf }}, []byte(strconv.Itoa(
		{{- if ne .Typ "int" -}}
		int( {{- .Str -}} . {{- .Fld -}} )
		{{- else -}}
		{{- .Str -}} . {{- .Fld -}}
		{{- end -}}
	 ))...)
	{{ .Buf }} = append( {{ .Buf }}, ',')
	{{ if .Omitempty }} } {{end }}
`))

var uintTemplate = template.Must(template.New("uint_template").Parse(`
	// Write {{ .Str -}} . {{- .Fld }}
	{{ if .Omitempty }} if {{ .Str -}} . {{- .Fld }} != 0 { {{end}}
	{{ .Buf }} = append({{ .Buf }}, goOpt {{- .Lnm -}} {{- .Fld }}...)
	{{ .Buf }} = append( {{- .Buf -}} , []byte(strconv.FormatUint(
		{{- if ne .Typ "uint64" -}}
		uint64( {{- .Str -}} . {{- .Fld -}} )
		{{- else -}}
		{{- .Str -}} . {{- .Fld -}}
		{{- end -}}
	, 10))...)
	{{ .Buf }} = append( {{ .Buf }}, ',')
	{{ if .Omitempty }} } {{end }}
`))

var strTemplate = template.Must(template.New("str_template").Parse(`
    // Write {{ .Str -}} . {{- .Fld }}
	{{ if .Omitempty }} if len( {{- .Str -}} . {{- .Fld -}} ) != 0 { {{end}}
	{{ .Buf }} = append({{ .Buf }}, goOpt {{- .Lnm -}} {{- .Fld }}...)
	{{ .Buf }} = append({{ .Buf }}, '"')
	{{ .Buf }} = append({{ .Buf }}, []byte( {{- .Str -}} . {{- .Fld -}} )...)
	{{ .Buf }} = append({{ .Buf }}, '"')
	{{ .Buf }} = append({{ .Buf }}, ',')
	{{ if .Omitempty }} } {{end }}
`))

var otherTemplate = template.Must(template.New("other_template").Parse(`
    // Write {{ .Str -}} . {{- .Fld }}
	{{ .Buf }} = append({{ .Buf }}, goOpt {{- .Lnm -}} {{- .Fld }}...)
	if bytes, err := json.Marshal( {{- .Str -}} . {{- .Fld -}} ); err != nil {
		return nil, err
	} else {
	 	{{ .Buf }} = append( {{ .Buf }} , bytes...)
	}
	{{ .Buf }} = append({{ .Buf }}, ',')
`))

type jsonTemplate struct {
	Buf string // buffer name
	Lnm string // struct long name
	Str string // struct short name
	Fld string // field name
	Typ string // field type name

	Omitempty bool // for wrapping omitempty checks
	StrCast   bool // for wrapping output in quotes
}

func JSONMarshalers(pkg *parse.Package, w string) {
	for _, s := range pkg.Structs {
		for _, f := range s.Fields {
			var buf bytes.Buffer
			tags := f.Tags["json"]
			v := jsonTemplate{
				Buf:       w,
				Lnm:       s.Name,
				Str:       s.ShortName,
				Fld:       f.Name,
				Typ:       f.Type.Name,
				Omitempty: len(tags) == 2 && tags[1] == "omitempty",
				StrCast:   len(tags) == 2 && tags[1] == "string",
			}

			if len(tags) > 0 && len(tags[0]) > 0 {
				f.TagName = tags[0]
			} else {
				f.TagName = f.Name
			}

			switch f.Type.Name {
			case "int", "int8", "int16", "int32", "int64":
				intTemplate.Execute(&buf, v)
				f.Type.Func = buf.String()

			case "uint", "uint8", "uint16", "uint32", "uint64":
				uintTemplate.Execute(&buf, v)
				f.Type.Func = buf.String()

			case "string":
				strTemplate.Execute(&buf, v)
				f.Type.Func = buf.String()

			default:
				otherTemplate.Execute(&buf, v)
				f.Type.Func = buf.String()
			}
		}
	}
}
