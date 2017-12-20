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

TODO...
