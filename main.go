package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	rice "github.com/GeertJohan/go.rice"
	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
)

//go:generate rice embed-go

func BinDataRenderer() *ginview.ViewEngine {
	viewEngine := goview.New(goview.DefaultConfig)
	viewEngine.SetFileHandler(func(config goview.Config, tplFile string) (content string, err error) {
		box := rice.MustFindBox("views")
		file, err := box.String(tplFile + config.Extension)
		if err != nil {
			return "", err
		}
		return string(file), nil
	})
	return &ginview.ViewEngine{
		ViewEngine: viewEngine}
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
	router.HTMLRender = BinDataRenderer()
	router.GET("/", func(c *gin.Context) {
		query := c.Query("query")
		if len(query) > 0 {
			items := get_itunes_search(query)
			c.HTML(http.StatusOK, "index", gin.H{"query": query, "items": items})
		} else {
			c.HTML(http.StatusOK, "index", gin.H{})
		}

	})
	router.Run()
}
