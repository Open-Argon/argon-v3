package main

import (
	"fmt"
	"net"
)

func ArSocket(args ...any) (any, ArErr) {
	if len(args) != 2 {
		return ArObject{}, ArErr{
			TYPE:    "SocketError",
			message: "Socket takes exactly 2 arguments",
			EXISTS:  true,
		}
	} else if typeof(args[0]) != "string" {
		return ArObject{}, ArErr{
			TYPE:    "SocketError",
			message: "Socket type must be a string",
			EXISTS:  true,
		}
	} else if typeof(args[1]) != "number" {
		return ArObject{}, ArErr{
			TYPE:    "SocketError",
			message: "Socket port must be a number",
			EXISTS:  true,
		}
	}
	networktype := ArValidToAny(args[0]).(string)
	port := args[1].(number)
	if port.Denom().Int64() != 1 {
		return ArObject{}, ArErr{
			TYPE:    "SocketError",
			message: "Socket port must be an integer",
			EXISTS:  true,
		}
	}
	ln, err := net.Listen(networktype, ":"+fmt.Sprint(port.Num().Int64()))
	if err != nil {
		return ArObject{}, ArErr{
			TYPE:    "SocketError",
			message: err.Error(),
			EXISTS:  true,
		}
	}
	return Map(anymap{
		"accept": builtinFunc{
			"accept",
			func(args ...any) (any, ArErr) {
				conn, err := ln.Accept()
				if err != nil {
					return ArObject{}, ArErr{
						TYPE:    "SocketError",
						message: err.Error(),
						EXISTS:  true,
					}
				}
				return Map(anymap{
					"read": builtinFunc{
						"read",
						func(args ...any) (any, ArErr) {
							if len(args) != 1 {
								return ArObject{}, ArErr{
									TYPE:    "SocketError",
									message: "Socket.read() takes exactly 1 argument",
									EXISTS:  true,
								}
							} else if typeof(args[0]) != "number" {
								return ArObject{}, ArErr{
									TYPE:    "SocketError",
									message: "Socket.read() argument must be a number",
									EXISTS:  true,
								}
							}
							buf := make([]byte, args[0].(number).Num().Int64())
							n, err := conn.Read(buf)
							if err != nil {
								return ArObject{}, ArErr{
									TYPE:    "SocketError",
									message: err.Error(),
									EXISTS:  true,
								}
							}
							return ArString(string(buf[:n])), ArErr{}
						},
					},
					"writeData": builtinFunc{
						"writeData",
						func(args ...any) (any, ArErr) {
							if len(args) != 1 {
								return ArObject{}, ArErr{
									TYPE:    "SocketError",
									message: "Socket.writeData() takes exactly 1 argument",
									EXISTS:  true,
								}
							}
							data := ArValidToAny(args[0])
							switch x := data.(type) {
							case []any:
								bytes := []byte{}
								for _, v := range x {
									bytes = append(bytes, byte(v.(number).Num().Int64()))
								}
								conn.Write(bytes)
								return nil, ArErr{}
							}
							return nil, ArErr{
								TYPE:    "SocketError",
								message: "Socket.writeData() argument must be a array of numbers",
							}
						},
					},
					"write": builtinFunc{
						"write",
						func(args ...any) (any, ArErr) {
							if len(args) != 1 {
								return ArObject{}, ArErr{
									TYPE:    "SocketError",
									message: "Socket.write() takes exactly 1 argument",
									EXISTS:  true,
								}
							} else if typeof(args[0]) != "string" {
								return ArObject{}, ArErr{
									TYPE:    "SocketError",
									message: "Socket.write() argument must be a string",
									EXISTS:  true,
								}
							}
							data := ArValidToAny(args[0]).(string)
							conn.Write([]byte(data))
							return nil, ArErr{}
						},
					},
					"close": builtinFunc{
						"close",
						func(args ...any) (any, ArErr) {
							conn.Close()
							conn = nil
							return nil, ArErr{}
						},
					},
					"isClosed": builtinFunc{
						"isClosed",
						func(args ...any) (any, ArErr) {
							return conn == nil, ArErr{}
						},
					},
					"RemoteAddr": builtinFunc{
						"RemoteAddr",
						func(args ...any) (any, ArErr) {
							return ArString(conn.RemoteAddr().String()), ArErr{}
						},
					},
					"LocalAddr": builtinFunc{
						"LocalAddr",
						func(args ...any) (any, ArErr) {
							return ArString(conn.LocalAddr().String()), ArErr{}
						},
					},
				}), ArErr{}
			},
		},
		"close": builtinFunc{
			"close",
			func(args ...any) (any, ArErr) {
				ln.Close()
				ln = nil
				return nil, ArErr{}
			},
		},
		"isClosed": builtinFunc{
			"isClosed",
			func(args ...any) (any, ArErr) {
				return ln == nil, ArErr{}
			},
		},
	}), ArErr{}
}
