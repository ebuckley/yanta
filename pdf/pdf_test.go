package pdf

import (
	"fmt"
	"log"
	"testing"
)

func TestDownload(t *testing.T) {
	cmd := getPdfCommand("www.github.com")
	dat, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal("could not get ouput", err)
		t.Fail()
	}
	fmt.Print("got all this data!?", dat)
	// err = ioutil.WriteFile("/tmp/dlpdf.pdf", dat, 0644)
	// if err != nil {
	// log.Fatal("could not write file for some reason", err)
	// }
}
