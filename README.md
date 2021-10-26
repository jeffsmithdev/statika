# Statika

Statika is simple static site generator(SSG) written in go emphasizing convention over configuration.  This is a newer version of my [original statika](https://github.com/jeffsmithdev/statika).
Statika is (generally) not intended for public consumption.  It just does what I want.  It probably won't do what you want.  
That being said....

## Build

* clone
* go get
* go run statika.go examples/simple.com -b

## Usage

```shell
$ statika <project-dir>   # The default action is to watch
$ statika simple.com -c  # Clean the output
$ statika simple.com -b  # Manually build
$ statika simple.com -w  # Watch for changes then build
$ statika simple.com -s  # Run a local development server
```