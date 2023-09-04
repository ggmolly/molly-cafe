package minify

import (
	"os/exec"
	"strings"
)

// Pipes the HTML to minify command (if it exists)
// returns the passed HTML if minify command does not exist
func MinifyHTML(html string) (string, error) {
	cmd := exec.Command("minify", "--html")
	cmd.Stdin = strings.NewReader(html)
	out, err := cmd.Output()
	if err != nil {
		return html, err
	}
	return string(out), nil
}
