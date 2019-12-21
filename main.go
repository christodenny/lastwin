package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	renstromFuzzy "github.com/renstrom/fuzzysearch/fuzzy"
	sahilmFuzzy "github.com/sahilm/fuzzy"
)

var (
	port         = "8080"
	fuzzy        = "renstrom"
	allTeams     = map[string]Team{}
	allTeamNames = []string{}

	cfbQueryURL = "http://www.espn.com/college-football/team/schedule/_/id/%s/season/%s"
	cfbEspnLink = "http://www.espn.com/college-football/team/_/id/%s"
	cfbImgLink  = "http://a.espncdn.com/combiner/i?img=/i/teamlogos/ncaa/500/%s.png&h=200&w=200"
)

func loadConfigs() {
	p := os.Getenv("PORT")
	if p != "" {
		log.Printf("Using port %s\n", p)
		port = p
	} else {
		log.Println("$PORT not provided, using default port 8080")
	}
	f := os.Getenv("FUZZY")
	if f == "renstrom" || f == "sahilm" {
		log.Printf("Using fuzzy search by %s\n", f)
		fuzzy = f
	} else {
		log.Println("$FUZZY not provided, using default fuzzy search by renstrom")
	}
}

func loadTeams() {
	for k, v := range getCfbTeams() {
		allTeams[k] = v
		allTeamNames = append(allTeamNames, k)
	}
}

func getRankedTeams(query string) []string {
	matchedTeamNames := []string{}
	if fuzzy == "renstrom" {
		matches := renstromFuzzy.RankFind(query, allTeamNames)
		sort.Sort(matches)
		for i := range matches {
			matchedTeamNames = append(matchedTeamNames, matches[i].Target)
		}
	} else {
		matches := sahilmFuzzy.Find(query, allTeamNames)
		for i := range matches {
			matchedTeamNames = append(matchedTeamNames, matches[i].Str)
		}
	}
	return matchedTeamNames
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello there")
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
	// espnLink := fmt.Sprintf(cfbEspnLink, team.ID)
	// imgLink := fmt.Sprintf(cfbImgLink, team.ID)
	for year := time.Now().Year(); year >= 2001; year-- {
		queryURL := fmt.Sprintf(cfbQueryURL, team.ID, strconv.Itoa(year))
		body := getHTML(queryURL)
		if strings.Contains(body, "clr-positive") {
			date := getDate(body)
			t, err := time.Parse("Jan 2 2006", date+" "+strconv.Itoa(year))
			if err != nil {
				fmt.Fprintf(w, "Bad date parse: %s", err)
				return
			}
			daysSinceWin := int(time.Since(t).Hours()) / 24
			fmt.Fprintf(w, "%s last won %d days ago", teamname, daysSinceWin)
			return
		}
	}
	fmt.Fprintf(w, "Your team sucks lol")
}

func main() {
	loadTeams()
	r := mux.NewRouter()

	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/{teamname}", lastWinHandler)

	log.Printf("Starting server on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
