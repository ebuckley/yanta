package sync

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/ebuckley/yanta/site"
)

func printCommand(cmd *exec.Cmd) {
	log.Printf("==> Executing: %s\n", strings.Join(cmd.Args, " "))
}

func publish() *exec.Cmd {
	cmd := exec.Command("git", "push", "publish", "master")
	cmd.Env = os.Environ()
	printCommand(cmd)
	return cmd
}

// PublishHandler publishes your notes
func PublishHandler(s *site.Site) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		out, err := publish().CombinedOutput()
		if err != nil {
			log.Println("error publishing", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		fmt.Fprint(w, string(out))
	}
}
