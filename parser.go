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
	return fmt.Sprintf("href: %s\ntext: %s\n", l.Href, l.Text)
}

// Parse accepts an HTML string and returns
// a slice of links parsed from it.
func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, fmt.Errorf("error parsing html: %s", err)
	}

	links := make([]Link, 0)

	parseNode(doc, &links)

	return links, nil
}

func parseNode(n *html.Node, links *[]Link) {
	if n.Type == html.ElementNode && n.Data == "a" {
		l := Link{}
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				l.Href = attr.Val
				break
			}
		}
		// search all descendants of this a tag to find all text nodes
		// and concat them in t
		var t string
		extractText(n, &t)
		// replace all \n, \t and extra whitespace
		re := regexp.MustCompile(`\s+`)
		t = re.ReplaceAllString(t, " ")
		t = strings.TrimSpace(t)
		l.Text = t
		*links = append(*links, l)
		return
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		parseNode(c, links)
	}
}

// DFS all descendants of n to find all text nodes
// and concat them in t
func extractText(n *html.Node, t *string) {
	if n.Type == html.TextNode {
		*t += n.Data
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractText(c, t)
	}
}

// func ParseHTML(r io.Reader) (*html.Node, error) {
// 	doc, err := html.Parse(r)
// 	if err != nil {
// 		return nil, fmt.Errorf("error parsing html: %s", err)
// 	}

// 	return doc, nil
// }

// // PrintNode prints the node and all its descendants
// func PrintNode(n *html.Node, padding string) {
// 	fmt.Printf("%s%s\n", padding, n.Data)

// 	for c := n.FirstChild; c != nil; c = c.NextSibling {
// 		PrintNode(c, padding+"  ")
// 	}
// }
