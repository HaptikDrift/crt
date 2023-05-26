/*

The following code is a proof of concept:
Project Title: crt.sh
Goal or Aim:
 * As a Proof Of Concept to run against online application "crt.sh" in order to find extra dns entries of exsisting domains and or sub-domains.
ToDo:
 -

 written by Haptik Drift
 <haptikdrift@gmail.com>
*/

package main

/* All imports needed in the main function */
import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

/*	########################################################################################################	*/
/*
	###
	# Start MAIN Function
	###
*/

/*
basic useage example;
$ go run ./crt.sh.go prada.com
*/
func main() {
	uri := os.Args[1]
	slice1 := CrtGet(uri)
	for _, i := range slice1 {
		fmt.Println(i)
	}
}

/*
	###
	# End MAIN Function
	###
*/
/*	########################################################################################################	*/

/*	########################################################################################################	*/
/*
	###
	# Functions used inside the main loop
	###
*/

/* Remove dupicate IP from a created slice */
func RemoveDuplicatesFromSlice(s []string) []string {
	m := make(map[string]bool)
	for _, item := range s {
		if _, ok := m[item]; ok {
		} else {
			m[item] = true
		}
	}
	var result []string
	for item, _ := range m {
		result = append(result, item)
	}
	return result
}

/* Connect to the online CRT application and attempt to return page with table of additional domains */
func CrtGet(v string) []string {
	url := fmt.Sprintf("https://crt.sh/?q=%s", v)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Host", "crt.sh")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.5672.93 Safari/537.36")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,**;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Add("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	str2 := strings.ToLower(string(bodyBytes))
	row := regexp.MustCompile(`<tr[^>]*>(?:.|\n)*?<\/tr>`)
	cell := regexp.MustCompile(`<td[^>]*>(?:.|\n)*?<\/td>`)
	incell := regexp.MustCompile(`<td[^>]*>((?:.|\n)*?)<\/td>`)
	monster := []string{}
	rows := row.FindAllString(str2, -1)
	for _, r := range rows {
		celll := cell.FindAllString(r, -1)
		row := []string{}
		for _, c := range celll {
			mycell := incell.FindStringSubmatch(c)
			row = append(row, strings.TrimSpace(mycell[1]))
		}
		if len(row) > 5 {
			monster = append(monster, row[4])
			row5 := strings.Split(row[5], "<br>")
			monster = append(monster, row5...)

		} else {
			continue
		}
	}
	newmonster := RemoveDuplicatesFromSlice(monster)
	return newmonster
}
