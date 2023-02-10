package convert

import (
	"github.com/goccy/go-json"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func JsonDecode(in []byte, out interface{}) error {
	return json.Unmarshal(in, out)
}
func JsonEncode(in interface{}) ([]byte, error) {
	return json.Marshal(in)
}

func ProtoMessage(input []byte, out proto.Message) error {
	return proto.UnmarshalOptions{
		DiscardUnknown: true,
	}.Unmarshal(input, out)
}

func ProtoJsonToMessage(input []byte, out proto.Message) error {
	return protojson.UnmarshalOptions{
		DiscardUnknown: true,
	}.Unmarshal(input, out)
}

func ProtoBytes(in proto.Message) ([]byte, error) {
	return proto.Marshal(in)
}

func ProtoToJsonBytes(in proto.Message) ([]byte, error) {
	m := protojson.MarshalOptions{
		Indent:          "  ",
		UseProtoNames:   true,
		EmitUnpopulated: true,
	}
	return m.Marshal(in)
}
