package example

import (
	"encoding/json"
	"testing"

	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/assert"
)

type BasicOrig struct {
	Number int    `json:"int"`
	Small  int16  ``
	Large  uint32 `json:"-,"`
	Text   string `json:"txt,omitempty"`
}

func testStructsBasic() (Basic, BasicOrig) {
	var orig BasicOrig
	var opt Basic

	f := fuzz.New()
	f.Fuzz(&orig)

	opt.Number = orig.Number
	opt.Small = orig.Small
	opt.Large = orig.Large
	opt.Text = orig.Text

	return opt, orig
}

func TestMarshalBasic(t *testing.T) {
	opt, orig := testStructsBasic()

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

func BenchmarkMarshalBasic(b *testing.B) {
	opt, orig := testStructsBasic()

	b.Run("go-opt", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b, _ := opt.MarshalJSON()
			GoOptRecycleBasic(b)
		}
	})

	b.Run("json", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			json.Marshal(orig)
		}
	})
}

type SubstructOrig struct {
	Text string ``
	Sub  struct {
		Text string ``
		Num  int32  ``
	} `json:"sub"`
}

func testStructsSubstruct() (Substruct, SubstructOrig) {
	var orig SubstructOrig
	var opt Substruct

	f := fuzz.New()
	f.Fuzz(&orig)

	opt.Text = orig.Text
	opt.Sub = orig.Sub

	return opt, orig
}

func TestMarshalSubstruct(t *testing.T) {
	opt, orig := testStructsSubstruct()

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

func BenchmarkMarshalSubstruct(b *testing.B) {
	opt, orig := testStructsSubstruct()

	b.Run("go-opt", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b, _ := opt.MarshalJSON()
			GoOptRecycleSubstruct(b)
		}
	})

	b.Run("json", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			json.Marshal(orig)
		}
	})
}
