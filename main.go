package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

func main() {
	var domain []string
	var urls []string

	// 域
	domain = []string{
		"www.whfwdh.com",
	}

	// 爬取的根URL
	urls = []string{
		"http://www.whfwdh.com/book/index/index/videoid/65e904a9a-0872-d9c1-417d-759e97bbf8b.html",
	}

	for _, url := range urls {
		crawl(domain, url)
	}
}

/// 爬取
func crawl(domain []string, url string) {
	// 初始化默认的收集器
	c := colly.NewCollector(
		// 允许访问的域
		colly.AllowedDomains(domain...),
	)

	// 爬取到页面内的链接 a[href] 时回调
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		// 爬取子链接
		//  只会访问允许域下的子链接
		c.Visit(e.Request.AbsoluteURL(link))
	})

	// 爬取到音频文件时回调
	c.OnHTML("audio source[src]", func(e *colly.HTMLElement) {
		src := e.Attr("src")
		// 音频文件 URL
		var fileUrl string
		if strings.Index(src, "http") == 0 {
			fileUrl = src
		} else {
			fileUrl = e.Request.URL.Scheme + "://" + e.Request.URL.Host + src
		}

		fmt.Printf("File URL: %s\n", fileUrl)
		// 下载文件
		download(fileUrl)
	})

	// 请求链接时回调
	c.OnRequest(func(r *colly.Request) {
		//fmt.Println("Visiting: ", r.URL.String())
	})

	// 开始爬取
	err := c.Visit(url)
	if err != nil {
		fmt.Println("error on c.Visit()")
	}
}

// 下载文件
func download(fileUrl string) error {
	// 请求数据
	resp, err := http.Get(fileUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 文件名
	fileName := getFileName(fileUrl)
	fmt.Printf("File Name: %s\n", fileName)

	// 创建文件用于保存
	out, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	// 将获取的响应流写入到文件流
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// 截取文件名
func getFileName(url string) string {
	separator := "/"
	pathIndex := strings.LastIndex(url, separator)

	var name string
	if pathIndex == -1 {
		name = url
	} else {
		path := strings.Split(url, separator)
		name = path[len(path)-1]
	}

	return name
}
