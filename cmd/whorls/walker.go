//
// walker.go
// Copyright (C) 2017 Kevin Lyda <kevin@phrye.com>
//
// Distributed under terms of the GPL license.
//

package main

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/template"
)

// Whorls captures the mapping of file name to fingerprinted file name
type Whorls struct {
	Name, Whorl string
}

// Walker is the data structure for walking the source trees.
type Walker struct {
	Out, Package, NotFound, GoSrc string
	Mappings                      chan Whorls
	Done                          chan bool
}

func makeHash(path string) string {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("FATAL: Can't open '%s' (%s)", path, err)
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatalf("FATAL: Can't hash '%s' (%s)", path, err)
	}

	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))[:6]
}

func copyFile(dst, src string) {
	d, errDst := os.Create(dst)
	if errDst != nil {
		log.Fatalf("FATAL: Can't open '%s' (%s)", dst, errDst)
	}
	defer d.Close()
	s, errSrc := os.Open(src)
	if errSrc != nil {
		log.Fatalf("FATAL: Can't open '%s' (%s)", src, errSrc)
	}
	defer s.Close()
	if _, err := io.Copy(d, s); err != nil {
		log.Fatalf("FATAL: Couldn't copy '%s' to '%s' (%s)", src, dst, err)
	}
}

// GenerateWhorl adds files to the asset directory.
func (w Walker) GenerateWhorl() func(string, os.FileInfo, error) error {
	return func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			var whorl string
			digest := makeHash(path)
			slashless := strings.Replace(path, "/", "-", -1)
			dotsplit := strings.Split(slashless, ".")
			if len(dotsplit) > 1 {
				whorl = fmt.Sprintf("%s/%s-%s.%s",
					w.Out, strings.Join(dotsplit[:len(dotsplit)-1], "."),
					digest, dotsplit[len(dotsplit)-1])
			} else {
				whorl = fmt.Sprintf("%s/%s-%s",
					w.Out, digest, slashless)
			}
			copyFile(whorl, path)
			w.Mappings <- Whorls{Name: path, Whorl: whorl}
		}
		return nil
	}
}

// GenerateSource creates the source file.
func (w Walker) GenerateSource() {
	wv := WhorlItValues{
		Command:  strings.Join(os.Args, " "),
		Package:  w.Package,
		NotFound: w.NotFound}
	if gosrcF, err := os.Create(w.GoSrc); err != nil {
		log.Fatalf("FATAL: Can't open file: %s", err.Error())
	} else {
		template.Must(template.New("prefix").
			Parse(sourcePrefix)).
			Execute(gosrcF, wv)
		for {
			mapping, more := <-w.Mappings
			if more {
				fmt.Fprintf(gosrcF,
					"\t\"%s\": \"%s\",\n", mapping.Name, mapping.Whorl)
			} else {
				template.Must(template.New("suffix").
					Parse(sourceSuffix)).
					Execute(gosrcF, wv)
				gosrcF.Close()
				w.Done <- true
				return
			}
		}
	}
}
