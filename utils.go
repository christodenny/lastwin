package main

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/sahilm/fuzzy"
)

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
