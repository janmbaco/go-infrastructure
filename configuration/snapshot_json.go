package configuration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

func decodeSnapshot[T any](data []byte) (T, error) {
	var value T
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&value); err != nil {
		return value, err
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		return value, fmt.Errorf("configuration must contain exactly one JSON document")
	}
	return value, nil
}
