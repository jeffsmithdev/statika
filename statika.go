/*statika

statika is a simple SSG (Static Site Generator).

Examples:
$ statika <project-dir>   # The project dir should be in the form of a domain name (without subdomain)
$ statika simple.com -c  # Clean the output
$ statika simple.com -b  # Manually build
$ statika simple.com -w  # Watch for changes then build
$ statika simple.com -s  # Run a local development server
*/
package main

import (
	"bytes"
	"fmt"
	"github.com/flosch/pongo2/v4"
	"github.com/gernest/front"
	"github.com/ikeikeikeike/go-sitemap-generator/v2/stm"
	"github.com/jessevdk/go-flags"
	"github.com/joho/godotenv"
	"github.com/otiai10/copy"
	"github.com/radovskyb/watcher"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	minhtml "github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

const version = "0.0.1"

var (
	verbose      bool
	projectDir   string
	srcDir       string
	staticDir    string
	templatesDir string
	contentDir   string
	outputDir    string
	sections     []string
	sm           *stm.Sitemap
	m            *minify.M
)

type item struct {
	Content    string
	RawContent string
	Title      string
	Subtitle   string
	Slug       string
	Thumbnail  string
	Website    string
	Author     string
	Date       time.Time
	Tags       []string
	Categories []string
	Images     []string
}

func init() {

	if len(os.Args) == 1 {
		log.Fatal("The domain name (i.e. project folder) is required")
	}

	projectDir = os.Args[1]

	err := godotenv.Load(projectDir + "/.env")
	check(err)

	srcDir = filepath.Join(projectDir, "src/")
	outputDir = filepath.Join(projectDir, "output/")
	staticDir = filepath.Join(srcDir, "static/")
	templatesDir = filepath.Join(srcDir, "templates/html")
	contentDir = filepath.Join(srcDir, "content/")
	staticDir = filepath.Join(srcDir, "static/")

	sm = stm.NewSitemap(1)
	sm.SetVerbose(true)
	sm.SetDefaultHost("http://" + projectDir)
	sm.SetSitemapsPath("/")
	sm.SetCompress(false)
	sm.SetPublicPath(outputDir)
	sm.Create()

	m = minify.New()
	m.AddFunc("text/html", minhtml.Minify)
	m.AddFunc("text/css", css.Minify)
	m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
}

func main() {
	var opts struct {
		Server  []bool `short:"s" long:"server" description:"Run development server"`
		Build   []bool `short:"b" long:"build" description:"Build the site"`
		Watch   []bool `short:"w" long:"watch" description:"Watch for file changes and build site automatically"`
		Clean   []bool `short:"c" long:"clean" description:"Clean build directory"`
		Verbose []bool `short:"v" long:"verbose" description:"Enable verbose logging"`
	}
	_, err := flags.Parse(&opts)
	check(err)
	verbose = len(opts.Verbose) > 0

	if verbose {
		fmt.Println("Version: ", version)
	}

	if len(opts.Server) > 0 {
		fmt.Println("Serving...")
		server()
	} else if len(opts.Watch) > 0 {
		fmt.Println("Watching...")
		watch()
	} else if len(opts.Build) > 0 {
		fmt.Println("Building...")
		build()
	} else if len(opts.Clean) > 0 {
		fmt.Println("Cleaning...")
		clean()
	} else {
		log.Fatal("Please specify a flag to run task: server, watch, build, clean")
	}
}

func clean() {
	err := os.RemoveAll(outputDir)
	check(err)

	err = os.MkdirAll(outputDir, 0777)
	check(err)
}

func build() {
	start := time.Now()
	var items = make(map[string][]item)

	clean()

	err := copy.Copy(staticDir, outputDir)
	check(err)

	sections = getSections()

	for _, section := range sections {

		if section == "" {
			continue
		}

		if verbose {
			fmt.Println("Section: ", section)
		}

		sectionPath := filepath.Join(contentDir, section)
		var htmlListTpl = pongo2.Must(pongo2.FromFile(filepath.Join(templatesDir, section+"_list.html")))
		var htmlShowTpl = pongo2.Must(pongo2.FromFile(filepath.Join(templatesDir, section+"_show.html")))

		//var txtListTpl = pongo2.Must(pongo2.FromFile(filepath.Join(templatesDir, section+"_list.txt")))
		//var txtShowTpl = pongo2.Must(pongo2.FromFile(filepath.Join(templatesDir, section+"_show.txt")))

		files, err := ioutil.ReadDir(sectionPath)
		check(err)

		for _, file := range files {
			if verbose {
				fmt.Println("Processing file: ", file.Name())
			}

			contents, _ := ioutil.ReadFile(filepath.Join(sectionPath + "/" + file.Name()))
			matter := front.NewMatter()
			matter.Handle("---", front.YAMLHandler)
			frontMatter, body, err := matter.Parse(strings.NewReader(string(contents)))
			check(err)

			if draft := getBool(frontMatter, "draft"); draft == true {
				continue
			}

			var buf bytes.Buffer
			md := goldmark.New(
				goldmark.WithExtensions(extension.GFM),
				goldmark.WithRendererOptions(
					html.WithUnsafe(),
				),
			)
			if err := md.Convert([]byte(body), &buf); err != nil {
				panic(err)
			}

			item := item{
				Content:    buf.String(),
				RawContent: body,
				Slug:       removeExtension(file.Name()),
			}

			item.Title = get(frontMatter, "title")
			item.Subtitle = get(frontMatter, "subtitle")
			item.Website = get(frontMatter, "website")
			item.Thumbnail = get(frontMatter, "thumbnail")
			item.Author = get(frontMatter, "author")
			item.Images = getSlice(frontMatter, "images")
			item.Tags = getSlice(frontMatter, "tags")
			item.Categories = getSlice(frontMatter, "categories")
			item.Date = getDate(frontMatter, "date")
			check(err)

			items[section] = append(items[section], item)

			htmlContents, err := htmlShowTpl.Execute(pongo2.Context{"item": item})
			//txtContents, err := txtShowTpl.Execute(pongo2.Context{"item": item})
			check(err)

			if section == "pages" {
				outputFilePath := filepath.Join(outputDir, removeExtension(file.Name()))
				err = os.MkdirAll(outputFilePath, 0777)
				err = os.WriteFile(filepath.Join(outputFilePath, "index.html"), []byte(minifyHtml(htmlContents)), 0644)
				//err = os.WriteFile(filepath.Join(outputFilePath, "index.txt"), []byte(txtContents), 0644)
				sm.Add(stm.URL{{"loc", "/" + item.Slug}})
			} else {
				outputFilePath := filepath.Join(outputDir, section, removeExtension(file.Name()))
				err = os.MkdirAll(outputFilePath, 0777)
				check(err)
				err = os.WriteFile(filepath.Join(outputFilePath, "index.html"), []byte(minifyHtml(htmlContents)), 0644)
				//err = os.WriteFile(filepath.Join(outputFilePath, "index.txt"), []byte(txtContents), 0644)
				sm.Add(stm.URL{{"loc", "/" + section + "/" + item.Slug}})
			}
			check(err)

		}

		for _, section := range sections {
			sort.Slice(items[section], func(i, j int) bool {
				return items[section][i].Date.After(items[section][j].Date)
			})
		}

		if section != "pages" {
			outputFilePath := filepath.Join(outputDir, section)
			err = os.MkdirAll(outputFilePath, 0777)
			check(err)

			htmlContent, _ := htmlListTpl.Execute(pongo2.Context{"items": items[section]})
			err = os.WriteFile(filepath.Join(outputFilePath, "index.html"), []byte(minifyHtml(htmlContent)), 0644)
			check(err)

			//txtContent, _ := txtListTpl.Execute(pongo2.Context{"items": items[section]})
			//err = os.WriteFile(filepath.Join(outputFilePath, "index.txt"), []byte(txtContent), 0644)
			check(err)
		}
	}

	var htmlHomeTpl = pongo2.Must(pongo2.FromFile(filepath.Join(templatesDir, "home.html")))
	htmlContents, _ := htmlHomeTpl.Execute(pongo2.Context{"items": items})
	err = os.WriteFile(filepath.Join(outputDir, "index.html"), []byte(minifyHtml(htmlContents)), 0644)

	//var txtHomeTpl = pongo2.Must(pongo2.FromFile(filepath.Join(templatesDir, "home.txt")))
	//txtContents, _ := txtHomeTpl.Execute(pongo2.Context{"items": items})
	//err = os.WriteFile(filepath.Join(outputDir, "index.txt"), []byte(minifyHtml(txtContents)), 0644)

	sm.Finalize()
	duration := time.Since(start)
	fmt.Println("Finished building: ", duration)
}

func getSections() []string {
	dirList, err := ioutil.ReadDir(contentDir)
	if err != nil || len(dirList) == 0 {
		log.Fatal(err)
	}

	for _, fi := range dirList {
		sections = append(sections, fi.Name())
	}
	return sections
}

func minifyHtml(content string) string {
	minifiedHtml, err := m.String("text/html", content)
	check(err)
	return minifiedHtml
}

func watch() {
	build()
	w := watcher.New()
	w.SetMaxEvents(1)
	if err := w.AddRecursive(srcDir); err != nil {
		log.Fatalln(err)
	}
	go func() {
		for {
			select {
			case event := <-w.Event:
				fmt.Println("Change detected: ", event)
				build()
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}

func server() {
	fs := http.FileServer(http.Dir("./" + outputDir))
	http.Handle("/", fs)

	var hostname, host, port string

	if host = os.Getenv("SERVER_HOST"); host == "" {
		host = "localhost"
	}

	if port = os.Getenv("SERVER_PORT"); port == "" {
		port = "8001"
	}

	hostname = host + ":" + port

	log.Println("Listening on " + hostname)
	err := http.ListenAndServe(hostname, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func removeExtension(name string) string {
	return strings.TrimSuffix(name, filepath.Ext(name))
}

func get(frontMatter map[string]interface{}, key string) string {
	if frontMatter[key] != nil {
		return frontMatter[key].(string)
	} else {
		return ""
	}
}

func getBool(frontMatter map[string]interface{}, key string) bool {
	if frontMatter[key] != nil {
		return frontMatter[key].(bool)
	} else {
		return false
	}
}

func getDate(frontMatter map[string]interface{}, key string) time.Time {
	if frontMatter[key] != nil {
		d, err := time.Parse("2006-01-02", frontMatter[key].(string))
		check(err)
		return d
	} else {
		return time.Now()
	}
}

func getSlice(frontMatter map[string]interface{}, key string) []string {
	if frontMatter[key] != nil {
		return strings.Split(strings.ReplaceAll(frontMatter[key].(string), ", ", ","), ",")
	} else {
		return []string{}
	}
}
