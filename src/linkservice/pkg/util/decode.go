package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func DecodeJSON(r *http.Request, obj interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(obj)
	if err != nil {
		if err, ok := err.(*json.UnmarshalTypeError); ok {
			return fmt.Errorf("incorrect field %s value", err.Field)
		}
		if err == io.EOF {
			return errors.New("body is empty")
		}

		return err
	}

	return nil
}
