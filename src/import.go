package main

import (
	"bufio"
	"errors"
	"log"
	"os"
	"path/filepath"
)

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

func importMod(realpath string, origin string, main bool) ArErr {
	extention := filepath.Ext(realpath)
	path := realpath
	if extention == "" {
		path += ".ar"
	}
	ex, err := os.Getwd()
	if err != nil {
		return ArErr{"Import Error", err.Error(), 0, realpath, "", true}
	}
	executable, err := os.Executable()
	if err != nil {
		return ArErr{"Import Error", err.Error(), 0, realpath, "", true}
	}
	executable = filepath.Dir(executable)
	isABS := filepath.IsAbs(path)
	var pathsToTest []string
	if isABS {
		pathsToTest = []string{
			filepath.Join(path),
			filepath.Join(realpath, "init.ar"),
		}
	} else {
		pathsToTest = []string{
			filepath.Join(origin, realpath, "init.ar"),
			filepath.Join(origin, path),
			filepath.Join(origin, "modules", path),
			filepath.Join(origin, "modules", realpath, "init.ar"),
			filepath.Join(ex, path),
			filepath.Join(ex, "modules", realpath, "init.ar"),
			filepath.Join(ex, "modules", path),
			filepath.Join(executable, "modules", realpath, "init.ar"),
			filepath.Join(executable, "modules", path),
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
		return ArErr{"Import Error", "File does not exist: " + realpath, 0, realpath, "", true}
	}
	codelines := readFile(p)

	translated, translationerr := translate(codelines)
	if translationerr.EXISTS {
		return translationerr
	}
	_, runimeErr := run(translated, stack{vars})
	if runimeErr.EXISTS {
		return runimeErr
	}
	return ArErr{}
}
