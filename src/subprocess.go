package main

import (
	"bufio"
	"fmt"
	"os/exec"
)

func ArSubprocess(args ...any) (any, ArErr) {
	if len(args) != 1 {
		return nil, ArErr{
			TYPE:    "RuntimeError",
			message: fmt.Sprintf("subprocess() takes exactly 1 argument (%d given)", len(args)),
			EXISTS:  true,
		}
	} else if typeof(args[0]) != "array" {
		fmt.Println(args[0])
		return nil, ArErr{
			TYPE:    "TypeError",
			message: fmt.Sprintf("subprocess() argument must be an array, not %s", typeof(args[0])),
			EXISTS:  true,
		}
	} else if len(args[0].([]any)) == 0 {
		return nil, ArErr{
			TYPE:    "RuntimeError",
			message: "subprocess() argument must be an array of strings, not an empty array",
			EXISTS:  true,
		}
	}
	commandlist := []string{}
	for _, x := range args[0].([]any) {
		if typeof(x) != "string" {
			return nil, ArErr{
				TYPE:    "TypeError",
				message: fmt.Sprintf("subprocess() argument must be an array of strings, not %s", typeof(x)),
				EXISTS:  true,
			}
		}
		commandlist = append(commandlist, x.(ArObject).obj["__value__"].(string))
	}
	cmd := exec.Command(commandlist[0], commandlist[1:]...)
	return Map(
		anymap{
			"stdout": builtinFunc{
				"output",
				func(args ...any) (any, ArErr) {
					stdout, err := cmd.StdoutPipe()

					if err != nil {
						return nil, ArErr{
							TYPE:    "RuntimeError",
							message: err.Error(),
							EXISTS:  true,
						}
					}

					if err := cmd.Start(); err != nil {
						return nil, ArErr{
							TYPE:    "RuntimeError",
							message: err.Error(),
							EXISTS:  true,
						}
					}
					scanner := bufio.NewScanner(stdout) // create a new scanner to read the output
					scanner.Split(bufio.ScanRunes)
					return Map(anymap{
						"scan": builtinFunc{
							"scan",
							func(args ...any) (any, ArErr) {
								if len(args) != 0 {
									return nil, ArErr{
										TYPE:    "TypeError",
										message: fmt.Sprintf("scan() takes exactly 0 arguments (%d given)", len(args)),
										EXISTS:  true,
									}
								}
								return scanner.Scan(), ArErr{}
							}},
						"text": builtinFunc{
							"text",
							func(args ...any) (any, ArErr) {
								if len(args) != 0 {
									return nil, ArErr{
										TYPE:    "TypeError",
										message: fmt.Sprintf("text() takes exactly 0 arguments (%d given)", len(args)),
										EXISTS:  true,
									}
								}
								return scanner.Text(), ArErr{}
							}},
						"err": builtinFunc{
							"err",
							func(args ...any) (any, ArErr) {
								if len(args) != 0 {
									return nil, ArErr{
										TYPE:    "TypeError",
										message: fmt.Sprintf("err() takes exactly 0 arguments (%d given)", len(args)),
										EXISTS:  true,
									}
								}
								if err := scanner.Err(); err != nil {
									return err.Error(), ArErr{}
								}
								return nil, ArErr{}
							}},
						"wait": builtinFunc{
							"wait",
							func(args ...any) (any, ArErr) {
								if len(args) != 0 {
									return nil, ArErr{
										TYPE:    "TypeError",
										message: fmt.Sprintf("wait() takes exactly 0 arguments (%d given)", len(args)),
										EXISTS:  true,
									}
								}
								if err := cmd.Wait(); err != nil {
									return err.Error(), ArErr{}
								}
								return nil, ArErr{}
							},
						},
					}), ArErr{}
				}},
			"output": builtinFunc{
				"output",
				func(args ...any) (any, ArErr) {
					out, err := cmd.Output()
					if err != nil {
						return nil, ArErr{
							TYPE:    "RuntimeError",
							message: err.Error(),
							EXISTS:  true,
						}
					}
					return string(out), ArErr{}
				}},
		},
	), ArErr{}
}
