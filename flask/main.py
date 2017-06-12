'''
Created on Jun 10, 2017

@author: ohad
'''
from flask import Flask, Response
from apscheduler.schedulers.background import BackgroundScheduler
from apscheduler.triggers.interval import IntervalTrigger

import sys
reload(sys)
sys.setdefaultencoding('UTF-8')

from config import port
from meteodata import meteodatajson
import meteodata
from cropsdata import loadCrops

app = Flask(__name__, static_url_path='')

@app.route('/', methods=['GET'])
def index():
    return app.send_static_file('index.html')

@app.route('/meteodata.json', methods=['GET'])
def dumpmeteodatajson():
    return Response(meteodatajson(), mimetype='application/json') 

@app.route('/crop/<crop>')
def crop_complete():
    pass

@app.route('/compute/<crop>/<stage>/<float:latitude>/<float:longitude>')
@app.route('/compute/<crop>/<stage>/<float:latitude>/<float:longitude>/<int:days>')
def compute(plant, stage, latitude, longitude, days = 1):
    return str(12 * days) + " mm / d"

@app.route('/computestation/<crop>/<stage>/<station>/<int:days>')
def computeByStation(plant, stage, station, days = 1):
    return str(12 * days) + " mm / d"

if __name__ == "__main__":
    scheduler = BackgroundScheduler()
    scheduler.start()
    scheduler.add_job(
        func=meteodata.loadMeteorologicalData,
        trigger=IntervalTrigger(minutes=10),
        id='meteodata.reloadData.minuts10',
        name='reload meteorological data every few minutes',
        replace_existing=True)
    scheduler.add_job(
        func=loadCrops,
        trigger=IntervalTrigger(days=1),
        id='loadCrops.day',
        name='reload crops data for today',
        replace_existing=True)
    app.run(port=port)