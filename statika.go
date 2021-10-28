/*statika

statika is a simple SSG (Static Site Generator) written in go.

Examples:
$ statika <project-dir>
$ statika simple.com -c  # Clean the output
$ statika simple.com -b  # Manually build
$ statika simple.com -w  # Watch for changes then build
$ statika simple.com -s  # Run a local development server
*/
package main

import (
	"fmt"
	_ "github.com/flosch/pongo2-addons"
	"github.com/jessevdk/go-flags"
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
	"statika/models"
	"statika/tasks"
	"statika/util"
)

const version = "0.0.1"

var (
	cfg  *models.Config
	opts = struct {
		SrcDir  string `short:"d" long:"dir" description:"The project directory (if not using SRC_DIR in .env)."`
		Server  []bool `short:"s" long:"server" description:"Run development server"`
		Build   []bool `short:"b" long:"build" description:"Build the site"`
		Watch   []bool `short:"w" long:"watch" description:"Watch for file changes and build site automatically"`
		Clean   []bool `short:"c" long:"clean" description:"Clean build directory"`
		Verbose []bool `short:"v" long:"Verbose" description:"Enable Verbose logging"`
	}{}
)

func init() {
	_, err := flags.Parse(&opts)
	util.Check(err)

	cfg = &models.Config{}
	cfg.Verbose = len(opts.Verbose) > 0

	if opts.SrcDir != "" {
		cfg.ProjectDir = opts.SrcDir
	} else if os.Getenv("SRC_DIR") != "" {
		cfg.ProjectDir = os.Getenv("SRC_DIR")
	} else {
		log.Fatal("The src directory is required.  Please specify either SRC_DIR in the .env file or use the -d parameter on the command line.")
	}

	err = godotenv.Load(cfg.ProjectDir + "/.env")
	util.Check(err)

	cfg.SrcDir = filepath.Join(cfg.ProjectDir, "src/")
	cfg.OutputDir = filepath.Join(cfg.ProjectDir, "output/")
	cfg.StaticDir = filepath.Join(cfg.SrcDir, "static/")
	cfg.TemplatesDir = filepath.Join(cfg.SrcDir, "templates/html")
	cfg.ContentDir = filepath.Join(cfg.SrcDir, "content/")
	cfg.StaticDir = filepath.Join(cfg.SrcDir, "static/")
}

func main() {
	if cfg.Verbose {
		fmt.Println("Version: ", version)
	}

	if len(opts.Server) > 0 {
		fmt.Println("Serving...")
		tasks.Server(cfg)
	} else if len(opts.Watch) > 0 {
		fmt.Println("Watching...")
		tasks.Watch(cfg)
	} else if len(opts.Build) > 0 {
		fmt.Println("Building...")
		tasks.Build(cfg)
	} else if len(opts.Clean) > 0 {
		fmt.Println("Cleaning...")
		tasks.Clean(cfg.OutputDir)
	} else {
		log.Fatal("Please specify a flag to run task: server, watch, build, clean")
	}
}
