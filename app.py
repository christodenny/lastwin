import requests
import difflib
from dateutil.parser import parse
from datetime import *
from flask import Flask, jsonify

teamIds = {}
resp = requests.get('http://www.espn.com/college-football/teams').text
start = 0
try:
    while True:
        start = resp.index('www.espn.com/college-football/team/_/id', start)
        idEnd = resp.index('/', start + 40)
        teamId = resp[start+40:idEnd]
        nameStart = resp.index('>', start) + 1
        nameEnd = resp.index('<', nameStart)
        teamName = resp[nameStart:nameEnd]
        teamIds[teamName] = teamId
        start += 1
except:
    print 'done getting team ids'

app = Flask(__name__)

@app.route('/')
def index():
    return 'hello world'

def getDate(body):
    lastWin = body.rfind('greenfont')
    rowteam = body[:lastWin].rfind('row team')
    dateStart = body[rowteam:].index(', ') + 2 + rowteam
    dateEnd = body[dateStart:].index('<') + dateStart
    return body[dateStart:dateEnd]

@app.route('/<teamname>')
def getTeamStats(teamname):
    teamname = difflib.get_close_matches(teamname, teamIds.keys(), 1, 0.2)[0]
    for year in xrange(2017, 2001, -1):
        resp = requests.get('http://www.espn.com/college-football/team/schedule/_/id/' + teamIds[teamname] + '/year/' + str(year)).text
        if 'greenfont' in resp:
            date = getDate(resp)
            lastWin = parse(date + ' ' + str(year))
            today = datetime.now()
            return str((today - lastWin).days)
            break

if __name__ == "__main__":
    app.run()
