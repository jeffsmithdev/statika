# Statika

Statika is simple static site generator(SSG) written in go emphasizing convention over configuration.  This is a newer version of my [original statica](https://github.com/jeffsmithdev/statica).
Statika is (generally) not intended for public consumption.  It just does what I want.  It probably won't do what you want.  That being said....

The statika.app example site is live here: https://www.statika.app/

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

## Example project folder

The following is an example project folder:

```bash
── .env
── src
    │── content
    │    │── blog
    │    │   ├── coffee-blog-one.md
    │    │   ├── coffee-blog-three.md
    │    │   └── coffee-blog-two.md
    │    └── pages
    │        └── about.md
    │        └── products.md
    ├── static
    │   └── favicon.ico
    │   └── humans.txt
    │   └── robots.txt
    └── templates
        └── html
            ├── blog_list.html
            ├── blog_show.html
            ├── blog_tags.html
            ├── layout.html
            ├── pages_home.html
            ├── pages_list.html
            └── pages_show.html
```
