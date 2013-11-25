package jsonrpc

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net"
	"os"
	"reflect"
	"strconv"
)


type JsonRpcServer struct {
	Port        string
	Protocol    string
	FunctionMap map[string]interface{}
}

type Parameter struct {
	ParamType  string
	ParamValue interface{}
}

type Request struct {
	ID         int
	MethodName string
	Param      []Parameter
}

type Response struct {
	ID         int
	MethodName string
	Error      string
	Param      []Parameter
}

//set the port you want to use
func (this *JsonRpcServer) SetPort(port int) {
	a := strconv.Itoa(port)
	this.Port = ":" + a
}

//set the protocol you want to use
func (this *JsonRpcServer) SetProtocol(protocol string) {
	this.Protocol = protocol
}

//Registe your function into this JsonRpcServer.All you need to take notice of is that
//the first param should be a string name hadn`t be used while the second param is
//a function type.
func (this *JsonRpcServer) RegisteFunc(name string, RegFunc interface{}) error {
	if this.FunctionMap == nil {
		this.FunctionMap = make(map[string]interface{})
	}
	if _, ok := this.FunctionMap[name]; ok {
		return errors.New("this name is already used , choose another name please!")
	}
	if rf := reflect.TypeOf(RegFunc); rf.Kind().String() != "func" {
		return errors.New("you need to registe a function,not any other type.")
	}
	this.FunctionMap[name] = RegFunc
	return nil
}

//this function should be used at the last step of your main function to make sure all those
//attributes neccessary you have already set .
func (this *JsonRpcServer) Serve() {
	switch this.Protocol {
	case "tcp", "tcp4":
		this.serveTcp()
	default:
		fmt.Fprint(os.Stderr, "Error occured : unexpected protocol")
		os.Exit(1)
	}
}

// TCPServer
func (this *JsonRpcServer) serveTcp() {
	tcpAddr, err := net.ResolveTCPAddr(this.Protocol, this.Port)
	checkErr(err)
	listener, err := net.ListenTCP("tcp4", tcpAddr)
	checkErr(err)
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			continue
		}
		go this.requestHandler(conn)
	}

}

func (this *JsonRpcServer) requestHandler(conn net.Conn) {
	defer conn.Close()
	defer func() {
		if rec := recover(); rec != nil {
			res := &Response{}
			res.Error = "panic occured in function call!"
			ress := make([]Response, 1)
			ress[0] = *res
			conn.Write(encodeJsonArray(ress))
		}
	}()
	fmt.Println("new client!Address is : ", conn.RemoteAddr())
	var ress []Response
	reader := bufio.NewReader(conn)
	jsonByte, err := reader.ReadBytes('#')
	bytesbuf := bytes.NewBuffer(jsonByte)
	fmt.Println("data is :", bytesbuf.String()) //delete!
	if err != nil {
		fmt.Println(err)
		res := new(Response)
		res.Error = err.Error()
		ress = append(ress, *res)
		conn.Write(encodeJsonArray(ress))
		return
	}
	reqs, err := decodeJsonArray(jsonByte)
	if err != nil {
		res := new(Response)
		res.Error = err.Error()
		ress = append(ress, *res)
		conn.Write(encodeJsonArray(ress))
		return
	}
	for i := 0; i < len(reqs); i++ {
		values, err := analysisParam(reqs[i])
		if err != nil {
			res := new(Response)
			res.Error = err.Error()
			ress = append(ress, *res)
			conn.Write(encodeJsonArray(ress))
			return
		}
		results := reflect.ValueOf(this.FunctionMap[reqs[i].MethodName]).Call(values)
		res := new(Response)
		res.MethodName = reqs[i].MethodName
		res.ID = reqs[i].ID
		for _, result := range results {
			param := new(Parameter)
			param.ParamType = result.Kind().String()
			param.ParamValue = result.Interface()
			res.Param = append(res.Param, *param)
		}
		ress = append(ress, *res)
	}
	writeStr := convert2WriteString(encodeJsonArray(ress))
	conn.Write(writeStr)
}

func convert2WriteString(jsonArray []byte) []byte {
	bytes := make([]byte, 0, 512)
	bytes = append(bytes, '$')
	bytes = append(bytes, jsonArray...)
	bytes = append(bytes, '#')
	return bytes
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error occured : %s", err.Error())
		os.Exit(1)
	}
}
