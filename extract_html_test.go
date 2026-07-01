package main

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadingExtraction(t *testing.T) {
	//TEST:html basic
	inputBody := "<html><body><h1>Test Title</h1></body></html>"
	res := getHeadingFromHTML(inputBody)
	require.NotNil(t, res)
	assert.Equal(t, "Test Title", res)

	//TEST:html no heading
	inputBody = `<html><body><p></p></body></html>`
	res = getHeadingFromHTML(inputBody)
	require.NotNil(t, res)
	assert.Equal(t, "", res)

	//TEST:html h2 fallback
	inputBody = "<html><body><h2>Fallback Title</h2></body></html>"
	res = getHeadingFromHTML(inputBody)
	require.NotNil(t, res)
	assert.Equal(t, "Fallback Title", res)

	//TEST:html chose h1
	inputBody = "<html><body><h1>yes</h1><h2>Fallback Title</h2></body></html>"
	res = getHeadingFromHTML(inputBody)
	require.NotNil(t, res)
	assert.Equal(t, "yes", res)
}

func TestParageaphExtraction(t *testing.T) {
	//TEST:simple html input
	inputBody := `<html><body>
		<p>Outside paragraph.</p>
		<main>
			<p>Main paragraph.</p>
		</main>
	</body></html>`
	res := getFirstParagraphFromHTML(inputBody)
	require.NotNil(t, res)
	assert.Equal(t, "Main paragraph.", res)

	//TEST:simple html input
	inputBody = `<html><body>
		<p>First paragraph outside main.</p>
		<p>Second paragraph outside main.</p>
	</body></html>`
	res = getFirstParagraphFromHTML(inputBody)
	require.NotNil(t, res)
	assert.Equal(t, "First paragraph outside main.", res)
}

func TestGetURLsFromHTML(t *testing.T) {
	//TEST: Abs url
	inputURL := "https://crawler-test.com"
	inputBody := `<html><body><a href="https://crawler-test.com"><span>Boot.dev</span></a></body></html>`
	baseURL, err := url.Parse(inputURL)
	require.NoError(t, err)
	urls, err := getURLsFromHTML(inputBody, baseURL)
	require.NoError(t, err)
	assert.Equal(t, []string{"https://crawler-test.com"}, urls)

	//TEST: rel url
	inputURL = "https://crawler-test.com"
	inputBody = `<html><body><a href="/path/one"><span>Boot.dev</span></a></body></html>`
	baseURL, err = url.Parse(inputURL)
	require.NoError(t, err)
	urls, err = getURLsFromHTML(inputBody, baseURL)
	require.NoError(t, err)
	assert.Equal(t, []string{"https://crawler-test.com/path/one"}, urls)

	//TEST: rel and abs url
	inputURL = "https://crawler-test.com"
	inputBody = `
<html>
	<body>
		<a href="/path/one">
			<span>Boot.dev</span>
		</a>
		<a href="https://other.com/path/one">
			<span>Boot.dev</span>
		</a>
	</body>
</html>`
	baseURL, err = url.Parse(inputURL)
	require.NoError(t, err)
	urls, err = getURLsFromHTML(inputBody, baseURL)
	require.NoError(t, err)
	assert.Equal(t, []string{"https://crawler-test.com/path/one", "https://other.com/path/one"}, urls)

	//TEST: no href
	inputURL = "https://crawler-test.com"
	inputBody = `<html>
	<body>
		<a>
			<span>Boot.dev</span>
		</a>
	</body>
</html>`
	baseURL, err = url.Parse(inputURL)
	require.NoError(t, err)
	urls, err = getURLsFromHTML(inputBody, baseURL)
	require.NoError(t, err)
	assert.Nil(t, urls)

	//TEST: bad html
	inputURL = "https://crawler-test.com"
	inputBody = `<html body>
	<a href="path/one">
		<span>Boot.dev</span>
	</a>
</html body>`
	baseURL, err = url.Parse(inputURL)
	require.NoError(t, err)
	urls, err = getURLsFromHTML(inputBody, baseURL)
	require.NoError(t, err)
	assert.Equal(t, []string{"https://crawler-test.com/path/one"}, urls)

	//TEST: invalid href url
	inputURL = "https://crawler-test.com"
	inputBody = `<html>
	<body>
		<a href=":\\invalidURL">
			<span>Boot.dev</span>
		</a>
	</body>
</html>`
	baseURL, err = url.Parse(inputURL)
	require.NoError(t, err)
	urls, err = getURLsFromHTML(inputBody, baseURL)
	require.NoError(t, err)
	assert.Nil(t, urls)
}

