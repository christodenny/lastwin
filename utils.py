import requests
import difflib
import html

def getTeams():
    teamMap = {}
    try:
        resp = requests.get('http://www.espn.com/college-football/teams', timeout=15).text
        searchString = '/college-football/team/_/id/'
        skipLength = len(searchString)
        teamSet = set()
        cursor = 0
        while True:
            cursor = resp.index(searchString, cursor)
            idEnd = resp.index('/', cursor + skipLength)
            teamId = resp[cursor+skipLength : idEnd]
            # print(teamId)
            if teamId in teamSet:
                # skip over <a> ending angle bracket
                cursor = resp.index('<h2', cursor)
                nameStart = resp.index('>', cursor) + 1
                nameEnd = resp.index('<', nameStart)
                teamName = html.unescape(resp[nameStart:nameEnd])
                # print(teamName)
                # ex: teamMap['texas longhorns'] = (251, 'Texas Longhorns')
                teamMap[teamName.lower()] = (teamId, teamName)
            else:
                teamSet.add(teamId)
            cursor += 1
    except ValueError:
        print('{} teams captured'.format(len(teamMap)))
        return teamMap
    except:
        print('getTeams: unexpected error.')
        raise

# Fetch teams on startup
TeamsByName = getTeams()
fuzzymatch = lambda text, num : difflib.get_close_matches(text.lower(), TeamsByName.keys(), num, 0.2)
