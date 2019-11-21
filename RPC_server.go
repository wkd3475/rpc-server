package main

import (
	"fmt"
	"net"
	"net/rpc"
	"os"
)

type Request int

type Args struct {
	Root string
}

type Reply struct {
	Files []string
}

func OSReadDir(root string) ([]string, error) {
	var files []string
	f, err := os.Open(root)
	if err != nil {
		return files, err
	}
	fileInfo, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		files = append(files, file.Name())
	}
	return files, nil
}

func (c *Request) Ls(args *Args, reply *Reply) error {
	var (
		files []string
		err   error
	)
	location := args.Root
	files, err = OSReadDir(location)
	if err != nil {
		panic(err)
	}

	reply.Files = files

	return nil
}

func main() {
	rpc.Register(new(Request))
	ln, err := net.Listen("tcp", ":6000")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()

		if conn != nil {
			fmt.Println("hello")
		}
		if err != nil {
			continue
		}
		defer conn.Close()

		go rpc.ServeConn(conn)
	}
}
