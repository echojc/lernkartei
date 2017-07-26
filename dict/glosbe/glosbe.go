package glosbe

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/echojc/lernkartei/dict"
	"github.com/ericchiang/css"
	"golang.org/x/net/html"
)

const (
	keywordComparative    = "Comparative forms"
	keywordSuperlative    = "Superlative forms"
	keywordPredicative    = "predicative"
	keywordConjugation    = "Conjugation of"
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

func NewWord(word string) ([]dict.Word, error) {
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
		w.Base = extractBase(e)
		if w.Base != word {
			continue
		}
		w.Definitions, _ = extractDefinitions(e)
		w.PartOfSpeech, _ = extractPartOfSpeech(e)
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

	sel, err := css.Compile("li strong")
	if err != nil {
		return nil, err
	}

	var out []string
	for _, n := range sel.Select(e.definitions) {
		out = append(out, text(n))
	}

	if len(out) < 4 {
		return out, nil
	}
	return out[:4], nil
}

func extractPartOfSpeech(e entry) (dict.PartOfSpeech, error) {
	if e.grammar == nil {
		return "", ErrUnknownPartOfSpeech
	}

	for n := e.grammar.FirstChild; n != nil; n = n.NextSibling {
		switch strings.TrimSpace(n.FirstChild.Data) {
		case "verb":
			return dict.Verb, nil
		case "noun":
			return dict.Noun, nil
		case "adjective":
			return dict.Adjective, nil
		}
	}

	return "", ErrUnknownPartOfSpeech
}

func extractForms(e entry) ([]string, error) {
	if e.grammar == nil {
		return nil, ErrUnknownPartOfSpeech
	}

	for n := e.grammar.FirstChild; n != nil; n = n.NextSibling {
		switch strings.TrimSpace(n.FirstChild.Data) {
		case "adjective":
			return extractAdjectiveForms(n)
		case "noun":
			return extractNounForms(n)
		case "verb":
			return extractVerbForms(n)
		}
	}

	return nil, ErrUnknownPartOfSpeech
}

func extractVerbForms(n *html.Node) ([]string, error) {
	sel, err := css.Compile("tr")
	if err != nil {
		return nil, err
	}

	nextTable := false
	for n = n.FirstChild; n != nil; n = n.NextSibling {
		if strings.Contains(n.Data, keywordConjugation) {
			nextTable = true
		}

		if !nextTable || !isTag(n, "table") {
			continue
		}

		// found table, range over rows to find things we want
		index := -1
		out := []string{"", "", " "}
		for _, r := range sel.Select(n) {
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
		return out, nil
	}

	return nil, ErrNoConjugations
}

func extractNounForms(n *html.Node) ([]string, error) {
	sel, err := css.Compile("td")
	if err != nil {
		return nil, err
	}

	// best effort extract
	forms := make([]string, 2)
	for i, n := range sel.Select(n) {
		if i > 4 {
			break
		}

		switch i {
		case 1:
			forms[0] = strings.TrimSpace(text(n))
		case 2:
			forms[0] = fmt.Sprintf("%s %s", forms[0], strings.TrimSpace(text(n)))
		case 3:
			forms[1] = strings.TrimSpace(text(n))
		case 4:
			forms[1] = fmt.Sprintf("%s %s", forms[1], strings.TrimSpace(text(n)))
		}
	}

	return forms, nil
}

func extractAdjectiveForms(n *html.Node) ([]string, error) {
	sel, err := css.Compile("td")
	if err != nil {
		return nil, err
	}

	forms := make([]string, 2)
	index := -1
	for n = n.FirstChild; n != nil; n = n.NextSibling {
		switch true {
		case strings.Contains(n.Data, keywordComparative):
			index = 0
		case strings.Contains(n.Data, keywordSuperlative):
			index = 1
		}

		if index < 0 || !isTag(n, "table") {
			continue
		}

		cell := sel.Select(n)[0]
		matches := regexpAdjective.FindStringSubmatch(text(cell))
		if matches == nil {
			continue
		}

		forms[index] = strings.TrimSpace(matches[1])
	}

	return forms, nil
}

type entry struct {
	heading     *html.Node
	definitions *html.Node
	grammar     *html.Node
}

func extractEntries(root *html.Node) ([]entry, error) {
	sel, err := css.Compile("#phraseTranslation h3")
	if err != nil {
		return nil, err
	}

	var es []entry
	for _, n := range sel.Select(root) {
		var e entry
		e.heading = n

		for n = n.NextSibling; n != nil && !isTag(n, "h3"); n = n.NextSibling {
			if isTag(n, "ul") {
				e.definitions = n
			} else if hasClass(n, "additional-data") {
				e.grammar = n
			}
		}

		es = append(es, e)
	}

	if es == nil {
		return nil, ErrNoEntries
	}

	return es, nil
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
