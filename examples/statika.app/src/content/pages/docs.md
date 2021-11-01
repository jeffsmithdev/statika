---
title: Docs
date: 2021-10-20
draft: false
---

## Intro

Statika is a simple static site generator written in go.  It uses convention over configuration, so it expects the 
directory structure and naming to match specific conventions. Because of this, the only config needed is to
specify a domain and url.

## Usage

```shell
$ statika -d <project-dir>   # The default action is to watch
$ statika -d simple.com -c  # Clean the output
$ statika -d simple.com -b  # Manually build
$ statika -d simple.com -w  # Watch for changes then build
$ statika -d simple.com -s  # Run a local development server
```

## Configuration

The config is specified in an .env file in the root of the project dir.  There are not many options but those that do
exist include:

```bash
# Required
DOMAIN="statika.app"
URL="https://statika.app"

# Optional
DEV_SERVER_HOST="localhost"
DEV_SERVER_PORT=8001
```

## Project Structure

The project folder assumes the following structure where different sections, i.e. blog, events. news, products, etc,
of the site will exist as folders under the content directory.  The pages folder is a special section and produces
content at the root of the site rather than within a section.  If a template file for a section, for example
pages_tags.html, does not exist, statika just skips that output.

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
