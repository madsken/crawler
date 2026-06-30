package main

import (
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
