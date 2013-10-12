golang-rpc-json
===============

A json-rpc server implement by Google GO.

Document
===============
If you registe a swap function to swap two integers.

The request should be a jsonarray like this :
[{"ID":123 , "MethodName":"swap" , "Param":[{"ParamType":"int","ParamValue":9} , {"ParamType":"int","ParamValue":8}]}]

And the response will be a jsonarray like this :
[{"ID":123 , "MethodName":"swap" , "Error" : null , "Param":[{"ParamType":"int","ParamValue":8} , {"ParamType":"int","ParamValue":9}]}]

"Error" in the response will not be null if errs occured in this remote function call.

Example
===============
```go
package main

import "jsonrpc"

func main(){
	rpcserver := new(jsonrpc.JsonRpcServer)
	
	//set the port you want to listen with a int-type-parameter
	
	rpcserver.SetPort(7777)
	
	//set the protocol you want to use, but now it just supports tcp4
	//you can use "tcp" or "tcp4".
	
	rpcserver.SetProtocol("tcp")
	
	//registe your function into this server .The type of function's
	//parameter , both input and return , should be in this set {int,int8,
	//in16,int32,int64,uint,uint8,uint16,uint32,uint64,string,bool}
	
	rpcserver.RegisteFunc("whatever name you want" , functionName)
	
	//start this server
	
	rpcserver.Serve()
}
```

Licence
===============
MIT

