package glosbe

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/andybalholm/cascadia"
	"github.com/echojc/lernkartei/dict"
	"golang.org/x/net/html"
)

const (
	keywordComparative    = "comparative forms"
	keywordSuperlative    = "superlative forms"
	keywordPredicative    = "predicative"
	keywordConjugation    = "conjugation of"
	keywordPastParticiple = "past participle"
	keywordAuxilliary     = "auxiliary"
	keywordPresent        = "present"
	keywordPreterite      = "preterite"

	baseURL = "https://glosbe.com/de/en"
	//baseURL = "http://localhost:8000"
)

var (
	ErrNoEntries           = errors.New("Could not find any entries")
	ErrUnknownPartOfSpeech = errors.New("Unknown part of speech")
	ErrNoConjugations      = errors.New("Could not find conjugations")
	ErrNoDefinitions       = errors.New("Could not find definitions")

	regexpAdjective       = regexp.MustCompile("er ist (.*)")
	regexpVerbThirdPerson = regexp.MustCompile("er (.*)")
)

func Lookup(word string) ([]dict.Word, error) {
	c := http.Client{
		Timeout: 10 * time.Second,
	}

	res, err := c.Get(fmt.Sprintf("%s/%s", baseURL, word))
	if err != nil {
		return nil, err
	}

	root, err := html.Parse(res.Body)
	if err != nil {
		return nil, err
	}

	entries, err := extractEntries(root)
	if err != nil {
		return nil, err
	}

	var ws []dict.Word
	for _, e := range entries {
		var w dict.Word
		w.PartOfSpeech, _ = extractPartOfSpeech(e)
		if w.PartOfSpeech == "" {
			continue
		}

		w.Base = extractBase(e)
		w.Definitions, _ = extractDefinitions(e)
		w.Forms, _ = extractForms(e)
		ws = append(ws, w)
	}

	return ws, nil
}

func extractBase(e entry) string {
	if e.heading == nil {
		return ""
	}
	return strings.TrimSpace(text(e.heading))
}

func extractDefinitions(e entry) ([]string, error) {
	if e.definitions == nil {
		return nil, ErrNoDefinitions
	}

	sel, err := cascadia.Compile("li strong")
	if err != nil {
		return nil, err
	}

	var out []string
	for _, n := range sel.MatchAll(e.definitions) {
		out = append(out, text(n))
	}

	if len(out) < 4 {
		return out, nil
	}
	return out[:4], nil
}

func extractPartOfSpeech(e entry) (dict.PartOfSpeech, error) {
	switch strings.TrimSpace(e.grammar.FirstChild.Data) {
	case "verb":
		return dict.Verb, nil
	case "noun":
		return dict.Noun, nil
	case "adjective":
		return dict.Adjective, nil
	}

	return "", ErrUnknownPartOfSpeech
}

func extractForms(e entry) ([]string, error) {
	switch strings.TrimSpace(e.grammar.FirstChild.Data) {
	case "adjective":
		return extractAdjectiveForms(e.grammar)
	case "noun":
		return extractNounForms(e.grammar)
	case "verb":
		return extractVerbForms(e.grammar)
	}

	return nil, ErrUnknownPartOfSpeech
}

func extractVerbForms(n *html.Node) ([]string, error) {
	sel, err := cascadia.Compile("tr")
	if err != nil {
		return nil, err
	}

	nextTable := false
	for n = n.FirstChild; n != nil; n = n.NextSibling {
		if strings.Contains(strings.ToLower(n.Data), keywordConjugation) {
			nextTable = true
		}

		if !nextTable || !isTag(n, "table") {
			continue
		}

		// found table, range over rows to find things we want
		index := -1
		out := []string{"", "", " "}
		for _, r := range sel.MatchAll(n) {
			var first, last *html.Node
			for first = r.FirstChild; first.Type != html.ElementNode; first = first.NextSibling {
			}
			for last = r.LastChild; last.Type != html.ElementNode; last = last.PrevSibling {
			}

			keyword := text(first)
			data := strings.TrimSpace(text(last))

			switch true {
			case strings.Contains(keyword, keywordPastParticiple):
				out[2] = data + out[2]
			case strings.Contains(keyword, keywordAuxilliary):
				out[2] = out[2] + data
			case strings.Contains(keyword, keywordPresent):
				index = 0
			case strings.Contains(keyword, keywordPreterite):
				index = 1
			default:
				matches := regexpVerbThirdPerson.FindStringSubmatch(keyword)
				if matches != nil && index >= 0 {
					out[index] = strings.TrimSpace(matches[1])
				}
			}
		}

		for i := len(out) - 1; i >= 0; i-- {
			out[i] = strings.TrimSpace(out[i])
			if out[i] == "" {
				out = append(out[0:i], out[i+1:]...)
			}
		}

		return out, nil
	}

	return nil, ErrNoConjugations
}

