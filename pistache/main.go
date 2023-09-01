package main

import (
	"bytes"
	"fmt"
	"os"
	"regexp"

	"github.com/akamensky/argparse"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
	"go.abhg.dev/goldmark/anchor"
	"go.abhg.dev/goldmark/toc"
)

const (
	HTML_BASE = `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="utf-8">
		<title>%s</title>
		<link rel="stylesheet" href="//cdn.mana.rip/pistache.css">
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
	</body>
	`
)

var (
	ADDITIONNAL_LANGUAGE = `<script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.8.0/languages/%s.min.js"></script>`

	Title *string
	Desc  *string
	URL   *string
)

// TODO:
//	- Inject custom badges at the beginning of the page (tags) -> goldmark extension (custom MD tags)
//  - Silk icons as emojis? (fork goldmark-emoji ?)
//  - Pipe output to minify

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

func init() {
	args := argparse.NewParser("pistache", "Quickly generate blog posts from markdown files")
	Title = args.String("t", "title", &argparse.Options{Required: true, Help: "Title of the blog post (used for the HTML title and the OG title)"})
	Desc = args.String("d", "description", &argparse.Options{Required: true, Help: "Description of the blog post (used for the OG description)"})
	// TODO: deprecate this
	URL = args.String("u", "url", &argparse.Options{Required: false, Help: "URL of the blog post (used for the OG URL)"})
	err := args.Parse(os.Args)
	if err != nil {
		panic(err)
	}
}

func main() {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM, &anchor.Extender{}),
		goldmark.WithRendererOptions(html.WithUnsafe()),
	)
	md.Parser().AddOptions(
		parser.WithAutoHeadingID(),
		parser.WithASTTransformers(
			util.Prioritized(&toc.Transformer{
				Title: "Table of content",
			}, 100),
		),
	)
	file, err := os.OpenFile("test.html", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	var buffer bytes.Buffer
	mdData, err := os.ReadFile("test.md")
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
		*URL,
		*Desc,
		injectLanguages(langs...),
		buffer.String(),
	)
	if _, err := file.WriteString(html); err != nil {
		panic(err)
	}
}
