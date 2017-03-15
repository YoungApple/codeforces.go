package main

import (
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"os"
	"strings"
)

type Problem struct {
	name       string
	url        string
	tags       []string
	accepted   int
	submission int
}

func (p *Problem) String() string {
	return p.url + ",tags=" + strings.Join(p.tags, "|") + ",name=" + p.name
}

var baseUrl = "http://codeforces.com/problemset/page/"
var pageNum = 34

func parseProblem(p *goquery.Selection) Problem {
	// name and link
	link, _ := p.Children().Eq(1).Find("a").Attr("href")

	var name string
	var tags []string
	p.Children().Eq(1).Find("a").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			name = strings.TrimSpace(s.Text())
		} else {
			tags = append(tags, strings.TrimSpace(s.Text()))
		}
	})

	// TODO(yangshuguo): AC ratio
	return Problem{name, strings.TrimSpace(link), tags, 0, 0}
}

func fetchUrl(url string, pch chan Problem, finished chan int) {
	fmt.Println("Fetching url", url)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("table.problems").First().Children().Children().Each(func(i int, s *goquery.Selection) {
		if i != 0 {
			pch <- parseProblem(s)
		}
	})
	finished <- 1
}

func write(pch chan Problem, finished chan int) {
	var i = 0
	var count = pageNum
	for {
		select {
		case p := <-pch:
			fmt.Println(fmt.Sprintf("%d, %s", i, p.String()))
			i++
		case <-finished:
			count--
			if count == 0 {
				return
			}
		}
	}
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) > 5 {
		fmt.Println("Error! Too many arguments!")
		os.Exit(1)
	}

	var pch = make(chan Problem)
	defer close(pch)

	var finished = make(chan int)
	defer close(finished)

	for i := 0; i < pageNum; i++ {
		go fetchUrl(fmt.Sprintf("%s%d", baseUrl, i), pch, finished)
	}

	write(pch, finished)
}
