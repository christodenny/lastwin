package main

import (
	"fmt"
	"html"
	"strings"
)

var (
	cfbTeamsURL = "http://www.espn.com/college-football/teams"
	nflTeamsURL = "http://www.espn.com/nfl/teams"
)

// Team is a struct storing team name and id
type Team struct {
	ID       string
	Name     string
	Division string
}

func getCfbTeams() map[string]Team {
	body := getHTML(cfbTeamsURL)
	searchString := "/college-football/team/_/id/"
	skipLength := len(searchString)
	teamSet := map[string]bool{}
	teamMap := map[string]Team{}
	cursor := 0
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
			teamMap[strings.ToLower(teamName)] = Team{teamID, teamName, "cfb"}
		} else {
			cursor++
		}
	}
	fmt.Printf("%d CFB teams captured\n", len(teamMap))
	return teamMap
}

func getNflTeams() map[string]Team {
	body := getHTML(nflTeamsURL)
	searchString := "/nfl/team/_/name/"
	skipLength := len(searchString)
	teamSet := map[string]bool{}
	teamMap := map[string]Team{}
	cursor := 0
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
			teamMap[strings.ToLower(teamName)] = Team{teamID, teamName, "nfl"}
		} else {
			cursor++
		}
	}
	fmt.Printf("%d NFL teams captured\n", len(teamMap))
	return teamMap
}
