package tasks

import (
	"fmt"
	"github.com/radovskyb/watcher"
	"log"
	"statika/models"
	"statika/util"
	"time"
)

func Watch(cfg *models.Config) {
	Build(cfg)
	w := watcher.New()
	w.SetMaxEvents(1)
	err := w.AddRecursive(cfg.SrcDir)
	util.Check(err)
	err = w.AddRecursive(cfg.StaticDir)
	util.Check(err)

	go func() {
		for {
			select {
			case event := <-w.Event:
				fmt.Println("Change detected: ", event)
				Build(cfg)
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	err = w.Start(time.Millisecond * 100)
	util.Check(err)
}
