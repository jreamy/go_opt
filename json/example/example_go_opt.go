package example

import (
	"encoding/json"
	"strconv"
	"sync"
)

var goOptBasicPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 1024)
	},
}

func goOptBasicBuffer() []byte {
	return goOptBasicPool.Get().([]byte)[:0]
}

// GoOptRecycleBasic returns the byte slice to
// the pool of slices available, increasing memory efficiency
func GoOptRecycleBasic(b []byte) {
	goOptBasicPool.Put(b)
}

// Basic struct field names
var goOptBasicNumber = []byte("\"int\":")
var goOptBasicSmall = []byte("\"Small\":")
var goOptBasicLarge = []byte("\"-\":")
var goOptBasicText = []byte("\"txt\":")

// MarshalJSON is generated json optimization
func (b Basic) MarshalJSON() ([]byte, error) {

	// Get reusable buffer
	buf := goOptBasicBuffer()

	buf = append(buf, '{')

	// Write b.Number

	buf = append(buf, goOptBasicNumber...)
	buf = append(buf, []byte(strconv.Itoa(b.Number))...)
	buf = append(buf, ',')

	// Write b.Small

	buf = append(buf, goOptBasicSmall...)
	buf = append(buf, []byte(strconv.Itoa(int(b.Small)))...)
	buf = append(buf, ',')

	// Write b.Large

	buf = append(buf, goOptBasicLarge...)
	buf = append(buf, []byte(strconv.FormatUint(uint64(b.Large), 10))...)
	buf = append(buf, ',')

	// Write b.Text
	if len(b.Text) != 0 {
		buf = append(buf, goOptBasicText...)
		buf = append(buf, '"')
		buf = append(buf, []byte(b.Text)...)
		buf = append(buf, '"')
		buf = append(buf, ',')
	}

	// Close the struct definition
	if len(buf) == 1 {
		buf = append(buf, '}')
	} else {
		// overwrite the last comma
		buf[len(buf)-1] = '}'
	}

	return buf, nil
}

var goOptSubstructPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 1024)
	},
}

func goOptSubstructBuffer() []byte {
	return goOptSubstructPool.Get().([]byte)[:0]
}

// GoOptRecycleSubstruct returns the byte slice to
// the pool of slices available, increasing memory efficiency
func GoOptRecycleSubstruct(b []byte) {
	goOptSubstructPool.Put(b)
}

// Substruct struct field names
var goOptSubstructText = []byte("\"Text\":")
var goOptSubstructSub = []byte("\"sub\":")

// MarshalJSON is generated json optimization
func (s Substruct) MarshalJSON() ([]byte, error) {

	// Get reusable buffer
	buf := goOptSubstructBuffer()

	buf = append(buf, '{')

	// Write s.Text

	buf = append(buf, goOptSubstructText...)
	buf = append(buf, '"')
	buf = append(buf, []byte(s.Text)...)
	buf = append(buf, '"')
	buf = append(buf, ',')

	// Write s.Sub
	buf = append(buf, goOptSubstructSub...)
	if bytes, err := json.Marshal(s.Sub); err != nil {
		return nil, err
	} else {
		buf = append(buf, bytes...)
	}
	buf = append(buf, ',')

	// Close the struct definition
	if len(buf) == 1 {
		buf = append(buf, '}')
	} else {
		// overwrite the last comma
		buf[len(buf)-1] = '}'
	}

	return buf, nil
}
