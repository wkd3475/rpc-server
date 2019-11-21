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
	rpc.Register(new(Request))            // Request 타입의 인스턴스를 생성하여 RPC 서버에 등록
	ln, err := net.Listen("tcp", ":6000") // TCP 프로토콜에 6000번 포트로 연결을 받음
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ln.Close() // main 함수가 종료되기 직전에 연결 대기를 닫음

	for {
		conn, err := ln.Accept() // 클라이언트가 연결되면 TCP 연결을 리턴

		if conn != nil {
			fmt.Println("hello")
		}
		if err != nil {
			continue
		}
		defer conn.Close() // main 함수가 끝나기 직전에 TCP 연결을 닫음

		go rpc.ServeConn(conn) // RPC를 처리하는 함수를 고루틴으로 실행
	}
}