func extractNounForms(n *html.Node) ([]string, error) {
	sel, err := cascadia.Compile("td + td")
	if err != nil {
		return nil, err
	}

	// best effort: start from the second td and extract siblings
	forms := []string{}
	ns := sel.MatchAll(n)
	if len(ns) == 0 {
		return forms, nil
	}

	for n, i := ns[0], 0; n != nil; n, i = nextElementSibling(n), i+1 {
		s := strings.TrimSpace(text(n))
		if len(forms) == i/2 {
			forms = append(forms, s)
		} else {
			forms[i/2] += " " + s
		}
	}

	return forms, nil
}

func extractAdjectiveForms(n *html.Node) ([]string, error) {
	sel, err := cascadia.Compile("td")
	if err != nil {
		return nil, err
	}

	forms := make([]string, 2)
	index := -1
	for n = n.FirstChild; n != nil; n = n.NextSibling {
		switch true {
		case strings.Contains(strings.ToLower(n.Data), keywordComparative):
			index = 0
		case strings.Contains(strings.ToLower(n.Data), keywordSuperlative):
			index = 1
		}

		if index < 0 || !isTag(n, "table") {
			continue
		}

		cell := sel.MatchAll(n)[0]
		matches := regexpAdjective.FindStringSubmatch(text(cell))
		if matches == nil {
			continue
		}

		forms[index] = strings.TrimSpace(matches[1])
	}

	for i := len(forms) - 1; i >= 0; i-- {
		if forms[i] == "" {
			forms = append(forms[0:i], forms[i+1:]...)
		}
	}

	return forms, nil
}

type entry struct {
	heading     *html.Node
	definitions *html.Node
	grammar     *html.Node
}

func extractEntries(root *html.Node) ([]entry, error) {
	sel, err := cascadia.Compile("#phraseTranslation .additional-data")
	if err != nil {
		return nil, err
	}

	var es []entry
	for _, d := range sel.MatchAll(root) {
		for n := d.FirstChild; n != nil; n = n.NextSibling {
			switch strings.TrimSpace(n.FirstChild.Data) {
			case "verb", "adjective", "noun":
				e := entry{
					grammar: n,
				}
				for n2 := n.Parent.PrevSibling; n2 != nil && !hasClass(n2, "additional-data"); n2 = n2.PrevSibling {
					if isTag(n2, "ul") {
						e.definitions = n2
					} else if isTag(n2, "h3") {
						e.heading = n2
					}
				}
				es = append(es, e)
			}
		}
	}

	if es == nil {
		return nil, ErrNoEntries
	}

	return es, nil
}

func nextElementSibling(n *html.Node) *html.Node {
	for n = n.NextSibling; n != nil; n = n.NextSibling {
		if n.Type == html.ElementNode {
			return n
		}
	}
	return nil
}

func isTag(n *html.Node, tag string) bool {
	return n.Type == html.ElementNode && n.Data == tag
}

func hasClass(n *html.Node, target string) bool {
	if n.Type != html.ElementNode {
		return false
	}

	for _, attr := range n.Attr {
		if attr.Key == "class" {
			for _, class := range strings.Fields(attr.Val) {
				if class == target {
					return true
				}
			}
			return false
		}
	}

	return false
}

func text(n *html.Node) string {
	var b bytes.Buffer
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		switch c.Type {
		case html.TextNode:
			b.WriteString(c.Data)
		default:
			b.WriteString(text(c))
		}
	}
	return b.String()
}
