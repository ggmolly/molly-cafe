package main

import (
	"bytes"
	"os"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/util"
	"go.abhg.dev/goldmark/anchor"
	"go.abhg.dev/goldmark/toc"
)

// TODO:
//	- Inject Pistache's CSS
//	- Syntax highlighting
//	- Inject silk icons
//	- Inject custom badges at the beginning of the page (tags)

func main() {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM, &anchor.Extender{}),
	)
	md.Parser().AddOptions(
		parser.WithAutoHeadingID(),
		parser.WithASTTransformers(
			util.Prioritized(&toc.Transformer{
				Title: "Contents",
			}, 100),
		),
	)
	var buf bytes.Buffer
	mdData, err := os.ReadFile("test.md")
	if err != nil {
		panic(err)
	}
	if err := md.Convert(mdData, &buf); err != nil {
		panic(err)
	}

	os.WriteFile("test.html", buf.Bytes(), 0644)
}
