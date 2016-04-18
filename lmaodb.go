package lmaodb

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
)

const (
	dataDir = "__data"
)

func getIDFilePath(model string) string {
	return path.Join(getModelDir(model), "__ID")
}

func getModelDir(model string) string {
	return path.Join(dataDir, "__"+model)
}

func getNextID(model string) (int64, error) {
	idFilePath := getIDFilePath(model)
	data, err := ioutil.ReadFile(idFilePath)
	if err != nil {
		return 0, err
	}
	id, err := strconv.Atoi(string(data))
	if err != nil {
		return 0, err
	}
	return int64(id), nil
}

func getRecordFilePath(id int64, model string) string {
	return path.Join(getModelDir(model), fmt.Sprintf("%010d", id))
}

func getType(r interface{}) string {
	t := reflect.TypeOf(r).String()
	return splitAndGetLast(t, ".")
}

func getTypeFromSlice(slice interface{}) string {
	t := reflect.TypeOf(slice).Elem().String()
	return splitAndGetLast(t, ".")
}

func readRecord(id int64, model string) ([]byte, error) {
	filePath := getRecordFilePath(id, model)
	return ioutil.ReadFile(filePath)
}

func setNextID(id int64, model string) error {
	idFilePath := getIDFilePath(model)
	data := []byte(fmt.Sprintf("%d", id))
	return ioutil.WriteFile(idFilePath, data, 0600)
}

func splitAndGetLast(str, sep string) string {
	parts := strings.Split(str, sep)
	return parts[len(parts)-1]
}

func writeRecord(id int64, model string, data []byte) error {
	filePath := getRecordFilePath(id, model)
	return ioutil.WriteFile(filePath, data, 0600)
}

func DeleteRecord(record interface{}) error {
	model := getType(record)
	ptrVal := reflect.ValueOf(record)
	if ptrVal.Kind() != reflect.Ptr {
		panic("record argument must be a pointer!")
	}
	structVal := reflect.Indirect(ptrVal)
	id := structVal.Field(0)
	if id.Int() == 0 {
		msg := fmt.Sprintf(
			"lmaodb: can't delete record that doesn't exist in %s table!",
			model,
		)
		return errors.New(msg)
	}
	data := []byte{}
	return writeRecord(id.Int(), model, data)
}

func GetRecord(id int64, record interface{}) error {
	model := getType(record)
	ptrVal := reflect.ValueOf(record)
	if ptrVal.Kind() != reflect.Ptr {
		panic("record argument must be a pointer!")
	}
	data, err := readRecord(id, model)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		msg := fmt.Sprintf(
			"lmaodb: record %d in table %s no longer exists!",
			id, model,
		)
		return errors.New(msg)
	}
	return json.Unmarshal(data, ptrVal.Interface())
}

func GetRecords(slice interface{}) error {
	model := getTypeFromSlice(slice)
	if reflect.ValueOf(slice).Kind() != reflect.Ptr {
		panic("slice argument must be a pointer!")
	}
	i := int64(1)
	jsonStrings := []string{}
	for {
		data, err := readRecord(i, model)
		if err != nil {
			break
		}
		if len(data) > 0 {
			jsonStrings = append(jsonStrings, string(data))
		}
		i++
	}
	allRecords := fmt.Sprintf("[%s]", strings.Join(jsonStrings, ","))
	return json.Unmarshal([]byte(allRecords), slice)
}

func InitModel(dummyRecord interface{}) error {
	model := getType(dummyRecord)
	dir := getModelDir(model)
	err := os.MkdirAll(dir, 0700)
	if err != nil {
		return err
	}
	idFilePath := getIDFilePath(model)
	_, err = os.Stat(idFilePath)
	if err != nil && os.IsNotExist(err) {
		err = ioutil.WriteFile(idFilePath, []byte("1"), 0600)
	}
	return err
}

func QueryRecords(slice interface{}, field string, value interface{}) error {
	err := GetRecords(slice)
	if err != nil {
		return err
	}
	sliceVal := reflect.Indirect(reflect.ValueOf(slice))
	len := sliceVal.Len()
	newSliceVal := reflect.MakeSlice(sliceVal.Type(), 0, 0)
	for i := 0; i < len; i++ {
		record := sliceVal.Index(i)
		fvalue := record.FieldByName(field)
		eq := false
		switch fvalue.Interface().(type) {
		case int:
			eq = int(fvalue.Int()) == value.(int)
		case int64:
			eq = fvalue.Int() == value.(int64)
		case string:
			eq = fvalue.String() == value.(string)
		}
		if eq {
			newSliceVal = reflect.Append(newSliceVal, record)
		}
	}
	sliceVal.Set(newSliceVal)
	return nil
}

func PutRecord(record interface{}) error {
	model := getType(record)
	ptrVal := reflect.ValueOf(record)
	if ptrVal.Kind() != reflect.Ptr {
		panic("record argument must be a pointer!")
	}
	var nextID int64
	var err error
	structVal := reflect.Indirect(ptrVal)
	id := structVal.Field(0)
	if id.Int() == 0 {
		nextID, err = getNextID(model)
		if err != nil {
			return err
		}
		id.SetInt(nextID)
	}
	data, err := json.Marshal(ptrVal.Interface())
	if err != nil {
		return err
	}
	err = writeRecord(id.Int(), model, data)
	if err != nil {
		return err
	}
	if nextID != 0 {
		nextID++
		err = setNextID(nextID, model)
		if err != nil {
			return err
		}
	}
	return nil
}

func PutRecords(slice interface{}) error {
	sliceVal := reflect.ValueOf(slice)
	len := sliceVal.Len()
	for i := 0; i < len; i++ {
		err := PutRecord(sliceVal.Index(i).Addr().Interface())
		if err != nil {
			return err
		}
	}
	return nil
}
