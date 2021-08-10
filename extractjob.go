package main

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"

	"encoding/csv"
	"os"
	"strconv"
	"strings"
)

type extractedJob struct {
	id       string
	title    string
	location string
	salary   string
	summary  string
}

var baseURL string = "https://kr.indeed.com/%EC%B7%A8%EC%97%85?q=%ED%8C%8C%EC%9D%B4%EC%8D%AC"

func main() {
	var jobs []extractedJob
	totalPages := getPages()

	for i := 0; i < totalPages+1; i++ {
		extractedJobs := getPage(i)
		jobs = append(jobs, extractedJobs...)
	}

	writeJobs(jobs)
}

func getPage(page int) []extractedJob {

	var jobs []extractedJob
	pageURL := baseURL + "&start=" + strconv.Itoa(page*10)

	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchCards := doc.Find(".tapItem")

	searchCards.Each(func(i int, card *goquery.Selection) {
		job := extractJob(card)
		jobs = append(jobs, job)
	})

	return jobs
	// 아래 코드 바뀐거니깐 공부해서 확인하자!!

}

func writeJobs(jobs []extractedJob) {
	file, err := os.Create("jobs.csv")

	checkErr(err)

	w := csv.NewWriter(file)

	defer w.Flush()

	headers := []string{"ID", "Title", "Location", "Salary", "Summary"}

	wErr := w.Write(headers)
	checkErr(wErr)

	for _, job := range jobs {
		jobSlice := []string{"https://kr.indeed.com/viewjob?" + job.id, job.title, job.location, job.salary, job.summary}
		jwErr := w.Write(jobSlice)
		checkErr(jwErr)
	}
}

func getPages() int {

	pages := 0

	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length()

		// fmt.Println(s.Find("a").Length())
	})

	return pages

}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status:", res.StatusCode)
	}
}

func extractJob(card *goquery.Selection) extractedJob {
	link, _ := card.Attr("href")
	link = link[8:]
	title := cleanString(card.Find(".jobTitle>span").Text())
	location := cleanString(card.Find(".companyLocation").Text())
	salary := cleanString(card.Find(".salary-snippet").Text())
	summary := card.Find(".job-snippet").Text()

	return extractedJob{id: link,
		title:    title,
		location: location,
		salary:   salary,
		summary:  summary}
}

func cleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}
