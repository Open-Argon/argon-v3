package main

var commentCompile = makeRegex("( *)//.*")

func isComment(code UNPARSEcode) bool {
	return commentCompile.MatchString(code.code)
}
