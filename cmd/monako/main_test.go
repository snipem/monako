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
	"runtime"
	"strings"
	"testing"

	"github.com/Flaque/filet"
	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func TestMainMonakoTest(t *testing.T) {

	targetDir := GetLocalTempDir(t)
	err := os.Chdir(targetDir)
	assert.NoError(t, err)

	monakoConfig, menuConfig := writeConfig("https://github.com/snipem/monako-test.git")

	os.Args = []string{
		"monako",
		"-fail-on-error",
		"-working-dir", targetDir,
		"-config", monakoConfig,
		"-menu-config", menuConfig}
	t.Logf("Running Monako with %s", os.Args)
	main()

	t.Run("Check for Hugo input files", func(t *testing.T) {

		assert.FileExists(t, filepath.Join(targetDir, "compose/config.toml"), "Hugo config is not present")
		assert.FileExists(t, filepath.Join(targetDir, "compose/content/monako_menu_directory/index.md"), "Menu is not present")
	})

	t.Run("Check for whitelisted files", func(t *testing.T) {
		assert.FileExists(t, filepath.Join(targetDir, "compose/content/docs/test/test_doc_asciidoc.adoc"))
	})

	t.Run("Check for blacklisted files", func(t *testing.T) {
		assert.NoFileExists(t, filepath.Join(targetDir, "compose/content/docs/test/test_doc_asciidoc_include_me.adoc"))
	})

	t.Run("Check for Frontmatter Markdown", func(t *testing.T) {
		contentBytes, err := ioutil.ReadFile(filepath.Join(targetDir, "compose/content/docs/test/test_doc_markdown.md"))
		assert.NoError(t, err)
		content := string(contentBytes)
		assert.Contains(t, content, "MonakoGitRemote: ")
	})

	t.Run("Check for Frontmatter Asciidoc", func(t *testing.T) {
		contentBytes, err := ioutil.ReadFile(filepath.Join(targetDir, "compose/content/docs/test/test_doc_asciidoc.adoc"))
		assert.NoError(t, err)
		content := string(contentBytes)
		assert.Contains(t, content, "MonakoGitRemote: ")
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

	// Provide the public folder over a webserver
	fs := http.FileServer(http.Dir(filepath.Join(targetDir, "compose/public/")))
	ts := httptest.NewServer(http.StripPrefix("/", fs))
	defer ts.Close()

	t.Run("Check if images and sources are served", func(t *testing.T) {

		content, err := getContentFromURL(ts, "/docs/test/test_doc_markdown/index.html")
		assert.NoError(t, err, "HTTP Call failed")

		srcs, err := getURLKeyValuesFromHTML(content, "src", ts.URL)
		if err != nil {
			log.Fatal(err)
		}
		hrefs, err := getURLKeyValuesFromHTML(content, "href", ts.URL)
		if err != nil {
			log.Fatal(err)
		}

		assert.NotNil(t, len(srcs), fmt.Sprintf("No links found in %s", ts.URL))
		assert.NotNil(t, len(srcs), fmt.Sprintf("No images found in %s", ts.URL))

		for _, url := range append(srcs, hrefs...) {
			if strings.HasPrefix(url.String(), ts.URL) {
				// Check only if it's served, ignore content
				_, err = getContentFromURL(ts, "")
				assert.NoError(t, err)
				t.Logf("Url %s is served correctly", url.String())
			}
		}

	})

	t.Run("Check contents of served page markdown", func(t *testing.T) {

		content, err := getContentFromURL(ts, "/docs/test/test_doc_markdown/index.html")
		assert.NoError(t, err, "HTTP Call failed")

		assert.Contains(t, content, "Ihr naht euch wieder, schwankende Gestalten!", "Does not contain Goethe")
		assert.Contains(t, content, "Test docs", "Does not contain Menu header")
		assert.Contains(t, content, "<h3 id=\"markdown-doc-3\">Markdown Doc 3</h3>", "Check rendered Markdown")

	})

	t.Run("Check contents of served page asciidoc", func(t *testing.T) {

		content, err := getContentFromURL(ts, "/docs/test/test_doc_asciidoc/index.html")
		assert.NoError(t, err, "HTTP Call failed")

		assert.Contains(t, content, "Ihr naht euch wieder, schwankende Gestalten!", "Does not contain Goethe")
		assert.Contains(t, content, "Test docs", "Does not contain Menu header")
		assert.Contains(t, content, "<a href=\"#_asciidoc_second_level\">Asciidoc Second Level</a>", "Check rendered Asciidoc")

		// TODO Remove on asciidoctor-fix removal in theme
		t.Run("Check for neccessary Asciidoc Fix Element", func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
			assert.NoError(t, err)
			assert.Len(t, doc.Find("body > main > aside.book-toc").Nodes, 1, "Asciidoc Page is missing aside.book-toc that is necessary for asciidoctor-fix")
		})

	})

	t.Run("Check for RSS feed", func(t *testing.T) {
		content, err := getContentFromURL(ts, "/docs/index.xml")
		assert.NoError(t, err, "HTTP Call failed")

		fp := gofeed.NewParser()
		feed, err := fp.ParseString(content)
		assert.NoError(t, err)

		assert.Equal(t, "Docs on Local Test Page", feed.Title)
		assert.Equal(t, "Monako", feed.Generator)

		assert.Len(t, feed.Items, 4, "Does not contain all matching documents from this test repo")

	})

	if !t.Failed() && runtime.GOOS != "windows" {
		// Only clean up when not failed
		// and not on Windows this is because of a filet bug (https://github.com/Flaque/filet/issues/3)
		filet.CleanUp(t)
	}
}

func TestMainSplitCalls(t *testing.T) {

	targetDir := GetLocalTempDir(t)
	err := os.Chdir(targetDir)
	assert.NoError(t, err)

	monakoConfig, menuConfig := writeConfig("https://github.com/snipem/monako-test.git")

	os.Args = []string{
		"monako",
		"-only-compose",
		"-fail-on-error",
		"-working-dir", targetDir,
		"-config", monakoConfig,
		"-menu-config", menuConfig}
	t.Logf("Running Monako with %s", os.Args)
	main()

	assert.DirExists(t, filepath.Join(targetDir, "compose"))
	assert.NoDirExists(t, filepath.Join(targetDir, "compose", "public"))

	os.Args = []string{
		"monako",
		"-only-generate",
		"-fail-on-error",
		"-working-dir", targetDir,
		"-config", monakoConfig,
		"-menu-config", menuConfig}
	t.Logf("Running Monako with %s", os.Args)
	main()
	assert.DirExists(t, filepath.Join(targetDir, "compose", "public"))
}

func TestFailOnNoComposeBeforeGenerate(t *testing.T) {

	t.Skip("Does not work since fatal can not be intercepted")
	targetDir := GetLocalTempDir(t)
	err := os.Chdir(targetDir)
	assert.NoError(t, err)

	monakoConfig, menuConfig := writeConfig("https://github.com/snipem/monako-test.git")
	os.Args = []string{
		"monako",
		"-only-generate",
		"-fail-on-error",
		"-working-dir", targetDir,
		"-config", monakoConfig,
		"-menu-config", menuConfig}
	t.Logf("Running Monako with %s", os.Args)
	main()
	assert.NoDirExists(t, filepath.Join(targetDir, "compose", "public"))
}

func getContentFromURL(ts *httptest.Server, url string) (string, error) {
	// res, err := http.Get(ts.URL)
	res, err := http.Get(ts.URL + url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	contentBytes, err := ioutil.ReadAll(res.Body)
	return string(contentBytes), err
}

func TestGetUrlKeyValuesFromHTML(t *testing.T) {

	t.Run("Test a href", func(t *testing.T) {

		urls, err := getURLKeyValuesFromHTML(`

		<a href="/local/path1">Link 1</a>
		<a href="/local/path2">Link 2</a>

		<a href="http://absolute.url/path3">Link 3</a>

		`, "href", "http://example.com")
		assert.NoError(t, err)

		assert.Equal(t, "http://example.com/local/path1", urls[0].String())
		assert.Equal(t, "http://example.com/local/path2", urls[1].String())
		assert.Equal(t, "http://absolute.url/path3", urls[2].String())
	})

	t.Run("Test img src", func(t *testing.T) {

		urls, err := getURLKeyValuesFromHTML(`

		<img src="/local/path1.jpg">Link 1</img>
		<img src="/local/path2.jpg">Link 2</img>

		<img src="http://absolute.url/path3.jpg">Link 3</img>

		`, "src", "http://example.com")
		assert.NoError(t, err)

		assert.Equal(t, "http://example.com/local/path1.jpg", urls[0].String())
		assert.Equal(t, "http://example.com/local/path2.jpg", urls[1].String())
		assert.Equal(t, "http://absolute.url/path3.jpg", urls[2].String())
	})

}

func TestGetVersion(t *testing.T) {
	// No newline in version string
	assert.NotContains(t, getVersion(), "\n")

	t.Run("Test for standard Version values", func(t *testing.T) {
		assert.Contains(t, getVersion(), "Monako Development Local")
		assert.Contains(t, getVersion(), "https://github.com/snipem/monako")
	})

	t.Run("Test if Version string can be overwritten", func(t *testing.T) {
		version = "v.13.34.67"
		commit = "LONGHASHID"

		assert.Contains(t, getVersion(), "Monako v.13.34.67 LONGHASHID")

	})
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

// writeConfig writes a temporary config for a repo in the local test folders
// of this project and returns the path to the monakoConfig and menuConfig
// Also if the MONAKO_TEST_REPO environment variable is set, it will use this
// environment variable for the repository stored in the Monako config.
func writeConfig(repo string) (string, string) {

	var testRepo string

	if os.Getenv("MONAKO_TEST_REPO") != "" {
		testRepo = os.Getenv("MONAKO_TEST_REPO")
	} else {
		testRepo = repo
	}

	monakoConfig := fmt.Sprintf(`
---
    baseURL : "https://example.com/"
    title : "Local Test Page"
    disableCommitInfo: false
  
    whitelist:
      - ".md"
      - ".adoc"
      - ".jpg"
      - ".jpeg"
      - ".svg"
      - ".gif"
      - ".png"
  
    origins:
  
    # Files have to be commited to appear!
    - src: %s
      branch: master
      docdir: .
      targetdir: docs/test/
      blacklist:
        - "test_doc_asciidoc_include_me.adoc"
`, testRepo)

	menuConfig := fmt.Sprintf(`
---
headless: true
---

- **Test docs**
	- [Markdown]({{<relref "test_doc_markdown.md">}})
<br />
	`)

	pathMonakoConfig := "config.testgenerated.yaml"
	err := ioutil.WriteFile(pathMonakoConfig, []byte(monakoConfig), os.FileMode(0600))
	if err != nil {
		log.Fatal(err)
	}

	pathMenuConfig := "menu.testgenerated.md"
	err = ioutil.WriteFile(pathMenuConfig, []byte(menuConfig), os.FileMode(0600))
	if err != nil {
		log.Fatal(err)
	}

	return pathMonakoConfig, pathMenuConfig
}

func GetLocalTempDir(t *testing.T) (tempdir string) {

	localTmpDir := filepath.Join("../../tmp/testdata/", t.Name())
	err := os.MkdirAll(localTmpDir, os.FileMode(0700))
	assert.NoError(t, err)

	return filet.TmpDir(t, localTmpDir)
}
