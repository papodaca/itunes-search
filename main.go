package main

import (
	"embed"
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
)

//go:embed views
var views embed.FS

//go:embed static
var static embed.FS

func loadTemplate(config goview.Config, tplFile string) (string, error) {
	file, err := views.ReadFile("views/" + tplFile + config.Extension)
	if err != nil {
		return "", err
	}
	return string(file), nil
}

const search_url = "https://itunes.apple.com/search"

func get_itunes_search(search_term string) []gin.H {

	u, err := url.Parse(search_url)
	if err != nil {
		panic(err)
	}

	q, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		panic(err)
	}
	q.Add("term", search_term)
	q.Add("entity", "podcast")
	u.RawQuery = q.Encode()

	final_url := u.String()
	resp, err := http.Get(final_url)
	if err != nil {
		panic(err)
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var result map[string]any
	json.Unmarshal(body, &result)
	if err != nil {
		panic(err)
	}
	count := int(result["resultCount"].(float64))
	items := result["results"].([]interface{})

	results := make([]gin.H, count)

	for ii := 0; ii < count; ii++ {
		item := items[ii].(map[string]any)
		results[ii] = gin.H{"collectionName": item["collectionName"].(string), "feedUrl": item["feedUrl"].(string)}
	}

	return results
}

func main() {
	router := gin.Default()
	router.SetTrustedProxies(nil)
	viewEngine := goview.New(goview.DefaultConfig)
	viewEngine.SetFileHandler(loadTemplate)
	router.HTMLRender = &ginview.ViewEngine{ViewEngine: viewEngine}
	staticFs, err := fs.Sub(static, "static")
	if err != nil {
		panic(err)
	}
	router.StaticFS("/static", http.FS(staticFs))
	router.GET("/", func(c *gin.Context) {
		query := c.Query("q")
		if len(query) > 0 {
			items := get_itunes_search(query)
			c.HTML(http.StatusOK, "index", gin.H{"query": query, "items": items})
		} else {
			c.HTML(http.StatusOK, "index", nil)
		}

	})
	router.Run()
}
