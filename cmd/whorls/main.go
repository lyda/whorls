//
// main.go
// Copyright (C) 2017 Kevin Lyda <kevin@phrye.com>
//
// Distributed under terms of the GPL license.
//

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	out      = flag.String("o", "assets", "directory where assets go")
	pkg      = flag.String("p", "assets", "name of the package")
	gosrc    = flag.String("f", "generated-assets.go", "generated file name")
	notfound = flag.String("not-found", "", "name if requested whorl missing")
	clear    = flag.Bool("clear", false,
		"delete contents of <out> dir first (default false)")
)

func prepOutDir() {
	if *clear {
		os.RemoveAll(*out)
	}
	err := os.MkdirAll(*out, 0755)
	if os.IsExist(err) && *clear {
		log.Fatalf("FATAL: unable to clear '%s'", *out)
	}
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr,
			"Usage of %s [options] src-dir1 [src-dir2 ... src-dirN]:\n\nOptions:\n",
			os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	prepOutDir()

	w := Walker{
		Out:      *out,
		Package:  *pkg,
		NotFound: *notfound,
		GoSrc:    *gosrc,
		Mappings: make(chan Whorls),
		Done:     make(chan bool),
	}

	go w.GenerateSource()
	for _, src := range flag.Args() {
		filepath.Walk(src, w.GenerateWhorl())
	}
	close(w.Mappings)
	<-w.Done
}
