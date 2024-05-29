package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

var imported = make(map[string]ArObject)
var translatedImports = make(map[string]translatedImport)
var importing = make(map[string]bool)

const modules_folder = "argon_modules"

func FileExists(filename string) bool {
	if info, err := os.Stat(filename); err == nil && !info.IsDir() {
		return true
	}
	return false
}

func readFile(path string) ([]UNPARSEcode, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	return output, nil
}

type translatedImport struct {
	translated []any
	p          string
	path       string
	origin     string
}

var runTranslatedImport func(translatedImport, ArObject, bool) (ArObject, ArErr)
var ex string
var exc string
var exc_dir string

func init() {
	runTranslatedImport = __runTranslatedImport
	ex, _ = os.Getwd()
	exc, _ = os.Executable()
	exc_dir = filepath.Dir(exc)
}

func __runTranslatedImport(translatedImport translatedImport, global ArObject, main bool) (ArObject, ArErr) {

	if _, ok := imported[translatedImport.p]; ok {
		return imported[translatedImport.p], ArErr{}
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
			"origin": ArString(translatedImport.origin),
			"cwd":    ArString(ex),
			"exc":    ArString(exc),
			"file": Map(anymap{
				"name": ArString(filepath.Base(translatedImport.p)),
				"path": ArString(translatedImport.p),
			}),
			"main": main,
		}),
	})
	imported[translatedImport.p] = local
	_, runimeErr := run(translatedImport.translated, stack{global, localvars, local})
	if runimeErr.EXISTS {
		return ArObject{}, runimeErr
	}
	return local, ArErr{}
}

func translateImport(realpath string, origin string, topLevelOnly bool) (translatedImport, ArErr) {
	extention := filepath.Ext(realpath)
	path := realpath
	if extention == "" {
		path += ".ar"
	}
	isABS := filepath.IsAbs(path)
	var pathsToTest []string
	if isABS {
		pathsToTest = []string{
			filepath.Join(path),
			filepath.Join(realpath, "init.ar"),
		}
	} else {
		pathsToTest = []string{
			filepath.Join(exc_dir, path),
			filepath.Join(exc_dir, realpath, "init.ar"),
			filepath.Join(exc_dir, modules_folder, path),
			filepath.Join(exc_dir, modules_folder, realpath, "init.ar"),
		}
		var currentPath string = origin
		var oldPath string = ""
		for currentPath != oldPath {
			pathsToTest = append(pathsToTest,
				filepath.Join(currentPath, path),
				filepath.Join(currentPath, realpath, "init.ar"),
				filepath.Join(currentPath, modules_folder, path),
				filepath.Join(currentPath, modules_folder, realpath, "init.ar"))
			if topLevelOnly {
				break
			}
			oldPath = currentPath
			currentPath = filepath.Dir(currentPath)
		}
		fmt.Println(pathsToTest)
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
		return translatedImport{}, ArErr{TYPE: "Import Error", message: "File does not exist: " + path, EXISTS: true}
	} else if importing[p] {
		return translatedImport{}, ArErr{TYPE: "Import Error", message: "Circular import: " + path, EXISTS: true}
	} else if _, ok := translatedImports[p]; ok {
		return translatedImports[p], ArErr{}
	}
	importing[p] = true
	codelines, err := readFile(p)
	if err != nil {
		return translatedImport{}, ArErr{TYPE: "Import Error", message: "Could not read file: " + path, EXISTS: true}
	}

	importing[p] = true
	translated, translationerr := translate(codelines)
	importing[p] = false

	if translationerr.EXISTS {
		return translatedImport{}, translationerr
	}

	translatedImports[p] = translatedImport{translated, p, path, origin}
	return translatedImports[p], ArErr{}
}
