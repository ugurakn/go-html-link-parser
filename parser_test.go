package link

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

const ex1 = `
<html>
<body>
  <h1>Hello!</h1>
  <a href="/other-page">A link to another page</a>
</body>
</html>
`

const ex2 = `
<html>
<head>
  <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/font-awesome/4.7.0/css/font-awesome.min.css">
</head>
<body>
  <h1>Social stuffs</h1>
  <div>
    <a href="https://www.twitter.com/joncalhoun">
      Check me out on twitter
      <i class="fa fa-twitter" aria-hidden="true"></i>
    </a>
    <a href="https://github.com/gophercises">
      Gophercises is on <strong>Github</strong>!
    </a>
  </div>
</body>
</html>
`

const ex4 = `
<html>
<body>
  <a href="/dog-cat">dog cat <!-- commented text SHOULD NOT be included! --></a>
</body>
</html>
`

var r1 = strings.NewReader(ex1)
var doc1, _ = html.Parse(r1)
var r2 = strings.NewReader(ex2)
var doc2, _ = html.Parse(r2)
var r4 = strings.NewReader(ex4)
var doc4, _ = html.Parse(r4)

func TestGetLinkNodes(t *testing.T) {
	testCases := []struct {
		doc  *html.Node
		name string
		want int
	}{
		{doc1, "single_a", 1},
		{doc2, "multiple_a", 2},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			nodes := getLinkNodes(tc.doc)
			assert.Equal(t, tc.want, len(nodes))
		})
	}
}

func TestExtractText(t *testing.T) {
	testCases := []struct {
		node *html.Node
		name string
		want string
	}{
		{getLinkNodes(doc1)[0], "single", "A link to another page"},
		{getLinkNodes(doc2)[0], "with_extra_after", "Check me out on twitter"},
		{getLinkNodes(doc2)[1], "with_extra_middle", "Gophercises is on Github!"},
		{getLinkNodes(doc4)[0], "with_comment", "dog cat"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, extractText(tc.node))
		})
	}
}

func TestBuildLink(t *testing.T) {
	nodes := getLinkNodes(doc1)

	lnk := buildLink(nodes[0])

	expected := Link{
		Href: "/other-page",
		Text: "A link to another page",
	}

	assert.Equal(t, expected, lnk)
}

func TestParse(t *testing.T) {
	r := strings.NewReader(ex1)
	links, err := Parse(r)
	if assert.NoError(t, err) {
		expected := Link{
			Href: "/other-page",
			Text: "A link to another page",
		}

		assert.Equal(t, 1, len(links))
		assert.Equal(t, expected, links[0])
	}
}
