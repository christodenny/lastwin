import requests
from dateutil.parser import parse
from datetime import *
from flask import Flask, jsonify, render_template, request
from utils import TeamsByName, fuzzymatch

app = Flask(__name__)

@app.route('/')
def index():
    return render_template('index.html')

def getDateFromSeason(body):
    lastWin = body.rfind('clr-positive')
    rowIdx = body[:lastWin].rfind('tr')
    dateStart = rowIdx + body[rowIdx:].index(', ') + 2
    dateEnd = dateStart + body[dateStart:].index('<')
    return body[dateStart:dateEnd]

def getDate(body):
    if '>Preseason<' in body:
        regularSeason = body[0:body.index('>Preseason<')]
        preSeason = body[body.index('>Preseason<'):]
        if 'clr-positive' in regularSeason:
            return getDateFromSeason(regularSeason)
        else:
            return getDateFromSeason(preSeason)
    else:
        return getDateFromSeason(body)

@app.route('/<teamname>')
def getTeamStats(teamname):
    teamname = fuzzymatch(teamname.lower(), 1)[0]
    team_id, teamname, division = TeamsByName[teamname]
    for year in range(datetime.now().year, 2001, -1):
        if division == "cfb":
            baseUrl = 'http://www.espn.com/college-football/team/schedule/_/id/'
            espnLink = 'http://www.espn.com/college-football/team/_/id/' + team_id
            imgLink = 'http://a.espncdn.com/combiner/i?img=/i/teamlogos/ncaa/500/' + team_id + '.png&h=200&w=200'
        else:
            baseUrl = 'https://www.espn.com/nfl/team/schedule/_/name/'
            espnLink = 'http://www.espn.com/nfl/team/_/id/' + team_id
            imgLink = 'http://a.espncdn.com/combiner/i?img=/i/teamlogos/nfl/500/' + team_id + '.png&h=200&w=200'

        resp = requests.get(baseUrl + team_id + '/season/' + str(year), timeout=15).text
        if 'clr-positive' in resp:
            date = getDate(resp)
            lastWin = parse(date + ' ' + str(year))
            today = datetime.now()
            return render_template('results.html', count=(today - lastWin).days, school=teamname, school_id=team_id, espnLink=espnLink, imgLink=imgLink)

@app.route('/autocomplete')
def autocomplete():
    text = request.args.get('text')
    data = [{"id": TeamsByName[name][0], "name": TeamsByName[name][1]} for name in fuzzymatch(text, 3)]
    return jsonify(data)

if __name__ == "__main__":
    app.run()
