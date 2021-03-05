package client

import (
	"errors"
	"fmt"
	"html/template"
	"hummingbard/config"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type BasePage struct {
	Name         string      `json:"name"`
	LoggedInUser interface{} `json:"logged_in_user"`
	Nonce        string      `json:"nonce"`
}

type Template struct {
	*template.Template
}

var fMap = template.FuncMap{
	"InsertJS":       insertJS,
	"InsertCSS":      insertCSS,
	"FormatTime":     formatTime,
	"Map":            mapp,
	"FileSize":       FileSize,
	"ToString":       ToString,
	"IsLastItem":     IsLastItem,
	"StripMXCPrefix": StripMXCPrefix,
	"AspectRatio":    AspectRatio,
	"IsUserProfile":  isUserProfile,
	"RandomString":   randomString,
	"Title":          title,
	"Sum":            sum,
	"Concat":         concat,
	"Trunc":          truncate,
	"HasColon":       hasColon,
	"Repeat":         repeat,
}

func hasColon(s string) bool {
	return strings.Contains(s, ":")
}

func truncate(s string, i int) string {

	runes := []rune(s)
	if len(runes) > i {
		return string(runes[:i])
	}

	return s
}

func concat(values ...string) string {
	return strings.Trim(strings.Join(values, ""), "")
}

func sum(i, g int) int {
	return i + g
}

func repeat(s string, g int) string {
	x := ``
	for i := 0; i < g; i++ {
		x += s
	}
	return x
}

func title(s string) string {
	return strings.Title(s)
}

func randomString(i int) string {
	return RandomString(i)
}

func isUserProfile(s string) bool {
	conf, err := config.Read()
	if err != nil {
		log.Println(err)
	}
	return strings.Contains(s, "@") && strings.Contains(s, conf.Client.Domain)
}

func AspectRatio(x, y string) string {
	height, err := strconv.Atoi(x)
	if err != nil {
		log.Println(err)
	}
	width, err := strconv.Atoi(y)
	if err != nil {
		log.Println(err)
	}

	return fmt.Sprintf(`%d %d`, height, width)
}

func mapp(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dict call")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

func formatTime(t int64) string {
	ut := time.Unix(t, 0)
	log.Println(ut)
	log.Println(ut)
	log.Println(ut)
	return fmt.Sprintf(`%s`, ut)
}

func insertJS(name string) template.HTML {

	root := "static/js"
	files, err := ioutil.ReadDir(root)
	if err != nil {
	}

	var scr string

	for _, file := range files {
		x := strings.Split(file.Name(), ".")
		if name == x[0] && x[len(x)-1] == "js" {
			scr = fmt.Sprintf("/static/js/%s", file.Name())
			return template.HTML(scr)
		}
	}
	scr = fmt.Sprintf("/static/js/%s.js", "missing")
	return template.HTML(scr)
}

func insertCSS(name string) template.HTML {

	root := "static/css"
	files, err := ioutil.ReadDir(root)
	if err != nil {
		panic(err)
	}

	var scr string

	for _, file := range files {
		x := strings.Split(file.Name(), ".")
		if name == x[0] && x[len(x)-1] == "css" {
			scr = fmt.Sprintf("/static/css/%s", file.Name())
			return template.HTML(scr)
		}
	}
	return template.HTML("")
}

func NewTemplate() (*Template, error) {

	tmpl, err := findAndParseTemplates([]interface{}{"templates"}, fMap)
	if err != nil {
		panic(err)
	}

	return tmpl, err
}

func (c *Client) ReloadTemplates() {
	tmpl, err := findAndParseTemplates([]interface{}{"templates"}, fMap)

	if err != nil {
		log.Printf("parsing: %s", err)
	}
	c.Templates = tmpl
}

func (t *Template) execute(wr io.Writer, name string, data interface{}) error {

	pdat := reflect.TypeOf(data)

	newdat := reflect.New(pdat)

	t.ExecuteTemplate(wr, name, newdat)

	return nil
}

func findAndParseTemplates(rootDir interface{}, funcMap template.FuncMap) (*Template, error) {
	root := template.New("")

	tempo := &Template{root}

	var err error

	for _, x := range rootDir.([]interface{}) {
		cleanRoot := filepath.Clean(x.(string))
		pfx := len(cleanRoot) + 1
		err = filepath.Walk(cleanRoot, func(path string, info os.FileInfo, e1 error) error {
			if !info.IsDir() && strings.HasSuffix(path, ".html") {
				if e1 != nil {
					return e1
				}

				b, e2 := ioutil.ReadFile(path)
				if e2 != nil {
					return e2
				}

				name := path[pfx:]

				t := tempo.New(name).Funcs(funcMap)
				t, e2 = t.Parse(string(b))
				if e2 != nil {
					return e2
				}
			}

			return nil
		})
	}

	return tempo, err
}

func IsLastItem(index, length int) bool {
	return index == length-1
}

func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

func FileSize(t float64) string {
	suffixes := []string{"Bytes", "KB", "MB", "GB"}

	base := math.Log(float64(t)) / math.Log(1024)
	getSize := Round(math.Pow(1024, base-math.Floor(base)), .5, 2)
	getSuffix := suffixes[int(math.Floor(base))]

	return strconv.FormatFloat(getSize, 'f', -1, 64) + " " + string(getSuffix)
}

func ToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case uint:
		return fmt.Sprint(v)
	case float64:
		return fmt.Sprint(v)
	default:
		return ""
	}
}
