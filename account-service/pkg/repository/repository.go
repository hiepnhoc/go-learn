package repository

import (
	"2margin.vn/account-service/pkg/logger"
	"context"
	"fmt"
	dapr "github.com/dapr/go-sdk/client"
	"github.com/doug-martin/goqu/v9"
	"github.com/goccy/go-json"
	"github.com/gookit/goutil/strutil"

	"github.com/doug-martin/goqu/v9/exp"
	"github.com/gookit/goutil/errorx"
	"reflect"
	"time"
)

type Config struct {
	Name string `yaml:"name"`
}

const (
	exec  = "exec"
	query = "query"
)

const (
	sqlKey = "sql"
)

// https://github.com/doug-martin/goqu/issues/175

type Repository interface {
	LoadEntity(ctx context.Context, table string, entity interface{}, where goqu.Expression, order exp.OrderedExpression) (bool, error)
	Exec(ctx context.Context, exec string) error
}

type repository struct {
	log    logger.Logger
	config *Config
	client dapr.Client
}

func NewRepository(log logger.Logger, cfg *Config, client dapr.Client) Repository {
	return &repository{log, cfg, client}
}

type columnInfos []colInfo

func (c columnInfos) columns() []interface{} {
	columns := make([]interface{}, len(c))
	for i, x := range c {
		columns[i] = x.db
	}
	return columns
}

type colInfo struct {
	value reflect.Value
	db    string
	index int
}

func (r *repository) LoadEntity(ctx context.Context, table string, entity interface{}, where goqu.Expression, order exp.OrderedExpression) (bool, error) {

	var (
		limit uint = 1
	)

	cols, err := createCols(entity)

	if err != nil {
		return false, errorx.Errorf("repository LoadEntity createCols entity %s failed - %v", typeOfObject(entity), err)
	}

	dataset := goqu.From(table).Select(cols.columns()...).Limit(limit)

	if where != nil {
		dataset.Where(where)
	}

	if order != nil {
		dataset.Order(order)
	}

	sql, _, err := dataset.ToSQL()

	if err != nil {
		return false, errorx.Errorf("repository LoadEntity get sql query table %v failed - %v ", table, err)
	}

	r.log.Debugf("Debug LoadEntity table %v : sql is %v", table, sql)

	metadata := map[string]string{
		sqlKey: sql,
	}

	invokeRequest := &dapr.InvokeBindingRequest{
		Name:      r.config.Name,
		Operation: query,
		Metadata:  metadata,
	}

	response, err := r.client.InvokeBinding(ctx, invokeRequest)

	r.log.Debugf("Debug LoadEntity table %v : sql response %v", table, response)

	var raws [][]interface{}

	err = decode(response.Data, &raws)

	if err != nil {
		return false, errorx.Errorf("repository LoadEntity decode %v %v failed - %v", response.Data, raws, err)
	}

	if len(raws) == 0 {
		return false, nil
	}

	err = scanCols(cols, raws[0])

	if err != nil {
		return false, errorx.Errorf("repository LoadEntity scanCols %v %v failed - %v", cols, raws[0], err)
	}

	return true, nil

}

func (r *repository) Exec(ctx context.Context, sql string) error {

	r.log.Infof("repository Exec : sql insert is %s", sql)

	metadata := map[string]string{
		sqlKey: sql,
	}

	invokeRequest := &dapr.InvokeBindingRequest{
		Name:      r.config.Name,
		Operation: exec,
		Metadata:  metadata,
	}

	_, err := r.client.InvokeBinding(ctx, invokeRequest)

	if err != nil {
		return errorx.Errorf("Repository Exec sql %s failed - %v ", sql, err)
	}

	return nil
}

func createCols(spec interface{}) (columnInfos, error) {
	const TagVal = "db"

	sType := reflect.TypeOf(spec)
	sVal := reflect.ValueOf(spec)

	if sType.Kind() == reflect.Ptr {
		sType = sType.Elem()
	}

	if sVal.Kind() == reflect.Ptr {
		sVal = sVal.Elem()
	}

	if sType.Kind() != reflect.Struct {
		return nil, errorx.Errorf("database createCols is not a valid pointer to struct")
	}

	var columns columnInfos

	for index := 0; index < sType.NumField(); index++ {
		fieldInfo := sType.Field(index)
		fieldVal := sVal.FieldByName(fieldInfo.Name)

		name, ok := fieldInfo.Tag.Lookup(TagVal)

		if !fieldVal.CanSet() || !ok {
			continue
		}

		column := colInfo{
			value: fieldVal,
			db:    name,
			index: index,
		}

		columns = append(columns, column)
		columns[index] = column
	}

	return columns, nil

}

func scanCols(columns columnInfos, raw []interface{}) error {
	defer func() error {
		if err := recover(); err != nil {
			return errorx.Errorf("repository scanCols column error : %s", err)
		}
		return nil
	}()

	var timeKind = reflect.TypeOf(time.Time{}).Kind()

	if len(columns) != len(raw) {
		return errorx.Errorf("repository scanCols is not a valid number len(columns) = %v , len(raw) = %v", len(columns), len(raw))
	}

	for _, column := range columns {

		switch column.value.Type().Kind() {
		case timeKind:
			var toTime, err = toTime(raw[column.index].(string), time.RFC3339)
			if err != nil {
				return errorx.Errorf("repository scanCols column %s toTime.parse failed : %s", column.db, err)
			}
			column.value.Set(reflect.ValueOf(toTime))
			break
		case reflect.Slice:
			if reflect.TypeOf(raw[column.index]).Kind() == reflect.Map {
				encode, err := encode(raw[column.index])
				if err != nil {
					return err
				}
				column.value.Set(reflect.ValueOf(encode))
			}
		default:
			column.value.Set(reflect.ValueOf(raw[column.index]).Convert(column.value.Type()))
		}

	}

	return nil

}

func toTime(s string, layouts ...string) (t time.Time, err error) {
	/*
		https://stackoverflow.com/questions/40939261/parse-strange-date-format
		https://go.dev/play/p/VO5413Z7-z
	*/
	return strutil.ToTime(s, layouts...)
}

func decode(in []byte, out interface{}) error {
	return json.Unmarshal(in, out)
}

func encode(in interface{}) ([]byte, error) {
	return json.Marshal(in)
}

func typeOfObject(x interface{}) string {
	return fmt.Sprint(reflect.TypeOf(x))
}
