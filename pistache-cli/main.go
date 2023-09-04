package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/akamensky/argparse"
	"github.com/bettercallmolly/molly-cafe/pistache/entities"
	"github.com/bettercallmolly/molly-cafe/pistache/minify"
	"github.com/fsnotify/fsnotify"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

const (
	HTML_BASE = `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="utf-8">
		<title>%s</title>
		<link rel="stylesheet" href="pistache.css">
		<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.8.0/styles/a11y-dark.min.css">
		<script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.8.0/highlight.min.js"></script>
		<meta property="og:title" content="%s">
		<meta property="og:type" content="website">
		<meta property="og:url" content="%s">
		<meta property="og:description" content="%s">
		<meta name="viewport" content="width=device-width,initial-scale=1">
		%s
	</head>
	<body>
		<script>hljs.highlightAll();</script>
		%s
		<footer class='pistache-footer'>
		<h6>
			<a href="https://tailwindcss.com">tailwindcss</a>
			&
			built with <a href="https://github.com/BetterCallMolly">molly</a>'s pistache toolkit
			-
			silk icons by <a href="https://frhun.de/silk-icon-scalable/preview/">frhun</a> (originals by Mark James)
			-
			generated in %.3fµs
		</h6>
		</footer>
	</body>
	`
)

var (
	ADDITIONNAL_LANGUAGE = `<script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.8.0/languages/%s.min.js"></script>`

	Title     *string
	Desc      *string
	InputFile *string
	WatchMode *bool
	URL       string

	// Pistache entities to replace
	entitiesFuncs = map[string]func(*string, string) string{
		"{$read_time}":  entities.ReadTime,
		"{$created_at}": entities.CreatedAt,
		"{$updated_at}": entities.UpdatedAt,
	}
)

// TODO:
//   - Pipe output to minify
func injectLanguages(lang ...string) string {
	if len(lang) == 0 {
		return ""
	}
	var output string
	for _, l := range lang {
		output += fmt.Sprintf(ADDITIONNAL_LANGUAGE, l)
	}
	return output
}

// taken from: https://stackoverflow.com/a/26722698
func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

func formatTitle(title string) string {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	min := strings.ToLower(strings.ReplaceAll(title, " ", "-"))
	s, _, _ := transform.String(t, min)
	return s
}

func pistacheEntities(html *string, path string) {
	// Loop through each entities, if not found, skip and warn
	for entity, f := range entitiesFuncs {
		if !strings.Contains(*html, entity) {
			log.Printf("entity %s not found, skipping", entity)
			continue
		}
		*html = strings.ReplaceAll(*html, entity, f(html, path))
	}
}

// Parses the markdown file and converts it to HTML, minifies if not in watch mode, and saves it to a file
func convertToHTML(md goldmark.Markdown) {
	start := time.Now()
	file, err := os.OpenFile(filepath.Base(URL), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	var buffer bytes.Buffer
	mdData, err := os.ReadFile(*InputFile)
	if err != nil {
		panic(err)
	}
	if err := md.Convert(mdData, &buffer); err != nil {
		panic(err)
	}
	// Get languages from code blocks
	reg := regexp.MustCompile(`<pre><code class=\"language-(\w+)\"`)
	langs := []string{}
	for _, match := range reg.FindAllStringSubmatch(buffer.String(), -1) {
		langs = append(langs, match[1])
	}
	html := fmt.Sprintf(
		HTML_BASE,
		*Title,
		*Title,
		URL,
		*Desc,
		injectLanguages(langs...),
		buffer.String(),
		time.Since(start).Seconds()*1000000,
	)
	pistacheEntities(&html, *InputFile)
	if !*WatchMode {
		log.Printf("Minifying HTML file... ")
		html, err = minify.MinifyHTML(html)
		if err != nil {
			log.Printf("Failed to minify: %s\n", err)
		}
	}
	if _, err := file.WriteString(html); err != nil {
		panic(err)
	}
	log.Println("Output file:", file.Name())
}

func init() {
	args := argparse.NewParser("pistache", "Quickly generate blog posts from markdown files")
	Title = args.String("t", "title", &argparse.Options{Required: true, Help: "Title of the blog post (used for the HTML title and the OG title)"})
	Desc = args.String("d", "description", &argparse.Options{Required: true, Help: "Description of the blog post (used for the OG description)"})
	InputFile = args.String("i", "input", &argparse.Options{Required: true, Help: "Input markdown file"})
	WatchMode = args.Flag("w", "watch", &argparse.Options{Required: false, Help: "Saves the HTML file and watches for changes in the markdown file (does not minify the HTML file)"})
	err := args.Parse(os.Args)
	if err != nil {
		panic(err)
	}
	URL = "https://mana.rip/pistache/" + formatTitle(*Title) + ".html"
	log.Println("Generated URL:", URL)
}

func main() {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithRendererOptions(html.WithUnsafe()),
	)
	md.Parser().AddOptions(
		parser.WithAutoHeadingID(),
	)
	if *WatchMode {
		convertToHTML(md)
		log.Println("Watching markdown file for changes... Press CTRL+C to stop.")
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			panic(err)
		}
		watcher.Add(*InputFile)
		defer watcher.Close()
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					convertToHTML(md)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	} else {
		convertToHTML(md)
	}
}
