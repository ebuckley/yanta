package site

import (
	"io/ioutil"
	"log"

	"github.com/ebuckley/yanta/context"
)

// Site retains the application state
type Site struct {
	Config     *Config
	Pages      map[string]*context.Page
	ConfigPath string
	CanPublish bool
}

// New will create a site state object
func New(opts ...option) *Site {
	cfg := setupConfig(opts...)
	s := new(Site)
	s.Config = cfg
	log.Println("created new site with config", cfg)
	s.Pages = make(map[string]*context.Page)
	s.CanPublish = true
	return s
}

// FetchPages updates the state of the site :)
func (s *Site) FetchPages() error {
	pages, err := context.GetHTMLs(s.Config.SitePath)
	if err != nil {
		return err
	}
	s.Pages = make(map[string]*context.Page)
	for _, p := range pages {
		s.Pages[p.Path] = p
	}
	return nil
}

// FetchPage will return the page, or fetch it and store it on the site
func (s *Site) FetchPage(p string) (*context.Page, error) {
	page := s.Pages[p]
	if page == nil {
		p, err := context.CreatePage(p)
		if err != nil {
			return nil, err
		}
		page = p
	} else {
		err := page.Refresh()
		if err != nil {
			return nil, err
		}
	}
	return page, nil
}

// UpsertPage is a function for helping the user update or create a page
func (s *Site) UpsertPage(path string, content string) (*context.Page, error) {
	err := ioutil.WriteFile(path, []byte(content), 0644)
	if err != nil {
		return nil, err
	}
	page := context.NewPage(path, []byte(content))
	s.Pages[page.Path] = page
	return page, nil
}