func TestGetIMGsFromHTML(t *testing.T) {
	//TEST: abs path
	inputURL := "https://crawler-test.com"
	inputBody := `<html><body><img src="https://crawler-test.com/logo.png" alt="Logo"></body></html>`
	baseURL, err := url.Parse(inputURL)
	require.NoError(t, err)
	urls, err := getImagesFromHTML(inputBody, baseURL)
	require.NoError(t, err)
	assert.Equal(t, []string{"https://crawler-test.com/logo.png"}, urls)

	//TEST: rel path
	inputURL = "https://crawler-test.com"
	inputBody = `<html><body><img src="/logo.png" alt="Logo"></body></html>`
	baseURL, err = url.Parse(inputURL)
	require.NoError(t, err)
	urls, err = getImagesFromHTML(inputBody, baseURL)
	require.NoError(t, err)
	assert.Equal(t, []string{"https://crawler-test.com/logo.png"}, urls)

	//TEST: multiple imgs
	inputURL = "https://crawler-test.com"
	inputBody = `<html><body>
		<img src="/logo.png" alt="Logo">
		<img src="https://cdn.boot.dev/banner.jpg">
	</body></html>`
	baseURL, err = url.Parse(inputURL)
	require.NoError(t, err)
	urls, err = getImagesFromHTML(inputBody, baseURL)
	require.NoError(t, err)
	assert.Equal(t, []string{"https://crawler-test.com/logo.png", "https://cdn.boot.dev/banner.jpg"}, urls)
}

func TestGetPageData(t *testing.T) {
	// TEST: basic html body and url
	inputURL := "https://crawler-test.com"
	inputBody := `<html>
  <body>
    <h1>Hello World</h1>
    <main><p>First paragraph inside main.</p></main>
    <a href="/about">About</a>
    <img src="/logo.png" alt="Logo">
  </body>
</html>`
	pageData := extractPageData(inputBody, inputURL)
	require.NotNil(t, pageData)
	assert.Equal(t, PageData{
		URL:            "https://crawler-test.com",
		Heading:        "Hello World",
		FirstParagraph: "First paragraph inside main.",
		OutgoingLinks:  []string{"https://crawler-test.com/about"},
		ImageURLs:      []string{"https://crawler-test.com/logo.png"},
	}, pageData)

	// TEST: fallback paragraph when no <main>
	inputURL = "https://crawler-test.com"
	inputBody = `<html>
  <body>
    <h1>Title</h1>
    <p>Outside paragraph wins.</p>
    <a href="/x">x</a>
    <img src="/img.png">
  </body>
</html>`
	pageData = extractPageData(inputBody, inputURL)
	require.NotNil(t, pageData)
	assert.Equal(t, PageData{
		URL:            "https://crawler-test.com",
		Heading:        "Title",
		FirstParagraph: "Outside paragraph wins.",
		OutgoingLinks:  []string{"https://crawler-test.com/x"},
		ImageURLs:      []string{"https://crawler-test.com/img.png"},
	}, pageData)

	// TEST: malformed HTML still parsed; absolute link and image
	inputURL = "https://crawler-test.com"
	inputBody = `<html body>
  <h1>Messy</h1>
  <a href="https://other.com/path">Other</a>
  <img src="https://cdn.boot.dev/banner.jpg">
</html body>`
	pageData = extractPageData(inputBody, inputURL)
	require.NotNil(t, pageData)
	assert.Equal(t, PageData{
		URL:            "https://crawler-test.com",
		Heading:        "Messy",
		FirstParagraph: "",
		OutgoingLinks:  []string{"https://other.com/path"},
		ImageURLs:      []string{"https://cdn.boot.dev/banner.jpg"},
	}, pageData)

	// TEST: no h1 and no paragraph
	inputURL = "https://crawler-test.com"
	inputBody = `<html>
  <body>
    <a href="/only-link">Only link</a>
    <img src="/only.png">
  </body>
</html>`
	pageData = extractPageData(inputBody, inputURL)
	require.NotNil(t, pageData)
	assert.Equal(t, PageData{
		URL:            "https://crawler-test.com",
		Heading:        "",
		FirstParagraph: "",
		OutgoingLinks:  []string{"https://crawler-test.com/only-link"},
		ImageURLs:      []string{"https://crawler-test.com/only.png"},
	}, pageData)

	// TEST: multiple links and images preserve order
	inputURL = "https://crawler-test.com"
	inputBody = `<html><body>
  <h1>t</h1>
  <main><p>p</p></main>
  <a href="/a1">a1</a>
  <a href="https://x.dev/a2">a2</a>
  <img src="/i1.png">
  <img src="https://x.dev/i2.png">
</body></html>`
	pageData = extractPageData(inputBody, inputURL)
	require.NotNil(t, pageData)
	assert.Equal(t, PageData{
		URL:            "https://crawler-test.com",
		Heading:        "t",
		FirstParagraph: "p",
		OutgoingLinks: []string{
			"https://crawler-test.com/a1",
			"https://x.dev/a2",
		},
		ImageURLs: []string{
			"https://crawler-test.com/i1.png",
			"https://x.dev/i2.png",
		},
	}, pageData)

	// TEST: invalid base URL → empty link/image slices
	inputURL = `:\\invalidBaseURL`
	inputBody = `<html>
  <body>
    <h1>Title</h1>
    <p>Paragraph</p>
    <a href="/path">path</a>
    <img src="/logo.png">
  </body>
</html>`
	pageData = extractPageData(inputBody, inputURL)
	require.NotNil(t, pageData)
	assert.Equal(t, PageData{
		URL:            `:\\invalidBaseURL`,
		Heading:        "Title",
		FirstParagraph: "Paragraph",
		OutgoingLinks:  nil,
		ImageURLs:      nil,
	}, pageData)
}
