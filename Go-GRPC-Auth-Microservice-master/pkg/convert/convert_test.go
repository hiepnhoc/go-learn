package convert

import (
	"testing"
)

func TestDecode(t *testing.T) {

	inputString := `{"key1":"value1", "key2": "value2"}`

	type inputStruct struct {
		Key1 string `json:"key1"`
		Key2 string `json:"key2"`
	}

	if err := JsonDecode([]byte(inputString), &inputStruct{}); err != nil {
		t.Fatalf("TestDecode: JsonDecode failed %v", err)
	}

	if err := JsonDecode([]byte("inputString"), &inputStruct{}); err == nil {
		t.Fatalf("TestDecode: JsonDecode mush have throw a  JSON input error")
	}

}

func TestEncode(t *testing.T) {

	type inputStruct struct {
		Key1 string `json:"key1"`
		Key2 string `json:"key2"`
	}

	input := &inputStruct{
		Key1: "value1",
		Key2: "value2",
	}
	want := `{"key1":"value1","key2":"value2"}`

	encode, err := JsonEncode(input)
	if err != nil {
		return
	}

	if err != nil {
		t.Fatalf("TestDecode: JsonDecode failed %v", err)
	}

	got := string(encode)

	if got != want {
		t.Fatalf("TestDecode: got %v want %s", got, want)
	}

}

//func TestToMessageAndToBytes(t *testing.T) {
//
//	createTrueIdRequest := v1.CreateTrueIdRequest{Service: "test", ReferenceId: "----"}
//
//	bytes, err := ProtoBytes(&createTrueIdRequest)
//
//	if err != nil {
//		t.Fatalf("TestToMessageAndToBytes: 1 ProtoBytes failed %v", err)
//	}
//	createTrueIdRequest1 := v1.CreateTrueIdRequest{}
//	if err := ProtoMessage(bytes, &createTrueIdRequest1); err != nil {
//		t.Fatalf("TestToMessageAndToBytes: 1 ProtoMessage failed %v", err)
//	}
//
//	if createTrueIdRequest.Service != createTrueIdRequest1.Service || createTrueIdRequest.ReferenceId != createTrueIdRequest1.ReferenceId {
//		t.Fatalf("TestToMessageAndToBytes: failed")
//	}
//
//}

//func TestToJsonMessageAndToBytes(t *testing.T) {
//
//	var (
//		service     = "test"
//		referenceId = "----"
//	)
//
//	type JsonTest struct {
//		Service     string `json:"service"`
//		ReferenceId string `json:"reference_id"`
//	}
//	json := JsonTest{
//		Service:     service,
//		ReferenceId: referenceId,
//	}
//
//	jsonBytes, err := JsonEncode(json)
//
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	createTrueIdRequest := v1.CreateTrueIdRequest{}
//
//	err = ProtoJsonToMessage(jsonBytes, &createTrueIdRequest)
//
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	if createTrueIdRequest.Service != service || createTrueIdRequest.ReferenceId != referenceId {
//		t.Fatalf("TestToJsonMessageAndToBytes : service, reference_id got %s,%s , service, reference_id want %s,%s ", createTrueIdRequest.Service, createTrueIdRequest.ReferenceId, service, referenceId)
//
//	}
//
//	t.Log(createTrueIdRequest)
//
//	jsonBytes, err = ProtoToJsonBytes(&createTrueIdRequest)
//
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	endTest := JsonTest{}
//
//	err = JsonDecode(jsonBytes, &endTest)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	t.Log(endTest)
//
//	if endTest.Service != service || endTest.ReferenceId != referenceId {
//		t.Fatalf("TestToJsonMessageAndToBytes : service, reference_id got %s,%s , service, reference_id want %s,%s ", endTest.Service, endTest.ReferenceId, service, referenceId)
//
//	}
//
//}
