package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

var (
	port         = "8080"
	ttl          = time.Second * 5 // Cache ttl in seconds
	allTeams     = map[string]Team{}
	allTeamNames = []string{}
	cache        = map[string]CacheEntry{}
	templates    = map[string]*template.Template{}

	urls = map[string]map[string]string{
		"cfb": {
			"query": "http://www.espn.com/college-football/team/schedule/_/id/%s/season/%s",
			"espn":  "http://www.espn.com/college-football/team/_/id/%s",
			"img":   "http://a.espncdn.com/combiner/i?img=/i/teamlogos/ncaa/500/%s.png",
		},
		"nfl": {
			"query": "https://www.espn.com/nfl/team/schedule/_/name/%s/season/%s",
			"espn":  "http://www.espn.com/nfl/team/_/id/%s",
			"img":   "http://a.espncdn.com/combiner/i?img=/i/teamlogos/nfl/500/%s.png",
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
	URL      string
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	templates["home"].ExecuteTemplate(w, "index.html", ResultData{URL: "lastwin.info"})
}

func lastWinHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamname := vars["teamname"]
	matchedTeamNames := getRankedTeams(teamname)
	if len(matchedTeamNames) == 0 {
		fmt.Fprintf(w, "Could not find a team")
		return
	}
	teamname = matchedTeamNames[0]
	team := allTeams[teamname]
	espnLink := fmt.Sprintf(urls[team.Division]["espn"], team.ID)
	imgLink := fmt.Sprintf(urls[team.Division]["img"], team.ID)
	// TODO: do smth about thundering herd with this cache
	if entry, ok := cache[teamname]; ok && time.Since(entry.lastRefresh) < ttl {
		resultData := ResultData{Count: entry.lastWin, School: teamname, SchoolID: team.ID, EspnLink: espnLink, ImgLink: imgLink, URL: "lastwin.info" + r.URL.Path + "?" + r.URL.RawQuery}
		log.Printf("Handling %s from cache\n", teamname)
		templates["results"].ExecuteTemplate(w, "results.html", resultData)
		return
	}
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
			// The server's local time may be in a different time zone than the
			// user's time, so convert both of them to UTC.
			userTime := time.Now().UTC()
			log.Printf("Current time in UTC: %s", userTime.Format(time.UnixDate))
			tzs, ok := r.URL.Query()["tz"]
			if ok && len(tzs[0]) > 0 {
				tz, err := strconv.Atoi(tzs[0])
				if err == nil {
					// Get the user's time as if it were in UTC, i.e. 10:30pm
					// PST would become 10:30pm UTC since the ESPN date is
					// interpreted as UTC
					userTime = userTime.Add(-time.Minute * time.Duration(tz))
					log.Printf("Current time in UTC after translation: %s", userTime.Format(time.UnixDate))
				}
			}
			daysSinceWin := int(userTime.Sub(t).Hours()) / 24
			if daysSinceWin < 0 {
				// Don't allow negative days (which may be possible if the
				// user's time zone is ahead of the location where the game
				// occurred)
				daysSinceWin = 0
			}
			resultData := ResultData{Count: daysSinceWin, School: teamname, SchoolID: team.ID, EspnLink: espnLink, ImgLink: imgLink, URL: "lastwin.info" + r.URL.Path + "?" + r.URL.RawQuery}
			cache[teamname] = CacheEntry{time.Now(), daysSinceWin}
			log.Printf("Handling %s from espn\n", teamname)
			templates["results"].ExecuteTemplate(w, "results.html", resultData)
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

func Gzip(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Force serving the gzip version for bootstrap.min.css since it's so big.
		// Everything else is small enough it doesn't matter.
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") || r.URL.Path != "bootstrap.min.css" {
			handler.ServeHTTP(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Content-Type", "text/css")
		r.URL.Path = r.URL.Path + ".gz"
		handler.ServeHTTP(w, r)
	})
}

func main() {
	t, err := template.ParseFiles("tmpl/base.html", "tmpl/index.html")
	if err != nil {
		log.Fatal("Error parsing template files")
	}
	templates["home"] = t
	t, err = template.ParseFiles("tmpl/base.html", "tmpl/results.html")
	if err != nil {
		log.Fatal("Error parsing template files")
	}
	templates["results"] = t
	loadTeams()
	loadConfigs()
	r := mux.NewRouter()

	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/autocomplete", autocompleteHandler)
	r.HandleFunc("/{teamname}", lastWinHandler)

	fileServer := Gzip(http.FileServer(http.Dir("static/")))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
