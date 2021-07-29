package main

import (
	"encoding/base64"
	"fmt"
	"github.com/axgle/mahonia"
	"strings"
)

// 异或
func xor(content string) string{
	test := make([]byte, len(content))
	for i:=0; i < len(content); i++ {
		//fmt.Println(content[i] ^ 0x77)
		test[i] = content[i] ^ 0x77
	}
	//fmt.Println(string(test))
	return string(test)
}

//base64 encode
func my_base_encode(s []byte)string{
	encodeStd := "SBk0EFGHIJKLMNrPQRATUXWVYZabcdefghijClmnvpqOstuowxyzD128456739#%"
	s64 := base64.NewEncoding(encodeStd).EncodeToString(s)
	s64 = strings.ReplaceAll(s64, "=", "-")
	//fmt.Println("self编码",s64)
	//fmt.Println("self.length", len(s64))
	return s64
}


//base64 decode
func my_base_decode(s string)string{
	encodeStd := "SBk0EFGHIJKLMNrPQRATUXWVYZabcdefghijClmnvpqOstuowxyzD128456739#%"
	s = strings.ReplaceAll(s,"-","=")
	decodeBytes, err := base64.NewEncoding(encodeStd).DecodeString(s)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println("self解码", decodeBytes)
	return string(decodeBytes)
}

// 数据加密
func encrypt(s string) string {
	// base encode
	//data := my_base_encode([] byte(s))
	//trans cstr to gostr
	//go_s := C.GoString(s)
	// xor
	data := xor(s)
	// base encode
	data = my_base_encode([] byte(data))
	// trans gostr to cstr
	//c_str := C.CString(data)
	//defer C.free(unsafe.Pointer(c_str))
	return data
}


// 数据解密
func decrypt(s string) string {
	//trans cstr to gostr
	//go_s := C.GoString(s)
	// base decode
	data := my_base_decode(s)
	// xor
	data = xor(data)
	// 编码
	dec := mahonia.NewDecoder("gbk")
	data = dec.ConvertString(data)
	//fmt.Println(data)
	// trans gostr to cstr
	//c_str := C.CString(data)
	//defer C.free(unsafe.Pointer(c_str))
	return data
}
