package main

import (
	wapc "github.com/wapc/wapc-guest-tinygo"
)

func main() {
	wapc.RegisterFunctions(wapc.Functions{
		"hello": hello,
	})
}

func hello(payload []byte) ([]byte, error) {
	wapc.ConsoleLog("Logging something")
	wapc.HostCall("myBinding", "sample", "hello", []byte("Simon"))
	return []byte("Hello"), nil
}
