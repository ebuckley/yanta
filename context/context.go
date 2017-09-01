package context

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-commonmark/markdown"
)

// Page is represents a markdown file, parsed into an HTML string, with a Path
type Page struct {
	Path string
	Html string
}

// Hash returns a unique identefier for this page
func (p *Page) Hash() string {
	h := sha1.New()
	io.WriteString(h, p.Html)
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

func (p *Page) String() string {
	return p.Path
}

// GetHTMLs is a function for building the whole page state thingy
func GetHTMLs(sitePath string) ([]*Page, error) {
	siteDir, err := os.Open(sitePath)
	if err != nil {
		return nil, err
	}

	var data []*Page

	err = eachFile(siteDir, sitePath, func(f *os.File, p string) error {
		pg := new(Page)
		html, err := getHTML(f, p)
		if err != nil {
			return err
		}
		pg.Html = html
		pg.Path = p
		data = append(data, pg)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return data, nil
}

func getHTML(f *os.File, fPath string) (string, error) {

	dat, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}
	md := markdown.New(markdown.HTML(true), markdown.Breaks(true))

	return md.RenderToString(dat), nil
}

// FileFunc is the action for iterating through all files
type FileFunc func(*os.File, string) error

// IterationError represents the net error output from the eachFile function
type IterationError []error

func (i IterationError) Error() string {
	return fmt.Sprintf("found a few iteration errors: %d", len(i))
}

func eachFile(folder *os.File, path string, fn FileFunc) (err error) {
	files, err := folder.Readdir(-1)
	if err != nil {
		return err
	}

	var errors = make([]error, 0)

	for _, file := range files {
		var childPath = filepath.Join(path, file.Name())
		f, err := os.Open(childPath)
		if err != nil {
			errors = append(errors, err)
		}
		if file.IsDir() == true {
			err = eachFile(f, childPath, fn)
			if err != nil {
				errors = append(errors, err)
			}
		} else {
			err = fn(f, childPath)
			if err != nil {
				errors = append(errors, err)
			}
		}
	}
	if len(errors) > 0 {
		log.Println("Error when iterating site files", errors)
		return IterationError(errors)
	}
	return nil
}
