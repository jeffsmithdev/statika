package tasks

import (
	"bytes"
	"fmt"
	"github.com/flosch/pongo2/v4"
	"github.com/gernest/front"
	"github.com/gosimple/slug"
	"github.com/ikeikeikeike/go-sitemap-generator/v2/stm"
	"github.com/otiai10/copy"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	minhtml "github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"statika/models"
	"statika/util"
	"strings"
	"time"
)

var m *minify.M
var sm *stm.Sitemap

type Pages map[string][]models.Page
type Tags map[string]map[string][]models.Page
type Templates map[string]map[string]map[string]*pongo2.Template // I'm sorry

func init() {
	m = minify.New()
	m.AddFunc("text/html", minhtml.Minify)
	m.AddFunc("text/css", css.Minify)
	m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
}

func Build(cfg *models.Config) {
	start := time.Now()
	pages := make(Pages)
	tags := make(Tags)
	templates := make(Templates)

	sm = stm.NewSitemap(1)
	sm.SetVerbose(true)
	sm.SetDefaultHost(os.Getenv("URL"))
	sm.SetSitemapsPath("/")
	sm.SetCompress(false)
	sm.SetPublicPath(cfg.OutputDir)
	sm.Create()

	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
	)

	Clean(cfg.OutputDir)

	err := copy.Copy(cfg.StaticDir, cfg.OutputDir)
	util.Check(err)

	sections := GetSections(cfg.ContentDir)

	for _, section := range sections {

		if section == "" {
			continue
		}

		if cfg.Verbose {
			fmt.Println("Section: ", section)
		}

		templates = loadTemplates(section, templates, cfg)
		tags[section] = make(map[string][]models.Page)

		sectionPath := filepath.Join(cfg.ContentDir, section)
		files, err := ioutil.ReadDir(sectionPath)
		util.Check(err)

		// create and save a page for each markdown file
		for _, file := range files {
			if cfg.Verbose {
				fmt.Println("Processing file: ", file.Name())
			}

			frontMatter, body := readMarkdownFile(filepath.Join(sectionPath + "/" + file.Name()))

			if draft := util.GetBool(frontMatter, "draft"); draft == true {
				continue
			}

			var buf bytes.Buffer
			if err := md.Convert([]byte(body), &buf); err != nil {
				util.Check(err)
			}

			page := models.Page{
				Content:     buf.String(),
				RawContent:  body,
				Slug:        removeExtension(file.Name()),
				Id:          util.GetString(frontMatter, "id"),
				Uuid:        util.GetString(frontMatter, "uuid"),
				Title:       util.GetString(frontMatter, "title"),
				Subtitle:    util.GetString(frontMatter, "subtitle"),
				Description: util.GetString(frontMatter, "subtitle"),
				Website:     util.GetString(frontMatter, "website"),
				Thumbnail:   util.GetString(frontMatter, "thumbnail"),
				Author:      util.GetString(frontMatter, "author"),
				Images:      util.GetSlice(frontMatter, "images"),
				Tags:        util.GetSlice(frontMatter, "tags"),
				Categories:  util.GetSlice(frontMatter, "categories"),
				Date:        util.GetDate(frontMatter, "date"),
			}

			htmlContents, err := templates[section]["html"]["show"].Execute(pongo2.Context{"page": page})
			util.Check(err)

			var outputFilePath string
			if section == "pages" {
				outputFilePath = filepath.Join(cfg.OutputDir, removeExtension(file.Name()))
				sm.Add(stm.URL{{"loc", "/" + page.Slug}})
			} else {
				outputFilePath = filepath.Join(cfg.OutputDir, section, removeExtension(file.Name()))
				sm.Add(stm.URL{{"loc", "/" + section + "/" + page.Slug}})
			}
			makeDir(outputFilePath)
			writeFile(filepath.Join(outputFilePath, "index.html"), htmlContents)

			pages[section] = append(pages[section], page)
			for _, tag := range page.Tags {
				tags[section][tag] = append(tags[section][tag], page)
			}
		}

		// sort all pages by date desc to list on index pages
		for _, section := range sections {
			sort.Slice(pages[section], func(i, j int) bool {
				return pages[section][i].Date.After(pages[section][j].Date)
			})
		}

		if section != "pages" {
			// write list page for section
			outputFilePath := filepath.Join(cfg.OutputDir, section)
			makeDir(outputFilePath)
			htmlContent, _ := templates[section]["html"]["list"].Execute(pongo2.Context{"pages": pages[section], "tags": tags[section]})
			writeFile(filepath.Join(outputFilePath, "index.html"), minifyHtml(htmlContent))
			util.Check(err)

			// write list page for each tag in this section
			for key, val := range tags[section] {
				tagPath := filepath.Join(outputFilePath, "tags", slug.Make(key))
				makeDir(tagPath)
				htmlContent, _ := templates[section]["html"]["list"].Execute(pongo2.Context{"pages": val, "tags": tags[section]})
				writeFile(filepath.Join(tagPath, "index.html"), minifyHtml(htmlContent))
				util.Check(err)
			}

			// write tag list page containing index and count of all tags
			tagContent, _ := templates[section]["html"]["tags"].Execute(pongo2.Context{"tags": tags[section]})
			writeFile(filepath.Join(outputFilePath, "tags", "index.html"), minifyHtml(tagContent))
			util.Check(err)
		}
	}

	// write the site's home page
	htmlContents, _ := templates["pages"]["html"]["home"].Execute(pongo2.Context{"pages": pages})
	writeFile(filepath.Join(cfg.OutputDir, "index.html"), minifyHtml(htmlContents))

	sm.Finalize()
	duration := time.Since(start)
	fmt.Println("Finished building: ", duration)
}

