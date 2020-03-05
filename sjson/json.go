package sjson

import (
	j "encoding/json"
	"io/ioutil"

	"github.com/syncfuture/go/u"
)

// Deserialize deserialize json to object pointer
func Deserialize(json string, objPtr interface{}) error {
	jb := []byte(json)

	err := j.Unmarshal(jb, objPtr)
	u.LogError(err)

	return err
}

// DeserializeFromFile deserialize json file to object pointer
func DeserializeFromFile(filename string, objPtr interface{}) error {
	jb, err := ioutil.ReadFile(filename)
	if u.LogError(err) {
		return err
	}

	err = j.Unmarshal(jb, objPtr)
	return err
}

// Serialize serialize object to json string
func Serialize(objPtr interface{}) (string, error) {
	jb, err := j.Marshal(objPtr)
	if u.LogError(err) {
		return "", err
	}
	return string(jb), err
}

// SerializeToFile serialize object to json file
func SerializeToFile(objPtr interface{}, filename string) error {
	jb, err := j.Marshal(objPtr)
	if u.LogError(err) {
		return err
	}
	err = ioutil.WriteFile(filename, jb, 0644)
	return err
}
