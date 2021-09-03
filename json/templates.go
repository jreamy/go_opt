package main

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"github.com/jreamy/go-opt/json/parse"
	"golang.org/x/tools/imports"
)

func UseTemplate(filename string, tmp *template.Template, pkg *parse.Package) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer f.Close()

	var buf bytes.Buffer
	if err = tmp.Execute(&buf, pkg); err != nil {
		fmt.Println(buf.String())
		return err
	}
	bytes, err := imports.Process("", buf.Bytes(), nil)
	if err != nil {
		fmt.Println(buf.String())
		fmt.Printf("%s\n", bytes)
	}

	f.Write(bytes)
	return err
}

var optTemplate = template.Must(template.New("go_opt").Parse(`package {{.Name}}


{{range $idx, $str := .Structs}}

var goOpt {{- .Name -}} Pool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 1024)
	},
}

func goOpt {{- .Name -}} Buffer() []byte {
	return goOpt {{- .Name -}} Pool.Get().([]byte)[:0]
}

// GoOptRecycle {{- .Name }} returns the byte slice to
// the pool of slices available, increasing memory efficiency
func GoOptRecycle {{- .Name -}} (b []byte) {
	goOpt {{- .Name -}} Pool.Put(b)
}

// {{ $str.Name }} struct field names {{range $i, $f := $str.Fields}}
var goOpt {{- $str.Name -}} {{- $f.Name}} = []byte("\" {{- $f.TagName -}} \":") {{end}}


// MarshalJSON is generated json optimization
func ( {{- .ShortName}} {{ .Name -}} ) MarshalJSON () ([]byte, error) {

	// Get reusable buffer
	buf := goOpt {{- .Name -}} Buffer()

	buf = append(buf, '{')

	{{range $idx, $f := .Fields}}
	{{$f.Type.Func}}
	{{end}}


	// Close the struct definition
	if len(buf) == 1 {
		buf = append(buf, '}')
	} else {
		// overwrite the last comma
		buf[len(buf) - 1] = '}'
	}

	return buf, nil
}
{{end}}

`))

var testTemplate = template.Must(template.New("go_opt_tests").Parse(`package {{.Name}}

import (
	"encoding/json"
	"testing"

	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/assert"
)

{{range .Structs}}
type {{.Name -}} Orig struct {
	{{range .Fields}}{{.Declaration}}
	{{end}}
}

func testStructs {{- .Name -}} () ({{- .Name }} , {{ .Name -}} Orig) {
	var orig {{.Name -}} Orig
	var opt {{.Name}}

	f := fuzz.New()
	f.Fuzz(&orig)

	{{ range .Fields }}
	opt. {{- .Name }} = orig. {{- .Name }}{{end}}

	return opt, orig
}

func TestMarshal {{- .Name -}} (t *testing.T) {
	opt, orig := testStructs {{- .Name -}} ()

	o, _ := opt.MarshalJSON()
	t.Logf("%s", o)

	origBytes, err := json.Marshal(orig)
	assert.NoError(t, err)
	optBytes, err := json.Marshal(opt)
	assert.NoError(t, err)

	assert.NotEmpty(t, origBytes)
	assert.NotEmpty(t, optBytes)

	assert.Equal(t, string(origBytes), string(optBytes))
}

func BenchmarkMarshal {{- .Name -}} (b *testing.B) {
	opt, orig := testStructs {{- .Name -}} ()

	b.Run("go-opt", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b, _ := opt.MarshalJSON()
			GoOptRecycle {{- .Name -}} (b)
		}
	})

	b.Run("json", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			json.Marshal(orig)
		}
	})
}

{{end}}
`))
