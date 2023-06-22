package main

import (
	"bufio"
	"errors"
	"log"
	"os"
	"path/filepath"
)

var imported = make(map[string]ArObject)
var importing = make(map[string]bool)

var modules_folder = "argon_modules"

func FileExists(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		return true

	} else if errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		return false
	}
}

func readFile(path string) []UNPARSEcode {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	output := []UNPARSEcode{}
	line := 1
	for scanner.Scan() {
		text := scanner.Text()
		output = append(output, UNPARSEcode{text, text, line, path})
		line++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return nil
	}
	return output
}

func importMod(realpath string, origin string, main bool, global ArObject) (ArObject, ArErr) {
	extention := filepath.Ext(realpath)
	path := realpath
	if extention == "" {
		path += ".ar"
	}
	ex, err := os.Getwd()
	if err != nil {
		return ArObject{}, ArErr{TYPE: "Import Error", message: "Could not get working directory", EXISTS: true}
	}
	exc, err := os.Executable()
	if err != nil {
		return ArObject{}, ArErr{TYPE: "Import Error", message: "Could not get executable", EXISTS: true}
	}
	executable := filepath.Dir(exc)
	isABS := filepath.IsAbs(path)
	var pathsToTest []string
	if isABS {
		pathsToTest = []string{
			filepath.Join(path),
			filepath.Join(realpath, "init.ar"),
		}
	} else {
		pathsToTest = []string{
			filepath.Join(origin, path),
			filepath.Join(origin, realpath, "init.ar"),
			filepath.Join(origin, modules_folder, path),
			filepath.Join(origin, modules_folder, realpath, "init.ar"),
			filepath.Join(ex, path),
			filepath.Join(ex, modules_folder, path),
			filepath.Join(ex, modules_folder, realpath, "init.ar"),
			filepath.Join(executable, modules_folder, path),
			filepath.Join(executable, modules_folder, realpath, "init.ar"),
		}
	}
	var p string
	var found bool
	for _, p = range pathsToTest {
		if FileExists(p) {
			found = true
			break
		}
	}

	if !found {
		return ArObject{}, ArErr{TYPE: "Import Error", message: "File does not exist: " + path, EXISTS: true}
	} else if importing[p] {
		return ArObject{}, ArErr{TYPE: "Import Error", message: "Circular import: " + path, EXISTS: true}
	} else if _, ok := imported[p]; ok {
		return imported[p], ArErr{}
	}
	importing[p] = true
	codelines := readFile(p)

	translated, translationerr := translate(codelines)
	if translationerr.EXISTS {
		return ArObject{}, translationerr
	}
	ArgsArArray := []any{}
	withoutarfile := []string{}
	if len(Args) > 1 {
		withoutarfile = Args[1:]
	}
	for _, arg := range withoutarfile {
		ArgsArArray = append(ArgsArArray, arg)
	}
	local := newscope()
	localvars := Map(anymap{
		"program": Map(anymap{
			"args":   ArArray(ArgsArArray),
			"origin": origin,
			"import": builtinFunc{"import", func(args ...any) (any, ArErr) {
				if len(args) != 1 {
					return nil, ArErr{"Import Error", "Invalid number of arguments", 0, realpath, "", true}
				}
				if _, ok := args[0].(string); !ok {
					return nil, ArErr{"Import Error", "Invalid argument type", 0, realpath, "", true}
				}
				return importMod(args[0].(string), filepath.Dir(p), false, global)
			}},
			"cwd": ex,
			"exc": exc,
			"file": Map(anymap{
				"name": filepath.Base(p),
				"path": p,
			}),
			"main": main,
		}),
	})
	_, runimeErr := ThrowOnNonLoop(run(translated, stack{global, localvars, local}))
	importing[p] = false
	if runimeErr.EXISTS {
		return ArObject{}, runimeErr
	}
	imported[p] = local
	return local, ArErr{}
}
