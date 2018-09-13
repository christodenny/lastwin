import requests
from dateutil.parser import parse
from datetime import *
from flask import Flask, jsonify, render_template, request
from utils import TeamsByName, fuzzymatch

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
    teamname = fuzzymatch(teamname.lower(), 1)[0]
    team_id, teamname = TeamsByName[teamname]
    for year in range(datetime.now().year, 2001, -1):
        resp = requests.get('http://www.espn.com/college-football/team/schedule/_/id/' + team_id + '/year/' + str(year), timeout=15).text
        if 'greenfont' in resp:
            date = getDate(resp)
            lastWin = parse(date + ' ' + str(year))
            today = datetime.now()
            return render_template('results.html', count=(today - lastWin).days, school=teamname, school_id=team_id)

@app.route('/autocomplete')
def autocomplete():
    text = request.args.get('text')
    data = [{"id": TeamsByName[name][0], "name": TeamsByName[name][1]} for name in fuzzymatch(text, 3)]
    return jsonify(data)

if __name__ == "__main__":
    app.run()
