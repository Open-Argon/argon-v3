package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

func ArSocketClient(args ...any) (any, ArErr) {
	if len(args) != 2 {
		return ArObject{}, ArErr{
			TYPE:    "Socket Error",
			message: "Socket takes exactly 2 arguments",
			EXISTS:  true,
		}
	} else if typeof(args[0]) != "string" {
		return ArObject{}, ArErr{
			TYPE:    "Socket Error",
			message: "Socket type must be a string",
			EXISTS:  true,
		}
	} else if typeof(args[1]) != "string" {
		return ArObject{}, ArErr{
			TYPE:    "Socket Error",
			message: "Socket address must be a string",
			EXISTS:  true,
		}
	}
	networktype := ArValidToAny(args[0]).(string)
	address := ArValidToAny(args[1]).(string)
	conn, err := net.Dial(networktype, address)
	if err != nil {
		return ArObject{}, ArErr{
			TYPE:    "Socket Error",
			message: fmt.Sprintf("Socket connection failed: %s", err.Error()),
			EXISTS:  true,
		}
	}
	return ArObject{
		obj: anymap{
			"read": builtinFunc{
				"read",
				func(args ...any) (any, ArErr) {
					if len(args) != 1 {
						return ArObject{}, ArErr{
							TYPE:    "Socket Error",
							message: "Socket.readData() takes exactly 1 argument",
							EXISTS:  true,
						}
					}
					if conn == nil {
						return ArObject{}, ArErr{
							TYPE:    "Socket Error",
							message: "Connection is closed",
							EXISTS:  true,
						}
					}
					buf := make([]byte, args[0].(number).Num().Int64())
					n, err := conn.Read(buf)
					if err != nil {
						return ArObject{}, ArErr{
							TYPE:    "Socket Error",
							message: fmt.Sprintf("Socket read failed: %s", err.Error()),
							EXISTS:  true,
						}
					}
					return ArBuffer(buf[:n]), ArErr{}
				}},
			"readUntil": builtinFunc{
				"readUntil",
				func(args ...any) (any, ArErr) {
					if len(args) != 1 {
						return ArObject{}, ArErr{
							TYPE:    "Socket Error",
							message: "Socket.readUntil() takes exactly 1 argument",
							EXISTS:  true,
						}
					}
					if conn == nil {
						return ArObject{}, ArErr{
							TYPE:    "Socket Error",
							message: "Connection is closed",
							EXISTS:  true,
						}
					}
					value := ArValidToAny(args[0])
					if typeof(value) != "buffer" {
						return ArObject{}, ArErr{
							TYPE:    "Type Error",
							message: fmt.Sprintf("Socket.readUntil() argument must be a buffer, not %s", typeof(value)),
							EXISTS:  true,
						}
					}
					endBuf := value.([]byte)
					reader := io.Reader(conn)
					buf := make([]byte, len(endBuf))
					var data []byte
					for {
						_, err := reader.Read(buf)
						if err != nil {
							return ArObject{}, ArErr{
								TYPE:    "Socket Error",
								message: fmt.Sprintf("Socket read failed: %s", err.Error()),
								EXISTS:  true,
							}
						}
						data = append(data, buf[0])
						if len(data) >= len(endBuf) {
							dataSlice := data[len(data)-len(endBuf):]
							for i := 0; i < len(endBuf); i++ {
								if dataSlice[i] != endBuf[i] {
									break
								}
								if i == len(endBuf)-1 {
									return ArBuffer(data), ArErr{}
								}
							}
						}
					}

				}},
			"write": builtinFunc{
				"write",
				func(args ...any) (any, ArErr) {
					if len(args) != 1 {
						return nil, ArErr{
							TYPE:    "Type Error",
							message: fmt.Sprintf("write() takes exactly 1 argument (%d given)", len(args)),
							EXISTS:  true,
						}
					}
					if typeof(args[0]) != "buffer" {
						return ArObject{}, ArErr{
							TYPE:    "Type Error",
							message: fmt.Sprintf("write() argument must be a buffer, not %s", typeof(args[0])),
							EXISTS:  true,
						}
					}
					if conn == nil {
						return ArObject{}, ArErr{
							TYPE:    "Socket Error",
							message: "Connection is closed",
							EXISTS:  true,
						}
					}
					args[0] = ArValidToAny(args[0])
					if typeof(args[0]) != "buffer" {
						return ArObject{}, ArErr{
							TYPE:    "Type Error",
							message: fmt.Sprintf("write() argument must be a buffer, not %s", typeof(args[0])),
							EXISTS:  true,
						}
					}
					_, err := conn.Write(args[0].([]byte))
					if err != nil {
						return ArObject{}, ArErr{
							TYPE:    "Socket Error",
							message: err.Error(),
							EXISTS:  true,
						}
					}
					return nil, ArErr{}
				}},
			"close": builtinFunc{
				"close",
				func(args ...any) (any, ArErr) {
					if conn == nil {
						return ArObject{}, ArErr{
							TYPE:    "Socket Error",
							message: "Connection is already closed",
							EXISTS:  true,
						}
					}
					err := conn.Close()
					if err != nil {
						return ArObject{}, ArErr{
							TYPE:    "Socket Error",
							message: err.Error(),
							EXISTS:  true,
						}
					}
					conn = nil
					return nil, ArErr{}
				},
			},
			"isClosed": builtinFunc{
				"isClosed",
				func(args ...any) (any, ArErr) {
					if conn == nil {
						return true, ArErr{}
					}
					conn.SetWriteDeadline(time.Now().Add(1 * time.Millisecond))
					_, err := conn.Write([]byte{})
					conn.SetWriteDeadline(time.Time{})
					if err != nil {
						conn.Close()
						conn = nil
						return true, ArErr{}
					}
					return false, ArErr{}

				},
			},
		}}, ArErr{}
}

