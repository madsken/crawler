package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeURL(t *testing.T) {
	//TEST:remove scheme
	res, err := normalizeURL("https://www.boot.dev/blog/path")
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, "www.boot.dev/blog/path", res)

	//TEST:remove trailing /
	res, err = normalizeURL("www.boot.dev/blog/path/")
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, "www.boot.dev/blog/path", res)

	//TEST:remove scheme and trailing /
	res, err = normalizeURL("http://www.boot.dev/blog/path/")
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, "www.boot.dev/blog/path", res)

	//TEST:convert to lower case
	res, err = normalizeURL("https://CRAWLER-TEST.com/PATH")
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, "crawler-test.com/path", res)

	//TEST:handle scheme, upper case and trailing /
	res, err = normalizeURL("http://CRAWLER-TEST.com/path/")
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, "crawler-test.com/path", res)

	//TEST:invalid url
	res, err = normalizeURL(`:\\invalidURL`)
	require.Error(t, err)
	assert.Equal(t, "", res)
}
