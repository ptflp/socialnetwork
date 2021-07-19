package decoder

import (
	"bytes"
	"io"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Decoder struct {
}

func NewDecoder() *Decoder {
	return &Decoder{}
}

func (d *Decoder) Decode(r io.Reader, val interface{}) error {
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(val); err != nil {
		return err
	}

	return nil
}

func (d *Decoder) Encode(w io.Writer, value interface{}) error {
	return json.NewEncoder(w).Encode(value)
}

func (d *Decoder) MapStructs(dest, src interface{}) error {
	var b bytes.Buffer

	err := d.Encode(&b, src)
	if err != nil {
		return err
	}

	err = d.Decode(&b, dest)
	if err != nil {
		return err
	}

	return nil
}
