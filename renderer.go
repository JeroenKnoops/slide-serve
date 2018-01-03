package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

const codeMarker = "##CODE##"

func renderSlide(content string, index int) string {
	scanner := bufio.NewScanner(strings.NewReader(content))
	slideMarkup := startSlide(index)
	code := ""

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "  ") || strings.HasPrefix(line, "\t") {
			line = strings.TrimPrefix(line, "  ")
			line = strings.TrimPrefix(line, "\t")
			code += line + "\n"
			slideMarkup += codeMarker
		} else if strings.HasPrefix(line, "#") {
			line = strings.TrimPrefix(line, "#")
			line = headline(line)
			slideMarkup += line
		} else {
			if strings.HasPrefix(line, ".") {
				line = strings.TrimPrefix(line, ".")
			} else {
				if strings.Contains(line, "*") {
					charScanner := bufio.NewScanner(strings.NewReader(line))
					line = ""
					charScanner.Split(bufio.ScanRunes)
					emOpen := false
					for charScanner.Scan() {
						c := charScanner.Text()
						if c == "*" && !emOpen {
							c = "<strong>"
							emOpen = true
						} else if c == "*" && emOpen {
							c = "</strong>"
							emOpen = false
						}
						line += c
					}
				}
			}
			slideMarkup += fmt.Sprintf(`
				<p>%s</p>
			`, line)
		}
	}

	marker := strings.Count(slideMarkup, codeMarker)
	slideMarkup = strings.Replace(slideMarkup, codeMarker, "", marker-1)
	// add highlights https://github.com/alecthomas/chroma
	if code != "" {
		getHighlightedMarkup(code)
		slideMarkup = strings.Replace(slideMarkup, codeMarker, fmt.Sprintf(`
			<pre>%s</pre>
		`, code), 1)
	}

	slideMarkup += endSlide()
	return slideMarkup
}

func headline(txt string) string {
	return fmt.Sprintf(`
		<h1>%s</h1>
	`, txt)
}

func startSlide(index int) string {
	return fmt.Sprintf(`
		<div class="slide slide-%d">
		<div class="slide-content">
	`, index)
}
func endSlide() string {
	return `
		</div>
		</div>
	`
}

func getHighlightedMarkup(code string) {

	lexer := lexers.Analyse(code)

	if lexer == nil {
		log.Println("could not find correct lexer")
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	style := styles.Get("swapoff")
	if style == nil {
		log.Println("using fallback styles")
		style = styles.Fallback
	}

	// formatter := formatters.Get("html")
	formatter := formatters.Get("noop")
	// formatter := html.New()
	if formatter == nil {
		formatter = formatters.Fallback
	}

	iterator, err := lexer.Tokenise(nil, code)
	err = formatter.Format(os.Stdout, style, iterator)
	if err != nil {
		log.Println("err", err)
	}
}