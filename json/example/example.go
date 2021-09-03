package example

// go-opt: json
type Basic struct {
	Number int `json:"int"`
	Small  int16
	Large  uint32 `json:"-,"`
	Text   string `json:"txt,omitempty"`
}

// go-opt: json
type Substruct struct {
	Text string
	Sub  struct {
		Text string
		Num  int32
	} `json:"sub"`
}
