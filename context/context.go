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
	Path    string
	Content string
	Html    string
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

// NewPage is a constructor for a page type
func NewPage(path string, content []byte) *Page {
	p := new(Page)
	p.Path = path
	p.Content = string(content)
	p.Html = getHTML(content)
	return p
}

// CreatePage is a function that fetches a page or returns an error
func CreatePage(path string) (*Page, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return NewPage(path, content), nil
}

// Refresh will fetch the updated content for this file
func (p *Page) Refresh() error {
	f, err := os.Open(p.Path)
	if err != nil {
		return err
	}
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	p.Content = string(content)
	p.Html = getHTML(content)
	return nil
}

// GetHTMLs is a function for building the whole page state thingy
func GetHTMLs(sitePath string) ([]*Page, error) {
	siteDir, err := os.Open(sitePath)
	if err != nil {
		return nil, err
	}

	var data []*Page

	err = eachFile(siteDir, sitePath, func(f *os.File, p string) error {
		content, err := getContent(f, p)
		if err != nil {
			return err
		}
		pg := NewPage(p, content)
		data = append(data, pg)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return data, nil
}

func getContent(f *os.File, fPath string) ([]byte, error) {
	return ioutil.ReadAll(f)
}
func getHTML(data []byte) string {
	md := markdown.New(markdown.HTML(true), markdown.Breaks(true))
	return md.RenderToString(data)
}

// FileFunc is the action for iterating through all files
type FileFunc func(*os.File, string) error

// IterationError represents the net error output from the eachFile function
type IterationError []error

func (i IterationError) Error() string {
	return fmt.Sprintf("found a few iteration errors: %d", len(i))
}

func eachFile(folder *os.File, path string, fn FileFunc) (err error) {
	if folder.Name() == ".git" {
		return nil
	}

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
