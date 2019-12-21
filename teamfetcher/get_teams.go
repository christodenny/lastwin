package teamfetcher

import (
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

var (
	collegeTeamsURL = "http://www.espn.com/college-football/teams"
)

// GetCfbTeams returns all college football teams.
func GetCfbTeams() string {
	resp, err := http.Get(collegeTeamsURL)
	if err != nil {
		return "failed to get http"
	}
	defer resp.Body.Close()
	rawText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "failed to read http response body"
	}
	body := string(rawText[:])
	searchString := "/college-football/team/_/id/"
	skipLength := len(searchString)
	teamSet := map[string]bool{}
	cursor := 0
	allTeams := ""
	for cursor < len(body) {
		index := strings.Index(body[cursor:], searchString)
		if index == -1 {
			break
		}
		cursor += index
		index = strings.Index(body[cursor+skipLength:], "/")
		if index == -1 {
			break
		}
		idEnd := cursor + skipLength + index
		teamID := body[cursor+skipLength : idEnd]
		if _, ok := teamSet[teamID]; !ok {
			// skip over <a> ending angle bracket
			cursor = strings.Index(body[cursor:], "<h2") + cursor
			nameStart := strings.Index(body[cursor:], ">") + cursor + 1
			nameEnd := strings.Index(body[nameStart:], "<") + nameStart
			teamName := html.UnescapeString(body[nameStart:nameEnd])
			cursor = nameEnd
			teamSet[teamID] = true
			allTeams += teamName + "\n"
		} else {
			cursor++
		}
	}
	fmt.Println(strconv.Itoa(len(teamSet)))
	return allTeams
}
