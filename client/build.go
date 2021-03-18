package client

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/fsnotify.v1"
)

func (c *Client) Build() {
	/*
		BuildCSSFiles()

		if c.Config.Mode == "development" {
			go CSSWatcher()
		}
	*/
}

func BuildCSSFiles() {
	css := "static/css"

	if _, err := os.Stat(css); os.IsNotExist(err) {
		os.MkdirAll(css, os.ModePerm)
	} else {
		os.RemoveAll(css)
		os.MkdirAll(css, os.ModePerm)
	}

	root := "ui/css"
	files, err := ioutil.ReadDir(root)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		x := strings.Split(file.Name(), ".")
		x = append(x, "")
		copy(x[2:], x[1:])
		x[1] = RandomString(20)
		y := strings.Join(x, ".")

		from, err := os.Open(root + "/" + file.Name())
		if err != nil {
			log.Fatal(err)
		}
		defer from.Close()

		to, err := os.OpenFile(css+"/"+y, os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer to.Close()

		_, err = io.Copy(to, from)
		if err != nil {
			log.Fatal(err)
		}

		name := filepath.Join(css, y)

		cmd := exec.Command("uglifycss", "--output", name, name)

		err = cmd.Run()
		if err != nil {
			log.Fatalf("cmd.Run() failed with %s\n", err)
		}
	}

}

func rebuildCSS(name string) {
	css := "static/css/"
	root := "ui/css/"

	// Get the filename without extension
	fn := strings.Split(name, ".")[0]

	// Remove all matching files from last build
	files, err := ioutil.ReadDir(css)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		// find matching files
		x := strings.Split(file.Name(), ".")
		if fn == x[0] {
			os.Remove(css + file.Name())
		}
	}

	x := strings.Split(name, ".")
	x = append(x, "")
	copy(x[2:], x[1:])
	x[1] = RandomString(20)
	y := strings.Join(x, ".")
	from, err := os.Open(root + "/" + name)
	if err != nil {
		log.Fatal(err)
	}
	defer from.Close()

	to, err := os.OpenFile(css+"/"+y, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		log.Fatal(err)
	}
}

func CSSWatcher() {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified CSS file:", event.Name)
					x := strings.Split(event.Name, "/")
					rebuildCSS(x[len(x)-1])
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("ui/css/")
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