func readMarkdownFile(path string) (frontMatter map[string]interface{}, body string) {
	contents, err := ioutil.ReadFile(path)
	util.Check(err)
	matter := front.NewMatter()
	matter.Handle("---", front.YAMLHandler)
	frontMatter, body, err = matter.Parse(strings.NewReader(string(contents)))
	util.Check(err)
	return
}

func makeDir(path string) {
	err := os.MkdirAll(path, 0777)
	util.Check(err)
}

func writeFile(path string, contents string) {
	err := os.WriteFile(filepath.Join(path), []byte(minifyHtml(contents)), 0644)
	util.Check(err)
}

func loadTemplates(section string, templates Templates, cfg *models.Config) Templates {

	types := []string{"list", "show", "tags", "home"}
	templates[section] = make(map[string]map[string]*pongo2.Template)
	templates[section]["html"] = make(map[string]*pongo2.Template)
	for _, t := range types {
		templates[section]["html"][t] = loadTemplate(section, t, cfg)
	}

	return templates
}

func loadTemplate(section string, t string, cfg *models.Config) *pongo2.Template {
	var tpl = &pongo2.Template{}
	localTplPath := filepath.Join(cfg.TemplatesDir, section+"_"+t+".html")
	if fileExists(localTplPath) {
		tpl = pongo2.Must(pongo2.FromFile(filepath.Join(cfg.TemplatesDir, section+"_"+t+".html")))
	}
	return tpl
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	var result bool
	if err == nil {
		result = true
	}
	return result
}

func GetSections(path string) []string {
	var sections []string
	dirList, err := ioutil.ReadDir(path)
	if err != nil || len(dirList) == 0 {
		log.Fatal(err)
	}

	for _, fi := range dirList {
		sections = append(sections, fi.Name())
	}
	return sections
}

func removeExtension(name string) string {
	return strings.TrimSuffix(name, filepath.Ext(name))
}

func minifyHtml(content string) string {
	minifiedHtml, err := m.String("text/html", content)
	util.Check(err)
	return minifiedHtml
}
