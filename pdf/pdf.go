package pdf

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ebuckley/marked/context"
)

const (
	// TmpDir is the directory of temporary pdf files
	TmpDir = "~/tmp"
)

func printCommand(cmd *exec.Cmd) {
	log.Printf("==> Executing: %s\n", strings.Join(cmd.Args, " "))
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
func Download(page *context.Page) ([]byte, error) {
	pdfPath := filepath.Join(TmpDir, page.Hash()+".pdf")
	cmd := GetPdfCommand("http://localhost:1337/"+page.Path, pdfPath)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	return ioutil.ReadFile(pdfPath)
}
