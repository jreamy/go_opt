package parse

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePackage(t *testing.T) {
	prefix := `package etc

	`

	match := `
// Etc is a tagged struct
// trigger
type Etc struct {
	thing string
}`

	function := `
// Etc is a tagged function
// trigger
func Etc() {}`

	unlabeled := `
// Etc is an untagged struct
type Etc struct {
	thing string
}`

	nested := `
// Etc is a tagged struct
// trigger
type Etc struct {
	thing struct {
		thing2 string
	}
}`

	tag := regexp.MustCompile(`// trigger`)

	cases := []struct {
		name    string
		src     string
		trigger int
	}{
		{"simple match", match, 1},
		{"two structs", match + strings.ReplaceAll(match, "Etc", "Etc2"), 2},
		{"function", function, 0},
		{"unlabeled struct", unlabeled, 0},
		{"nested structs", nested, 1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			src := prefix + tc.src

			// Mock file Parse
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "", src, parser.ParseComments)
			assert.NoError(t, err)

			// Mock package parse
			pkg := &ast.Package{
				Files: map[string]*ast.File{
					tc.name + ".go": f,
				},
			}

			p := &Package{Package: pkg}
			p.Parse(tag)

			assert.Equal(t, tc.trigger, len(p.Structs))

			for _, s := range p.Structs {
				fmt.Println(s.Declaration())
			}
		})
	}
}
