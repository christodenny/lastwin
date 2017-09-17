import requests
import difflib
from dateutil.parser import parse
from datetime import *
from flask import Flask, jsonify, render_template

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
        teamIds[teamName.lower()] = (teamId, teamName)
        start += 1
except:
    pass

app = Flask(__name__)

@app.route('/')
def index():
    return render_template('index.html')

def getDate(body):
    lastWin = body.rfind('greenfont')
    rowteam = body[:lastWin].rfind('row team')
    dateStart = body[rowteam:].index(', ') + 2 + rowteam
    dateEnd = body[dateStart:].index('<') + dateStart
    return body[dateStart:dateEnd]

@app.route('/<teamname>')
def getTeamStats(teamname):
    teamname = difflib.get_close_matches(teamname.lower(), teamIds.keys(), 1, 0.2)[0]
    team_id, teamname = teamIds[teamname]
    for year in range(2017, 2001, -1):
        resp = requests.get('http://www.espn.com/college-football/team/schedule/_/id/' + team_id + '/year/' + str(year)).text
        if 'greenfont' in resp:
            date = getDate(resp)
            lastWin = parse(date + ' ' + str(year))
            today = datetime.now()
            return render_template('results.html', count=(today - lastWin).days, school=teamname, school_id=team_id)

if __name__ == "__main__":
    app.run()
