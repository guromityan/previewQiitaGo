package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/labstack/echo"
)

const (
	template = `
<html lang='ja'>
    <head>
        <meta charset='utf-8'>
        <title> %s </title>
        <style type='text/css'>
            %s
        </style>
    </head>
    <body>
        <div class='markdown-body'>
						%s
        </div>
    </body>
</html>
`
	cssFilePath = "./resources/github.css"
)

func handler() func(c echo.Context) error {
	return func(c echo.Context) error {
		target := c.QueryParam("target")
		if target == "" {
			return c.String(http.StatusBadRequest, "query params 'target' is required.")
		}
		if !strings.HasPrefix(target, "https://qiita.com/") {
			return c.String(http.StatusBadRequest, "this service only supported qiita.com")
		}
		if !strings.HasSuffix(target, ".md") {
			target = target + ".md"
		}

		md, err := getMarkdownFromQiita(target)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed")
		}
		html, err := convertMD2HTML(md)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed")
		}
		return c.HTML(http.StatusOK, html)
	}
}

func getMarkdownFromQiita(target string) ([]byte, error) {
	resp, err := http.Get(target)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func convertMD2HTML(md []byte) (string, error) {
	css, err := ioutil.ReadFile(cssFilePath)
	if err != nil {
		return "", err
	}
	body := markdown.ToHTML(md, nil, nil)
	html := fmt.Sprintf(template, getPageTitle(md), string(css), string(body))
	return html, nil
}

func getPageTitle(md []byte) string {
	bytesReader := bytes.NewReader(md)
	bufReader := bufio.NewReader(bytesReader)
	bufReader.ReadLine() // skip first line
	v, _, _ := bufReader.ReadLine()
	title := strings.TrimPrefix(string(v), "title: ")
	return title
}

func main() {
	e := echo.New()
	e.GET("/", handler())
	e.Logger.Fatal(e.Start(":8080"))
}
