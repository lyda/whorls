# Whorls - the simple asset fingerprinter for Go

Whorls will generate a set of fingerprinted files (a hash of the
content makes up the file's name) and a function to use in templates
to match the name to the fingerprinted file.

## Usage

```
Usage of whorls [options] src-dir1 [src-dir2 ... src-dirN]:

Options:
  -C string
    	change to this directory first
  -clear
    	delete contents of <out> dir first (default false)
  -f string
    	generated file name (default "generated-assets.go")
  -not-found string
    	name if requested whorl missing
  -o string
    	directory where assets go (default "assets")
  -p string
    	name of the package (default "assets")
```

## Example

Assuming you have a `pages` package that handles the rendering of your
pages, put this in one of the files:

```go
//go:generate go get -u github.com/lyda/whorls/cmd/whorls
//go:generate whorls -clear -C static -p pages -f ../generated-assets.go -o assets -p pages -not-found missing-whorl-file css js
//go:generate go get -u github.com/mjibson/esc
//go:generate esc -o generated-pages.go -pkg pages static templates

package pages
```

This assumes that with this source file there is a directory called
`static` and in it, directories called `js` and `css`. The `whorls`
utility will:

1. change into `static` (`-C static`),
2. make a new directory called `assets` (`-o assets`),
3. after having deleted `assets` (`-clean`) and
4. make fingerprinted files from the `css` and `js` directories (`css js`).

The generated source file will be called `generated-assets.go`.
The `-f ../generated-assets.go` sets that - remember that the
`-C static` had us change directory first. That source file will
be part of the `pages` package (`-p pages`). The generated function
`WhorlIt` will return `missing-whorl-file` if it can't find an asset
name (`-not-found missing-whorl-file`).
