package models

import (
	"time"
)

type Config struct {
	Verbose       bool
	ProjectDir    string
	SrcDir        string
	StaticDir     string
	TemplatesDir  string
	ContentDir    string
	OutputDir     string
	DevServerHost string
	DevServerPort string
}

type Page struct {
	Content     string
	RawContent  string
	Title       string
	Subtitle    string
	Slug        string
	Thumbnail   string
	Website     string
	Author      string
	Date        time.Time
	Tags        []string
	Categories  []string
	Images      []string
	Id          string
	Uuid        string
	Description string
	Data        map[string]interface{}
}
