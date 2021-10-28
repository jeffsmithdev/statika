package util

import (
	"fmt"
	"log"
	"strings"
	"time"
)

func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func GetString(frontMatter map[string]interface{}, key string) string {
	defer handleError()

	if frontMatter[key] != nil {
		return frontMatter[key].(string)
	} else {
		return ""
	}
}

func GetBool(frontMatter map[string]interface{}, key string) bool {
	defer handleError()

	if frontMatter[key] != nil {
		return frontMatter[key].(bool)
	} else {
		return false
	}
}

func GetDate(frontMatter map[string]interface{}, key string) time.Time {
	defer handleError()

	if frontMatter[key] != nil {
		d, err := time.Parse("2006-01-02", frontMatter[key].(string))
		Check(err)
		return d
	} else {
		return time.Now()
	}
}

func GetSlice(frontMatter map[string]interface{}, key string) []string {
	defer handleError()

	if frontMatter[key] != nil {
		return strings.Split(strings.ReplaceAll(frontMatter[key].(string), ", ", ","), ",")
	} else {
		return []string{}
	}
}

func handleError() {
	if err := recover(); err != nil {
		fmt.Println("This is the error: ", err)
	}
}
