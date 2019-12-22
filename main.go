package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

var (
	port         = "8080"
	allTeams     = map[string]Team{}
	allTeamNames = []string{}

	urls = map[string]map[string]string{
		"cfb": {
			"query": "http://www.espn.com/college-football/team/schedule/_/id/%s/season/%s",
			"espn":  "http://www.espn.com/college-football/team/_/id/%s",
			"img":   "http://a.espncdn.com/combiner/i?img=/i/teamlogos/ncaa/500/%s.png&h=200&w=200",
		},
		"nfl": {
			"query": "https://www.espn.com/nfl/team/schedule/_/name/%s/season/%s",
			"espn":  "http://www.espn.com/nfl/team/_/id/%s",
			"img":   "http://a.espncdn.com/combiner/i?img=/i/teamlogos/nfl/500/%s.png&h=200&w=200",
		},
	}
)

// ResultData is struct used for html template
type ResultData struct {
	Count    int
	School   string
	SchoolID string
	EspnLink string
	ImgLink  string
}

func loadConfigs() {
	p := os.Getenv("PORT")
	if p != "" {
		log.Printf("Using port %s\n", p)
		port = p
	} else {
		log.Println("$PORT not provided, using default port 8080")
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

func homeHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("tmpl/base.html", "tmpl/index.html")
	if err != nil {
		fmt.Fprintf(w, "Error parsing template file")
		return
	}
	t.ExecuteTemplate(w, "index.html", nil)
}

func lastWinHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamname := strings.ToLower(vars["teamname"])
	matchedTeamNames := getRankedTeams(teamname)
	if len(matchedTeamNames) == 0 {
		fmt.Fprintf(w, "Could not find a team")
		return
	}
	teamname = matchedTeamNames[0]
	team := allTeams[teamname]
	espnLink := fmt.Sprintf(urls[team.Division]["espn"], team.ID)
	imgLink := fmt.Sprintf(urls[team.Division]["img"], team.ID)
	for year := time.Now().Year(); year >= 2001; year-- {
		queryURL := fmt.Sprintf(urls[team.Division]["query"], team.ID, strconv.Itoa(year))
		body := getHTML(queryURL)
		if strings.Contains(body, "clr-positive") {
			date := getDate(body)
			date = date + " " + strconv.Itoa(year)
			t, err := time.Parse("Jan 2 2006", date)
			if err != nil {
				fmt.Fprintf(w, "Error parsing date %s: %s", date, err)
				return
			}
			now := time.Now()
			tzs, ok := r.URL.Query()["tz"]
			if !ok || len(tzs[0]) < 1 {
				tz, err := strconv.Atoi(tzs[0])
				if err != nil {
					now = now.Add(-time.Minute * time.Duration(tz))
				}
			}
			daysSinceWin := int(now.Sub(t).Hours()) / 24
			resultData := ResultData{Count: daysSinceWin, School: teamname, SchoolID: team.ID, EspnLink: espnLink, ImgLink: imgLink}
			tmpl, err := template.ParseFiles("tmpl/base.html", "tmpl/results.html")
			if err != nil {
				fmt.Fprintf(w, "Error parsing template file")
				return
			}
			tmpl.ExecuteTemplate(w, "results.html", resultData)
			return
		}
	}
	fmt.Fprintf(w, "Your team sucks lol")
}

func autocompleteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	queries, ok := r.URL.Query()["text"]
	if !ok || len(queries[0]) < 1 {
		json.NewEncoder(w).Encode([]string{})
	}
	query := queries[0]
	matchedTeamNames := getRankedTeams(query)
	json.NewEncoder(w).Encode(matchedTeamNames)
}

func main() {
	loadTeams()
	loadConfigs()
	r := mux.NewRouter()

	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/autocomplete", autocompleteHandler)
	r.HandleFunc("/{teamname}", lastWinHandler)

	fileServer := http.FileServer(http.Dir("static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	log.Printf("Starting server on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
