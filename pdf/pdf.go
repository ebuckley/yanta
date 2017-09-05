package pdf

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ebuckley/yanta/context"
	"github.com/ebuckley/yanta/site"
)

func printCommand(cmd *exec.Cmd) {
	log.Printf("==> Executing: %s\n", strings.Join(cmd.Args, " "))
}

func hasFile(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// GetPdfCommand creates the cmd for creating a pdf
func GetPdfCommand(url string, path string) *exec.Cmd {
	// TODO make me configurable
	log.Println("downloading", url)
	nodePath, err := exec.LookPath("node")
	if err != nil {
		log.Panic("node not found on path")
	}

	cmd := exec.Command(nodePath, "/Users/Ersin/code/ebuckley/pdf-tool/index.js", "--url="+url, "--path="+path)
	printCommand(cmd)
	cmd.Env = os.Environ()
	return cmd
}

// Download a page pdf file
func Download(cfg *site.Config, page *context.Page) ([]byte, error) {
	pdfPath := filepath.Join(cfg.TmpDir, page.Hash()+".pdf")
	if !hasFile(pdfPath) {
		log.Println("pdf cache miss", page)
		cmd := GetPdfCommand("http://localhost:1337/page/"+page.Path, pdfPath)
		_, err := cmd.CombinedOutput()
		if err != nil {
			return nil, err
		}
	} else {
		log.Println("pdf cache hit!", page)
	}

	return ioutil.ReadFile(pdfPath)
}
