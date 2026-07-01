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
