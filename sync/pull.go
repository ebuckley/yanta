package sync

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/ebuckley/yanta/site"
)

func pull() *exec.Cmd {
	cmd := exec.Command("git", "pull", "--rebase", "publish", "master")
	cmd.Env = os.Environ()
	printCommand(cmd)
	return cmd
}

// PullHandler publishes your notes
func PullHandler(s *site.Site) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		out, err := pull().CombinedOutput()
		if err != nil {
			log.Println("error publishing", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		fmt.Fprint(w, string(out))
	}
}
