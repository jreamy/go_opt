package parse

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"regexp"
	"strings"
	"text/template"
)

// Package is a wrapping of the *ast.Package type that
// extracts useful information about the declared struct types
type Package struct {
	*ast.Package
	Structs []*Struct
}

// Struct is a wrapping of the *ast.StructType that extracts
// useful information about a declared struct
type Struct struct {
	*ast.StructType
	*ast.CommentGroup
	Name      string
	ShortName string
	Fields    []*Field
}

// Field is a wrapping of the *ast.Field type that extracts
// useful information about a struct field declaration
type Field struct {
	*ast.Field
	Name string
	Type *Type
	Tags map[string][]string

	TagName string // to be used by other packages
}

// Type contains information about a struct field type, either
// containing the name of the type or its struct definition
type Type struct {
	ast.Expr
	Name   string
	Struct *Struct
	Func   string // to be used by other packages
}

// Parse searches the given ast package for struct
// definitions with comments matching the given tag
func (p *Package) Parse(tag *regexp.Regexp) {

	// Walk the file, appending tagged struct definitions
	ast.Inspect(p.Package, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.GenDecl:
			if str := ParseStruct(x, tag); str != nil {
				p.Structs = append(p.Structs, str)
			}
		}
		return true
	})
}

// ParseStruct searches the given general declaration for
// a struct with a comment matching the given tag, returns
// nil if the declaration is not a struct or if the struct
// does not have a comment matching the tag regexp
func ParseStruct(g *ast.GenDecl, tag *regexp.Regexp) *Struct {
	var isTagged = false
	var s Struct

	// Walk the declaration, looking for a comment and struct definition
	ast.Inspect(g, func(n ast.Node) bool {

		switch x := n.(type) {
		case *ast.CommentGroup:
			s.CommentGroup = x

		case *ast.Comment:
			isTagged = isTagged || tag.Match([]byte(x.Text))

		case *ast.Ident:
			if isTagged {
				s.Name = x.Name
				if len(s.Name) != 0 {
					s.ShortName = strings.ToLower(x.Name[:1])
				}
			}

		case *ast.StructType:
			if isTagged {
				s.StructType = x
				s.Fields = ParseFields(x.Fields)
				return false
			}
		}
		return true
	})

	if !isTagged {
		return nil
	}

	return &s
}

// structDecl is a template for a struct declaration
var structDecl = template.Must(template.New("struct_declaration").Parse(`{{with .Name}}type {{.}}{{end}} struct {
	{{range .Fields}}{{.Declaration}}
	{{end}}
}`))

// Declaration recreates a declaration for a struct type
func (s *Struct) Declaration() string {
	if s == nil {
		return ""
	}

	var buf bytes.Buffer
	if err := structDecl.Execute(&buf, s); err != nil {
		panic(err)
	}

	if s.Name == "" {
		return buf.String()
	}

	src, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}

	return string(src)
}

// Declaration recreates a declaration for a struct field
func (f *Field) Declaration() string {
	if f == nil {
		return ""
	}

	var tag []string
	for key, vals := range f.Tags {
		tag = append(tag, fmt.Sprintf(`%s:"%s"`, key, strings.Join(vals, ",")))
	}

	return fmt.Sprintf("%s %s `%s`", f.Name, f.Type.Declaration(), strings.Join(tag, " "))
}

// Declaration recreates a type declaration for a struct field
func (t *Type) Declaration() string {
	if t == nil {
		return ""
	}

	if t.Name != "" {
		return t.Name
	}

	return t.Struct.Declaration()
}

// ParseFields searches the field list for exported fields
func ParseFields(fields *ast.FieldList) (fs []*Field) {
	if fields == nil {
		return
	}

	for _, f := range fields.List {
		for _, n := range f.Names {
			fs = append(fs, &Field{
				Field: f,
				Name:  n.Name,
				Type:  ParseType(f.Type),
				Tags:  parseTag(f.Tag),
			})
		}
	}

	return
}

// parseTag splits the given tag into a map of identifiers to
// comma separated values
func parseTag(t *ast.BasicLit) map[string][]string {
	if t == nil {
		return nil
	}

	items := strings.Split(t.Value, " ")
	tags := make(map[string][]string, len(items))

	for _, i := range items {
		tag := strings.Split(strings.Trim(i, "`"), ":")
		if len(tag) != 2 {
			continue
		}

		key, vals := tag[0], tag[1]
		vals = strings.Trim(vals, `"`)

		tags[key] = strings.Split(vals, ",")
	}

	return tags
}

func ParseType(expr ast.Expr) *Type {
	var t Type
	t.Expr = expr

	ast.Inspect(expr, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.Ident:
			t.Name = x.Name
			return false
		case *ast.StructType:
			t.Struct = &Struct{
				StructType: x,
				Fields:     ParseFields(x.Fields),
				Name:       "",
			}
			return false
		}
		fmt.Printf("%T %+v\n", n, n)
		return true
	})

	return &t
}
