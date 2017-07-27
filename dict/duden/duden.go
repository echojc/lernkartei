package duden

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
	softHyphen = "­"

	keywordGrammar        = "Grammatik"
	keywordNominative     = "Nominativ"
	keywordPastParticiple = "Partizip II"

	keywordNoun      = "Substantiv"
	keywordAdjective = "Adjektiv"
	keywordVerb      = "Verb"
)

var (
	ErrTooManyNodes        = errors.New("Found too many nodes")
	ErrMissingNode         = errors.New("Could not find node")
	ErrUnknownPartOfSpeech = errors.New("Unknown part of speech")
	ErrNoGrammar           = errors.New("Could not find grammar")
	ErrNoAdjectiveForms    = errors.New("Could not extract adjective forms")
	ErrNoNounForms         = errors.New("Could not extract noun forms")
	ErrNoVerbForms         = errors.New("Could not extract verb forms")
	ErrUnknownAuxilliary   = errors.New("Could not find auxilliary verb")

	regexpAdjective     = regexp.MustCompile(`Steigerungsformen:\s*(.*)\s*,\s*(.*)\s*`)
	regexpThirdSingular = regexp.MustCompile(`er/sie/es (.*)`)
	regexpAuxilliary    = regexp.MustCompile(`»(.*)«`)

	baseURL = "http://www.duden.de/rechtschreibung"
)

func Lookup(word string) (w dict.Word, err error) {
	c := http.Client{
		Timeout: 10 * time.Second,
	}

	res, err := c.Get(fmt.Sprintf("%s/%s", baseURL, word))
	if err != nil {
		return
	}

	root, err := html.Parse(res.Body)
	if err != nil {
		return
	}

	w.Base, err = extractBase(root)
	if err != nil {
		return
	}

	w.PartOfSpeech, err = extractPartOfSpeech(root)
	if err != nil {
		return
	}

	w.Forms, err = extractForms(root, w.PartOfSpeech)
	return
}

func extractBase(root *html.Node) (string, error) {
	sel, err := css.Compile("h1")
	if err != nil {
		return "", err
	}

	n, err := exactlyOne(sel, root)
	if err != nil {
		return "", err
	}

	return text(n), nil
}

func extractPartOfSpeech(root *html.Node) (dict.PartOfSpeech, error) {
	sel, err := css.Compile("strong.lexem")
	if err != nil {
		return "", err
	}

	for _, n := range sel.Select(root) {
		s := text(n)
		switch true {
		case strings.Contains(s, keywordNoun):
			return dict.Noun, nil
		case strings.Contains(s, keywordVerb):
			return dict.Verb, nil
		case strings.Contains(s, keywordAdjective):
			return dict.Adjective, nil
		}
	}

	return "", ErrUnknownPartOfSpeech
}

func extractForms(root *html.Node, pos dict.PartOfSpeech) ([]string, error) {
	g, err := grammarNode(root)
	if err != nil {
		return nil, err
	}

	switch pos {
	case dict.Adjective:
		return extractAdjectiveForms(g)
	case dict.Noun:
		return extractNounForms(g)
	case dict.Verb:
		return extractVerbForms(g)
	}

	return nil, ErrNoGrammar
}

func extractAdjectiveForms(g *html.Node) ([]string, error) {
	sel, err := css.Compile(".lexem")
	if err != nil {
		return nil, err
	}

	n, err := exactlyOne(sel, g)
	if err != nil {
		return nil, err
	}

	matches := regexpAdjective.FindStringSubmatch(text(n))
	if matches == nil {
		return nil, ErrNoAdjectiveForms
	}

	return matches[1:], nil
}

func extractNounForms(g *html.Node) ([]string, error) {
	sel, err := css.Compile("tbody tr")
	if err != nil {
		return nil, err
	}

	for _, n := range sel.Select(g) {
		c := n.FirstChild
		if strings.Contains(text(c), keywordNominative) {
			var out []string
			for c = c.NextSibling; c != nil; c = c.NextSibling {
				out = append(out, text(c))
			}
			return out, nil
		}
	}

	return nil, ErrNoNounForms
}

func extractVerbForms(g *html.Node) ([]string, error) {
	aux, err := extractAuxilliary(g)
	if err != nil {
		return nil, err
	}

	sel, err := css.Compile("tbody")
	if err != nil {
		return nil, err
	}

	n := sel.Select(g)
	if len(n) != 3 {
		return nil, ErrNoVerbForms
	}

	out := make([]string, 3)

	// present tense
	for r := n[0].FirstChild; r != nil; r = r.NextSibling {
		s := text(r.FirstChild.NextSibling)
		if matches := regexpThirdSingular.FindStringSubmatch(s); matches != nil {
			out[0] = matches[1]
		}
	}

	// past tense
	for r := n[1].FirstChild; r != nil; r = r.NextSibling {
		s := text(r.FirstChild.NextSibling)
		if matches := regexpThirdSingular.FindStringSubmatch(s); matches != nil {
			out[1] = matches[1]
		}
	}

	// perfect tense
	for r := n[2].FirstChild; r != nil; r = r.NextSibling {
		if strings.Contains(text(r.FirstChild), keywordPastParticiple) {
			out[2] = fmt.Sprintf("%s %s", text(r.LastChild), aux)
		}
	}

	for i := range out {
		if out[i] == "" {
			return out, ErrNoVerbForms
		}
		out[i] = strings.TrimSpace(out[i])
	}

	return out, nil
}

func extractAuxilliary(g *html.Node) (string, error) {
	sel, err := css.Compile(".lexem")
	if err != nil {
		return "", err
	}

	n, err := exactlyOne(sel, g)
	if err != nil {
		return "", ErrUnknownAuxilliary
	}

	matches := regexpAuxilliary.FindStringSubmatch(text(n))
	if matches == nil {
		return "", ErrUnknownAuxilliary
	}

	return normalizeAuxilliary(matches[1])
}

func normalizeAuxilliary(word string) (string, error) {
	if word == "hat" {
		return "haben", nil
	}
	if word == "ist" {
		return "sein", nil
	}
	return "", ErrUnknownAuxilliary
}

func text(n *html.Node) string {
	var b bytes.Buffer
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		switch c.Type {
		case html.TextNode:
			b.WriteString(strings.Replace(c.Data, softHyphen, "", -1))
		default:
			b.WriteString(text(c))
		}
	}
	return b.String()
}

func grammarNode(root *html.Node) (*html.Node, error) {
	sel, err := css.Compile("h2")
	if err != nil {
		return nil, err
	}

	for _, n := range sel.Select(root) {
		if strings.Contains(text(n), keywordGrammar) {
			// search up for the nearest <section>
			for !(n == nil || (n.Data == "section" && n.Type == html.ElementNode)) {
				n = n.Parent
			}
			return n, nil
		}
	}

	return nil, ErrNoGrammar
}

func exactlyOne(sel *css.Selector, n *html.Node) (*html.Node, error) {
	ns := sel.Select(n)
	if len(ns) > 1 {
		return nil, ErrTooManyNodes
	} else if len(ns) == 0 {
		return nil, ErrMissingNode
	}
	return ns[0], nil
}
