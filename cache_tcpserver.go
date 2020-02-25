package gocache

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
)

type TcpBnfHandler func(req *Request, resp *Response)
type TcpBnfServer struct {
	cache  Cache
	l      net.Listener
	router map[byte]TcpBnfHandler
}

type Request struct {
	op    byte
	key   string
	value []byte
}

type Response struct {
	value []byte
	err   error
}

func NewTcpBnfServer() *TcpBnfServer {
	server := &TcpBnfServer{cache: newMemoryCache(), router: make(map[byte]TcpBnfHandler)}
	server.router['S'] = server.handleSetCache
	return server
}

func (server *TcpBnfServer) Run(addr string) {
	li, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println(err)
		return
	}
	server.l = li
	for {
		conn, err := server.l.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go server.handleRequest(conn)
	}
}

func (server *TcpBnfServer) handleRequest(conn net.Conn) {
	defer conn.Close()

	req, err := server.decodeRequest(conn)
	if err != nil {
		log.Println(err)
		return
	}
	var resp Response
	server.router[req.op](req, &resp)
	server.sendResponse(conn, &resp)
}

func (server *TcpBnfServer) decodeRequest(r io.Reader) (*Request, error) {
	br := bufio.NewReader(r)
	op, err := br.ReadByte()
	if err != nil {
		return nil, err
	}
	if _, ok := server.router[op]; !ok {
		return nil, errors.New("bad operation")
	}

	klenString, err := br.ReadString(' ')
	if err != nil {
		return nil, err
	}
	klen, err := strconv.Atoi(strings.TrimSpace(klenString))
	if err != nil {
		return nil, err
	}

	vlenString, err := br.ReadString(' ')
	if err != nil {
		return nil, err
	}
	vlen, err := strconv.Atoi(strings.TrimSpace(vlenString))
	if err != nil {
		return nil, err
	}

	key := make([]byte, klen)
	_, err = io.ReadFull(br, key)
	if err != nil {
		return nil, err
	}

	value := make([]byte, vlen)
	_, err = io.ReadFull(br, value)
	if err != nil {
		return nil, err
	}

	return &Request{op: op, key: string(key), value: value}, nil
}

func (server *TcpBnfServer) handleSetCache(req *Request, resp *Response) {
	err := server.cache.Set(req.key, req.value)
	resp.err = err
}

func (server *TcpBnfServer) sendResponse(conn net.Conn, resp *Response) {
	var buf string
	if resp.err != nil {
		buf = fmt.Sprintf("-%d %s", len(resp.err.Error()), resp.err.Error())
	} else {
		buf = fmt.Sprintf("%d %s", len(resp.value), resp.value)
	}
	conn.Write([]byte(buf))
}
