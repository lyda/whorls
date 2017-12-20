//
// template.go
// Copyright (C) 2017 Kevin Lyda <kevin@phrye.com>
//
// Distributed under terms of the GPL license.
//

package main

// WhorlItValues are the replacements into the generated source template.
type WhorlItValues struct {
	Command, Package, NotFound string
}

var (
	sourcePrefix = `// Code generated by "{{.Command}}"; DO NOT EDIT.

package {{.Package}}

// WhorlIt will convert an asset path to a fingerprinted path.
func WhorlIt(path string) string {
	if value, ok := generatedWhorls[path]; ok {
		return value
	}
	return "{{.NotFound}}"
}

var generatedWhorls = map[string]string{
`
	sourceSuffix = "}\n"
)
