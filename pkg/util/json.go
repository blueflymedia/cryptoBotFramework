package util

import (
	"encoding/json"
	"io/ioutil"
)

func WriteJsonFile(p string, obj interface{}) error {
	out, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(p, out, 0644)
}
