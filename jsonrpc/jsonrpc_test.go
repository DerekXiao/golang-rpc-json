package jsonrpc

import (
	"fmt"
	"reflect"
	"testing"
)

func TestDecodeJson(t *testing.T) {
	var i = []byte(`$[{"Id":123,"MethodName":"abc","Param":[{"ParamType":"int8","ParamValue":9}]}]#`)
	k, err := decodeJsonArray(i)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(k)

	value, err := analysisParam(k[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(value)
}

func TestTransfer(t *testing.T) {
	bytes := make([]byte, 0, 512)
	var i = []byte(`$[{"Id":123,"MethodName":"swap","Param":[{"ParamType":"int","ParamValue":9},{"ParamType":"int","ParamValue":8}]}]#`)
	bytes = append(bytes, '$')
	bytes = append(bytes, i...)
	bytes = append(bytes, '#')

}

func TestHandler(t *testing.T) {
	var i = []byte(`$[{"Id":123,"MethodName":"swap","Param":[{"ParamType":"int","ParamValue":9},{"ParamType":"int","ParamValue":8}]}]#`)
	jsonrs := new(JsonRpcServer)
	var ress []Response
	// Swap := func(a int, b int) (int, int) {
	// 	a, b = b, a
	// 	return a, b
	// }
	jsonrs.RegisteFunc("swap", func(a int, b int) (int, int) {
		a, b = b, a
		return a, b
	})
	reqs, err := decodeJsonArray(i)
	fmt.Println("reqs= ", reqs, err)
	if err != nil {
		res := new(Response)
		res.Error = "error"
		ress = append(ress, *res)
		fmt.Println("err ress1 = ", ress)
		return
	}

	for i := 0; i < len(reqs); i++ {
		fmt.Println("req[i]=", reqs[i])
		values, err := analysisParam(reqs[i])
		if err != nil {
			res := new(Response)
			res.Error = err.Error()
			ress = append(ress, *res)
			fmt.Println("err ress2 = ", ress)
			return
		}

		results := reflect.ValueOf(jsonrs.FunctionMap[reqs[i].MethodName]).Call(values)
		res := new(Response)
		res.MethodName = reqs[i].MethodName
		for _, result := range results {
			param := new(Parameter)
			param.ParamType = result.Kind().String()
			param.ParamValue = result.Interface()
			res.Param = append(res.Param, *param)
		}
		ress = append(ress, *res)
	}
	fmt.Println(ress)
}

func TestServer(t *testing.T) {
	rpcserver := new(JsonRpcServer)

	//set the port you want to listen with a int-type-parameter

	rpcserver.SetPort(7777)

	//set the protocol you want to use, but now it just supports tcp4
	//you can use "tcp" or "tcp4".

	rpcserver.SetProtocol("tcp")

	//registe your function into this server .The type of function's
	//parameter , both input and return , should be in this set {int,int8,
	//in16,int32,int64,uint,uint8,uint16,uint32,uint64,string,bool}

	rpcserver.RegisteFunc("swap", func(a int, b int) (int, int) {
		a, b = b, a
		return a, b
	})

	//start this server

	rpcserver.Serve()
}
