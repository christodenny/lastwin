package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sahilm/fuzzy"
)

func loadConfigs() {
	p := os.Getenv("PORT")
	if p != "" {
		log.Printf("Using port %s\n", p)
		port = p
	} else {
		log.Println("$PORT not provided, using default port 8080")
	}
	t := os.Getenv("TTL")
	if t != "" {
		log.Printf("Using cache ttl of %s seconds\n", t)
		if val, err := strconv.Atoi(t); err != nil {
			ttl = time.Second * time.Duration(val)
		}
	} else {
		log.Println("$TTL not provided, using default cache ttl of 5 seconds")
	}
}

func loadTeams() {
	for k, v := range getCfbTeams() {
		allTeams[k] = v
		allTeamNames = append(allTeamNames, k)
	}
	for k, v := range getNflTeams() {
		allTeams[k] = v
		allTeamNames = append(allTeamNames, k)
	}
}

func getRankedTeams(query string) []string {
	matchedTeamNames := []string{}
	matches := fuzzy.Find(query, allTeamNames)
	for i := range matches {
		matchedTeamNames = append(matchedTeamNames, matches[i].Str)
	}
	return matchedTeamNames
}

func getHTML(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	rawText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return string(rawText[:])
}

func getDate(body string) string {
	if strings.Contains(body, ">Preseason<") {
		regularSeason := body[0:strings.Index(body, ">Preseason<")]
		preSeason := body[strings.Index(body, ">Preseason<"):]
		if strings.Contains(regularSeason, "clr-positive") {
			return getDateFromSeason(regularSeason)
		}
		return getDateFromSeason(preSeason)
	}
	return getDateFromSeason(body)
}

func getDateFromSeason(body string) string {
	lastWin := strings.LastIndex(body, "clr-positive")
	rowIdx := strings.LastIndex(body[:lastWin], "tr")
	dateStart := rowIdx + strings.Index(body[rowIdx:], ", ") + 2
	dateEnd := dateStart + strings.Index(body[dateStart:], "<")
	return body[dateStart:dateEnd]
}
