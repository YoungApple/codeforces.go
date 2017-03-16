package main

import (
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
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
	return fmt.Sprintf("%s,accepted=%d,submission=%d,name=%s,tags=%s",
		p.url, p.accepted, p.submission, p.name, strings.Join(p.tags, "|"))
}

const (
	codeforces = "http://codeforces.com"
	problemset = codeforces + "/problemset/page/"
	pageNum    = 1
)

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

	accepted, _ := strconv.Atoi(strings.TrimSpace(p.Children().Eq(3).Find("a").Text())[1:])
	statusLink, _ := p.Children().Eq(3).Find("a").Attr("href")
	submission := fetchSubmissionStatus(codeforces + statusLink)
	return Problem{name, strings.TrimSpace(link), tags, accepted, submission}
}

func codeforcesHashCode(cookie string) int {
	var hashcode = 0
	for i := 0; i < len(cookie); i++ {
		hashcode = (hashcode + (i+1)*(i+2)*(int)(cookie[i])) % 1009
		if i%3 == 0 {
			hashcode++
		}
		if i%2 == 0 {
			hashcode *= 2
		}
		if i > 0 {
			hashcode -= (int)(math.Floor(float64(cookie[int(math.Floor(float64(i)/2))]/
				2))) * (hashcode % 5)
		}
		for hashcode < 0 {
			hashcode += 1009
		}
		for hashcode >= 1009 {
			hashcode -= 1009
		}
	}
	return hashcode
}

func fetchSubmissionStatus(statusUrl string) int {
	response, err := http.Get(statusUrl)
	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		log.Fatal(err)
	}

	var cookie_seed string
	for _, cookie := range response.Cookies() {
		if (*cookie).Name == "39ce7" {
			cookie_seed = (*cookie).Value
		}
	}

	fmt.Println(response.Cookies())
	fmt.Println(statusUrl)
	fmt.Println("## 11 PageSize:" + doc.Find("div.pagination span.page-index").Last().Text())
	csrf_token, _ := doc.Find("form.status-filter input[name='csrf_token']").Attr("value")
	frameProblemIndex, _ := doc.Find("form.status-filter input[name='frameProblemIndex']").Attr("value")
	tta := codeforcesHashCode(cookie_seed)

	data := url.Values{
		"csrf_token":            {csrf_token},
		"action":                {"setupSubmissionFilter"},
		"frameProblemIndex":     {frameProblemIndex},
		"verdictName":           {"anyVerdict"},
		"programTypeForInvoker": {"anyProgramTypeForInvoker"},
		"comparisonType":        {"NOT_USED"},
		"judgedTestCount":       {},
		"_tta":                  {strconv.Itoa(tta)},
	}
	fmt.Println(data)

	response, err = http.PostForm(statusUrl, data)
	fmt.Println("post statusCode=" + response.Status)
	doc, err = goquery.NewDocumentFromResponse(response)
	if err != nil {
		log.Fatal(err)
	}

	//client := http.Client{}
	//request, err := http.NewRequest("POST", statusUrl, nil)
	if err != nil {
		log.Fatal(err)
	}
	//
	////request.Header.Add("Cookie", cookies)
	//
	//response, err = client.Do(request)

	//response, err = http.PostForm(statusUrl, data)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//doc, err = goquery.NewDocumentFromResponse(response)
	//fmt.Println(doc.Html())
	//fmt.Println("\n\n @@@@@@@@@@@@@@@@@")

	response, err = http.Get(statusUrl)

	os.Exit(1)
	return 0
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
		case <-pch:
			//case p := <-pch:
			//fmt.Println(strconv.Itoa(i) + p.String())
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

	for i := 1; i <= pageNum; i++ {
		go fetchUrl(problemset+strconv.Itoa(i), pch, finished)
	}

	write(pch, finished)
}
