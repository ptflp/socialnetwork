package cache

import (
	"bytes"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func Serialize(value interface{}) ([]byte, error) {
	if data, ok := value.([]byte); ok {
		return data, nil
	}

	var b bytes.Buffer
	encoder := json.NewEncoder(&b)
	if err := encoder.Encode(value); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func Deserialize(byt []byte, ptr interface{}) error {
	if data, ok := ptr.(*[]byte); ok {
		*data = byt
		return nil
	}

	b := bytes.NewBuffer(byt)
	decoder := json.NewDecoder(b)
	if err := decoder.Decode(ptr); err != nil {
		return err
	}

	return nil
}
