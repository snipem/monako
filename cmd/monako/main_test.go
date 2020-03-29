package main

// run: go test -v ./cmd/monako

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Flaque/filet"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func TestMain(t *testing.T) {

	targetDir := filet.TmpDir(t, "")

	os.Args = []string{
		"monako",
		"-fail-on-error",
		"-target-dir", targetDir,
		"-config", "../../test/configs/only_markdown/config.markdown.yaml",
		"-menu-config", "../../test/configs/only_markdown/config.menu.markdown.md"}
	main()

	t.Run("Check for Hugo input files", func(t *testing.T) {

		assert.FileExists(t, filepath.Join(targetDir, "compose/config.toml"), "Hugo config is not present")
		assert.FileExists(t, filepath.Join(targetDir, "compose/content/monako_menu_directory/index.md"), "Menu is not present")
	})

	t.Run("Check for generated test doc markdown page", func(t *testing.T) {
		assert.FileExists(t, filepath.Join(targetDir, "compose/public/docs/test/test_doc_markdown/index.html"), "Generated Test doc does not exist")

		contentBytes, err := ioutil.ReadFile(filepath.Join(targetDir, "compose/public/docs/test/test_doc_markdown/index.html"))
		content := string(contentBytes)

		assert.NoError(t, err, "Can't read file")
		assert.Contains(t, content, "<strong>Test docs</strong>", "Contains menu")

		assert.Contains(t, content, "<img src=\"../profile.png\" alt=\"Picture in same folder\" />", "Contains relative picture")
		assert.FileExists(t, filepath.Join(targetDir, "compose/public/docs/test/profile.png"), "Relative picture right placed")

		assert.FileExists(t, filepath.Join(targetDir, "compose/public/docs/test/subfolder/subfolderprofile.png"), "Relative subfolder picture right placed")
		assert.Contains(t, content, "<img src=\"../subfolder/subfolderprofile.png\" alt=\"Picture in sub folder\" />", "Contains relative picture")
	})

	fs := http.FileServer(http.Dir(filepath.Join(targetDir, "compose/public/")))
	ts := httptest.NewServer(http.StripPrefix("/", fs))
	defer ts.Close()

	t.Run("Check if images and sources are served", func(t *testing.T) {

		content, err := getContent(ts, "/docs/test/test_doc_markdown/index.html")
		assert.NoError(t, err, "HTTP Call failed")

		urls, err := getURLKeyValuesFromHTML(content, "src", ts.URL)
		if err != nil {
			log.Fatal(err)
		}

		// TODO For some reason imgages seem to be ignored in urls
		for _, url := range urls {
			if strings.HasPrefix(url.String(), ts.URL) {
				// t.Logf("Checking for local served url %s", url.String())
				// Check only if it's served, ignore content
				_, err = getContent(ts, "")
				assert.NoError(t, err, "URL is not served")
			}
		}

	})

	t.Run("Check contents of served page", func(t *testing.T) {

		content, err := getContent(ts, "/docs/test/test_doc_markdown/index.html")
		assert.NoError(t, err, "HTTP Call failed")

		assert.Contains(t, content, "Ihr naht euch wieder, schwankende Gestalten!", "Does not contain Goethe")
		assert.Contains(t, content, "Test docs", "Does not contain Menu header")
		assert.Contains(t, content, "<h3 id=\"markdown-doc-3\">Markdown Doc 3</h3>", "Check rendered Markdown")

	})

	if !t.Failed() {
		// Don't clean up when failed
		filet.CleanUp(t)
	}
}

func getContent(ts *httptest.Server, url string) (string, error) {
	// res, err := http.Get(ts.URL)
	res, err := http.Get(ts.URL + url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	contentBytes, err := ioutil.ReadAll(res.Body)
	return string(contentBytes), err
}

func TestURLKeyValues(t *testing.T) {

	cases := []struct {
		URL, Key, Base, Expected string
	}{
		{"<a href=\"/local/\"></a>", "href", "http://localhost", "http://localhost/local/"},
		{"<img src=\"../image.png\"></a>", "src", "http://localhost/content/docs/", "http://localhost/content/image.png"},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("Extract %s from %s", tc.Key, tc.URL), func(t *testing.T) {
			t.Parallel()
			urls, err := getURLKeyValuesFromHTML(tc.URL, tc.Key, tc.Base)
			assert.NoError(t, err, "Error at extraction")
			assert.Equal(t, tc.Expected, urls[0].String())
		})
	}

}

func getURLKeyValuesFromHTML(content string, key string, baseURL string) ([]*url.URL, error) {

	var urls []*url.URL
	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	z := html.NewTokenizer(strings.NewReader(content))

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document
			return urls, nil
		case tt == html.StartTagToken:
			t := z.Token()

			for _, a := range t.Attr {
				if a.Key == key {
					u, err := url.Parse(a.Val)
					if err != nil {
						log.Fatal(err)
					}
					absoluteURL := base.ResolveReference(u)
					urls = append(urls, absoluteURL)
					break
				}
			}
		}
	}

}
