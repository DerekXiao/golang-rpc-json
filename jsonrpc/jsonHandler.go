package jsonrpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

const (
	UINT    uint    = 0
	INT     int     = 0
	FLOAT32 float32 = 0
	FLOAT64 float64 = 0
	UINT8   uint8   = 0
	UINT16  uint16  = 0
	UINT32  uint32  = 0
	UINT64  uint64  = 0
	INT8    int8    = 0
	INT16   int16   = 0
	INT32   int32   = 0
	INT64   int64   = 0
	STRING  string  = ""
	BOOL    bool    = false
)

//Decode jsonArray into the Request struct.
func decodeJsonArray(jsonArray []byte) ([]Request, error) {
	var r []Request
	l := len(jsonArray)
	if jsonArray[0] != '$' || jsonArray[l-1] != '#' {
		return nil, errors.New("invalid json data")
	}
	newJsonArray := jsonArray[1 : l-1]
	err := json.Unmarshal(newJsonArray, &r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

//Encode Response struct to jsonArray
func encodeJsonArray(resps []Response) []byte {
	b, err := json.Marshal(resps)
	if err != nil {
		fmt.Println(err) //need to change it in future
		return nil
	}
	return b
}

//Put each parameter which is in the specific type into a []reflect.Value
// for function's calling.
func analysisParam(req Request) ([]reflect.Value, error) {
	value := make([]reflect.Value, len(req.Param))
	var err error = nil
	defer func() {
		if r := recover(); r != nil {
			value = nil
			err = errors.New("ParamType doesn't match the Value!")
		}
	}()
	for i, v := range req.Param {
		switch v.ParamType {
		case "uint":
			value[i] = reflect.ValueOf(v.ParamValue).Convert(reflect.TypeOf(UINT))
		case "int":
			value[i] = reflect.ValueOf(v.ParamValue).Convert(reflect.TypeOf(INT))
		case "int8":
			value[i] = reflect.ValueOf(v.ParamValue).Convert(reflect.TypeOf(INT8))
		case "int16":
			value[i] = reflect.ValueOf(v.ParamValue).Convert(reflect.TypeOf(INT16))
		case "int32":
			value[i] = reflect.ValueOf(v.ParamValue).Convert(reflect.TypeOf(INT32))
		case "int64":
			value[i] = reflect.ValueOf(v.ParamValue).Convert(reflect.TypeOf(INT64))
		case "uint8":
			value[i] = reflect.ValueOf(v.ParamValue).Convert(reflect.TypeOf(UINT8))
		case "uint16":
			value[i] = reflect.ValueOf(v.ParamValue).Convert(reflect.TypeOf(UINT16))
		case "uint32":
			value[i] = reflect.ValueOf(v.ParamValue).Convert(reflect.TypeOf(UINT32))
		case "uint64":
			value[i] = reflect.ValueOf(v.ParamValue).Convert(reflect.TypeOf(UINT64))
		case "string":
			value[i] = reflect.ValueOf(v.ParamValue).Convert(reflect.TypeOf(STRING))
		case "bool":
			value[i] = reflect.ValueOf(v.ParamValue).Convert(reflect.TypeOf(BOOL))
		case "float32":
			value[i] = reflect.ValueOf(v.ParamValue).Convert(reflect.TypeOf(FLOAT32))
		case "float64":
			value[i] = reflect.ValueOf(v.ParamValue).Convert(reflect.TypeOf(FLOAT64))
		default:
			return nil, errors.New("ParamType is not supported!")

		}
	}
	return value, err
}
