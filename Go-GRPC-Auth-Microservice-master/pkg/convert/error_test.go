package convert

import (
	"errors"
	"github.com/google/uuid"
	"testing"
)

func TestErrorConvert(t *testing.T) {

	myError := errors.New("TestErrorConvert " + uuid.NewString())

	errorConvert := NewError("TestErrorConvert", "JsonDecode", "123", "456", myError)

	t.Log(errorConvert.Error())

	type MyStructIn struct {
		field1 string
		field2 int
	}

	type MyStructOut struct {
		field1 string
		field2 int
	}

	myError = errors.New("TestErrorConvert" + uuid.NewString())

	errorConvert = NewError("TestErrorConvert", "JsonDecode", MyStructIn{field1: "1", field2: 2}, MyStructOut{}, myError)

	t.Log(errorConvert.Error())

	errorConvert = NewError("TestErrorConvert", "ProtoMessage", 123123, MyStructOut{}, myError)

	t.Log(errorConvert.Error())

	var outBytes []byte
	errorConvert = NewError("TestErrorConvert", "ProtoMessage", 123123, outBytes, myError)

	t.Log(errorConvert.Error())
}
