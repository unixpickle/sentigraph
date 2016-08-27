package main

import (
	"bytes"
	"compress/gzip"
	"io"

	"github.com/gopherjs/gopherjs/js"
)

func main() {
	js.Global.Set("onmessage", js.MakeFunc(messageHandler))
}

func messageHandler(this *js.Object, dataArg []*js.Object) interface{} {
	if len(dataArg) != 1 {
		panic("expected one argument")
	}
	data := dataArg[0].Get("data").Interface().([]byte)
	newData, err := decompressData(data)
	if err != nil {
		js.Global.Call("postMessage", []interface{}{nil, err.Error()})
	} else {
		js.Global.Call("postMessage", []interface{}{js.NewArrayBuffer(newData), nil})
	}
	js.Global.Call("close")
	return nil
}

func decompressData(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	var result bytes.Buffer
	if _, err := io.Copy(&result, reader); err != nil {
		return nil, err
	}
	return result.Bytes(), nil
}
