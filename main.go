package main

import (
	"flag"
	"github.com/go-logr/logr"
	"github.com/go-logr/stdr"
	"github.com/gokyle/fswatch"
	"github.com/patrickmn/go-cache"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	v     = flag.Int("v", 2, "log level")
	path  = flag.String("path", "", "path to monitor, watching working directory if not set")
	delay = flag.Duration("delay", time.Second, "delay between subsequent filesystem scans")

	lastTouchedCache = cache.New(time.Minute, time.Minute)
)

func main() {
	flag.Parse()
	stdr.SetVerbosity(*v)
	log := stdr.New(log.Default())

	if *path == "" {
		p, err := os.Getwd()
		if err != nil {
			log.Error(err, "failed getting working directory")
			os.Exit(1)
		}
		path = &p
	}
	log = log.WithValues("path", *path)

	fswatch.WatchDelay = *delay
	watcher := fswatch.NewAutoWatcher(*path)
	notifies := watcher.Start()
	log.V(1).Info("starting watch", "delay", *delay)
	for {
		notify, ok := <-notifies
		if !ok {
			log.V(1).Info("watched directory deleted, terminating")
			return
		}

		log := log.WithValues("file", notify.Path, "event", notify.Event)
		log.V(4).Info("observed event on file")
		touchFileOrParent(log, notify.Path)
	}
}

func touchFileOrParent(log logr.Logger, path string) bool {
	log = log.WithValues("touched", path)
	touchParent := func(log logr.Logger, err error, path string) bool {
		// no parent directory left to touch
		if path == "/" {
			log.V(1).Error(err, "failed touching /, dismissing")
			return false
		}
		log.V(1).Error(err, "failed touching, retrying on parent directory")
		return touchFileOrParent(log, filepath.Dir(path))
	}

	// test if it was us setting the timestamp, prevent loops
	if lastTouched, ok := lastTouchedCache.Get(path); ok {
		fileInfo, err := os.Stat(path)
		if err != nil {
			return touchParent(log, err, path)
		}

		if fileInfo.ModTime().Equal(lastTouched.(time.Time)) {
			log.V(3).Info("already touched, dismissing")
			return true
		}
	}

	currentTime := time.Now().Local()
	log.V(4).Info("touching")
	if err := os.Chtimes(path, currentTime, currentTime); err != nil {
		return touchParent(log, err, path)
	}
	lastTouchedCache.Set(path, currentTime, 0)
	log.V(2).Info("successfully touched")
	return true
}
