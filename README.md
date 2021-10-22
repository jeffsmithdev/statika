# Simpler Statica

A Simple(r) static site generator written in go. It's not intended for public consumption. 
It's a newer version of the [original statica](https://github.com/jeffsmithdev/statica).

## Usage

```shell
$ statica <project-dir>   # The default action is to watch
$ statica example.com -c  # Clean the output
$ statica example.com -b  # Manually build
$ statica example.com -w  # Watch for changes then build
$ statica example.com -s  # Run a local development server
```