func ArSocketServer(args ...any) (any, ArErr) {
	if len(args) != 2 {
		return ArObject{}, ArErr{
			TYPE:    "Socket Error",
			message: "Socket takes exactly 2 arguments",
			EXISTS:  true,
		}
	} else if typeof(args[0]) != "string" {
		return ArObject{}, ArErr{
			TYPE:    "Socket Error",
			message: "Socket type must be a string",
			EXISTS:  true,
		}
	} else if typeof(args[1]) != "number" {
		return ArObject{}, ArErr{
			TYPE:    "Socket Error",
			message: "Socket port must be a number",
			EXISTS:  true,
		}
	}
	networktype := ArValidToAny(args[0]).(string)
	port_num := args[1].(ArObject)
	port, err := numberToInt64(port_num)
	if err != nil {
		return ArObject{}, ArErr{
			TYPE:    "Socket Error",
			message: err.Error(),
			EXISTS:  true,
		}
	}
	ln, err := net.Listen(networktype, ":"+fmt.Sprint(port))
	if err != nil {
		return ArObject{}, ArErr{
			TYPE:    "Socket Error",
			message: err.Error(),
			EXISTS:  true,
		}
	}
	return Map(anymap{
		"accept": builtinFunc{
			"accept",
			func(args ...any) (any, ArErr) {
				if ln == nil {
					return ArObject{}, ArErr{
						TYPE:    "Socket Error",
						message: "Socket is closed",
						EXISTS:  true,
					}
				}
				conn, err := ln.Accept()
				if err != nil {
					return ArObject{}, ArErr{
						TYPE:    "Socket Error",
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
									TYPE:    "Socket Error",
									message: "Socket.readData() takes exactly 1 argument",
									EXISTS:  true,
								}
							}
							if conn == nil {
								return ArObject{}, ArErr{
									TYPE:    "Socket Error",
									message: "Connection is closed",
									EXISTS:  true,
								}
							}
							buf := make([]byte, args[0].(number).Num().Int64())
							n, err := conn.Read(buf)
							if err != nil {
								return ArObject{}, ArErr{
									TYPE:    "Socket Error",
									message: err.Error(),
									EXISTS:  true,
								}
							}
							return ArBuffer(buf[:n]), ArErr{}
						},
					},
					"readUntil": builtinFunc{
						"readUntil",
						func(args ...any) (any, ArErr) {
							if len(args) != 1 {
								return ArObject{}, ArErr{
									TYPE:    "Socket Error",
									message: "Socket.readUntil() takes exactly 1 argument",
									EXISTS:  true,
								}
							}
							if conn == nil {
								return ArObject{}, ArErr{
									TYPE:    "Socket Error",
									message: "Connection is closed",
									EXISTS:  true,
								}
							}
							value := ArValidToAny(args[0])
							if typeof(value) != "buffer" {
								return ArObject{}, ArErr{
									TYPE:    "Type Error",
									message: fmt.Sprintf("Socket.readUntil() argument must be a buffer, not %s", typeof(value)),
									EXISTS:  true,
								}
							}
							endBuf := value.([]byte)
							var data []byte
							buf := make([]byte, 1)
							lookingAt := 0
							for {
								n, err := io.ReadFull(conn, buf)
								if err != nil {
									return ArBuffer(data), ArErr{}
								}
								chunk := buf[:n]
								data = append(data, chunk...)
								for i := 0; i < n; i++ {
									if chunk[i] == endBuf[lookingAt] {
										lookingAt++
										if lookingAt == len(endBuf) {
											data = append(data, chunk...)
											return ArBuffer(data), ArErr{}
										}
									} else {
										lookingAt = 0
									}
								}

							}
						}},
					"clearTimeout": builtinFunc{
						"clearTimeout",
						func(args ...any) (any, ArErr) {
							if len(args) != 0 {
								return ArObject{}, ArErr{
									TYPE:    "Socket Error",
									message: "Socket.clearTimeout() takes exactly 0 arguments",
									EXISTS:  true,
								}
							}
							if conn == nil {
								return ArObject{}, ArErr{
									TYPE:    "Socket Error",
									message: "Connection is closed",
									EXISTS:  true,
								}
							}
							conn.SetDeadline(time.Time{})
							return ArObject{}, ArErr{}
						},
					},
					"setTimeout": builtinFunc{
						"setTimeout",
						func(args ...any) (any, ArErr) {
							if len(args) != 1 {
								return ArObject{}, ArErr{
									TYPE:    "Socket Error",
									message: "Socket.setTimeout() takes exactly 1 argument",
									EXISTS:  true,
								}
							}
							if typeof(args[0]) != "number" {
								return ArObject{}, ArErr{
									TYPE:    "Socket Error",
									message: "Socket timeout must be a number",
									EXISTS:  true,
								}
							}
							if conn == nil {
								return ArObject{}, ArErr{
									TYPE:    "Socket Error",
									message: "Connection is closed",
									EXISTS:  true,
								}
							}
							timeout := args[0].(number)
							if timeout.Denom().Int64() != 1 {
								return ArObject{}, ArErr{
									TYPE:    "Socket Error",
									message: "Socket timeout must be an integer",
									EXISTS:  true,
								}
							}
							err := conn.SetDeadline(time.Now().Add(time.Duration(timeout.Num().Int64()) * time.Millisecond))
							if err != nil {
								return ArObject{}, ArErr{
									TYPE:    "Socket Error",
									message: err.Error(),
									EXISTS:  true,
								}
							}
							return ArObject{}, ArErr{}
						},
					},
					"write": builtinFunc{
						"write",
						func(args ...any) (any, ArErr) {
							if len(args) != 1 {
								return ArObject{}, ArErr{
									TYPE:    "Socket Error",
									message: "Socket.writeData() takes exactly 1 argument",
									EXISTS:  true,
								}
							}
							if conn == nil {
								return ArObject{}, ArErr{
									TYPE:    "Socket Error",
									message: "Connection is closed",
									EXISTS:  true,
								}
							}
							data := ArValidToAny(args[0])
							switch x := data.(type) {
							case []any:
								bytes := []byte{}
								for _, v := range x {
									if typeof(v) != "number" && v.(number).Denom().Int64() != 1 {
										return ArObject{}, ArErr{
											TYPE:    "Socket Error",
											message: "Socket.writeData() argument must be a array of integers",
											EXISTS:  true,
										}
									}
									bytes = append(bytes, byte(v.(number).Num().Int64()))
								}
								conn.Write(bytes)
								return nil, ArErr{}
							case []byte:
								conn.Write(x)
								return nil, ArErr{}
							}
							return nil, ArErr{
								TYPE:    "Socket Error",
								message: "Socket.writeData() argument must be a array of numbers",
							}
						},
					},
					"close": builtinFunc{
						"close",
						func(args ...any) (any, ArErr) {
							if conn == nil {
								return ArObject{}, ArErr{
									TYPE:    "Socket Error",
									message: "Connection is already closed",
									EXISTS:  true,
								}
							}
							conn.Close()
							conn = nil
							return nil, ArErr{}
						},
					},
					"isClosed": builtinFunc{
						"isClosed",
						func(args ...any) (any, ArErr) {
							if conn == nil {
								return true, ArErr{}
							}
							conn.SetWriteDeadline(time.Now().Add(1 * time.Millisecond))
							_, err := conn.Write([]byte{})
							conn.SetWriteDeadline(time.Time{})
							if err != nil {
								conn.Close()
								conn = nil
								return true, ArErr{}
							}
							return false, ArErr{}

						},
					},
					"RemoteAddr": builtinFunc{
						"RemoteAddr",
						func(args ...any) (any, ArErr) {
							if conn == nil {
								return ArObject{}, ArErr{
									TYPE:    "Socket Error",
									message: "Connection is closed",
									EXISTS:  true,
								}
							}
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
				if ln == nil {
					return ArObject{}, ArErr{
						TYPE:    "Socket Error",
						message: "Socket is already closed",
						EXISTS:  true,
					}
				}
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
