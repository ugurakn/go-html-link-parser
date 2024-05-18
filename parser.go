package link

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// Link represents a link (<a href="/...">...</a>)
// in an HTML document
type Link struct {
	Href string
	Text string
}

func (l Link) String() string {
	return fmt.Sprintf("[\n\thref: %s\n\ttext: %s\n]\n", l.Href, l.Text)
}

// Parse accepts an HTML string and returns
// a slice of links parsed from it.
func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, fmt.Errorf("error parsing html: %s", err)
	}

	linkNodes := getLinkNodes(doc)

	// extract href and text from linkNodes
	// build links out of them
	var links []Link
	for _, n := range linkNodes {
		links = append(links, buildLink(n))
	}

	return links, nil
}

func buildLink(n *html.Node) Link {
	var l Link
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			l.Href = attr.Val
			break
		}
	}

	t := extractText(n)
	// remove all \n, \t etc and extra whitespace
	re := regexp.MustCompile(`\s+`)
	t = re.ReplaceAllString(t, " ")
	t = strings.TrimSpace(t)
	l.Text = t

	return l
}

func getLinkNodes(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n}
	}

	var ret []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, getLinkNodes(c)...)
	}
	return ret
}

// DFS all descendants of n to find all text nodes
// and concat them in t
func extractText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}

	var ret string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret += extractText(c)
	}
	return ret
}
