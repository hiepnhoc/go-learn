package repository

import (
	"2margin.vn/account-service/config"
	"2margin.vn/account-service/pkg/logger"
	"context"
	dapr "github.com/dapr/go-sdk/client"
	"sync"
	"testing"
)

var client = newFakeClient()

func newFakeClient() *fakeClient {
	return &fakeClient{}
}

func TestDatabase(t *testing.T) {

	type TestSpec struct {
		Id   string `db:"id"`
		Name string `db:"name"`
	}

	testSpec := TestSpec{Id: "_"}

	columnQueries, _ := createCols(&testSpec)

	t.Log(columnQueries)
	t.Log(columnQueries.columns())

	if len(columnQueries.columns()) != 2 {
		t.Fatalf("len is 2")
	}

	values := []interface{}{"id-raw", "name-raw"}

	err := scanCols(columnQueries, values)
	if err != nil {
		t.Fatalf("Scan err %v", err)
	}

	if testSpec.Id != "id-raw" {
		t.Fatalf("len is id-raw")
	}

}

type fakeEntity struct {
	Col1 string `db:"col1"`
	Col2 string `db:"col2"`
}

type fakeClient struct {
	dapr.Client
	sync.Mutex
	rawNumber int
}

func (f *fakeClient) ChangeRawNumber(number int) {
	f.Lock()
	f.rawNumber = number
	f.Unlock()
}

func (f *fakeClient) InvokeBinding(ctx context.Context, in *dapr.InvokeBindingRequest) (out *dapr.BindingEvent, err error) {
	if in.Operation == query {
		var data [][]interface{}

		for i := 0; i < f.rawNumber; i++ {
			data = append(data, []interface{}{"col1", "col2"})
		}

		dateBytes, _ := encode(data)
		return &dapr.BindingEvent{Data: dateBytes}, nil
	}

	return &dapr.BindingEvent{}, nil
}

func TestLoadEntity(t *testing.T) {

	var log = logger.NewAppLogger(&config.Config{})
	log.InitLogger()
	var rep = NewRepository(log, &Config{}, client)

	entity := &fakeEntity{}

	ok, err := rep.LoadEntity(context.Background(), "test", entity, nil, nil)

	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatalf("TestLoadEntity : LoadEntity want to not found")
	}

	client.ChangeRawNumber(1)

	ok, err = rep.LoadEntity(context.Background(), "test", entity, nil, nil)

	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fatalf("TestLoadEntity : LoadEntity want to found")
	}

}

func TestExec(t *testing.T) {

	var log = logger.NewAppLogger(&config.Config{})
	log.InitLogger()

	var rep = NewRepository(log, &Config{}, client)

	err := rep.Exec(context.Background(), "select * from dual")

	if err != nil {
		t.Fatal(err)
	}
}
