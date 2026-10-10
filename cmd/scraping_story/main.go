package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/microcosm-cc/bluemonday"
)

func main() {
	url := "https://ncode.syosetu.com/n9619hu/1"

	// User-Agentを設定
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// タイトル取得
	title := doc.Find(".p-novel__title").Text()

	// 著者取得（c-announceの中の2番目のaタグ）
	author := doc.Find(".c-announce a").Eq(2).Text()

	// 本文取得
	var content []string
	doc.Find(".p-novel__text p").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		if text != "" {
			content = append(content, text)
		}
	})

	// bluemondayでサニタイズ
	p := bluemonday.NewPolicy()
	p.AllowElements("p")
	cleanedContent := p.Sanitize(strings.Join(content, "\n"))

	// 出力
	fmt.Printf("タイトル: %s\n", strings.TrimSpace(title))
	fmt.Printf("著者: %s\n", strings.TrimSpace(author))
	fmt.Printf("\n本文:\n%s\n", cleanedContent)
}
