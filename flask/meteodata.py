'''
Created on Jun 10, 2017

@author: ohad
'''

from urllib2 import urlopen
from time import time
from json import loads, dumps
from traceback import print_exc
from schedule import every
from logging import getLogger

# import sys
# reload(sys)
# sys.setdefaultencoding('UTF-8')

meteodata = []
stations = []
raw = []

def loadMeteorologicalData():
    global raw, meteodata, stations
    logger = getLogger("__main__")
    try:
        logger.debug("reloading stations data")
        rawtext = urlopen('http://www.israelmeteo.mobi/Ajax/getStations').read()
        open('raw' + str(time()) + '.txt', "wt").write(rawtext)
        stations = loads(rawtext)
        temp_meteodata = []
        for station in stations:
#             print station["name"]
            if not isinstance(station, dict) or not station.has_key("monitors"
                    ) or not isinstance(station["monitors"], list):
                logger.error("no monitors for station")
                continue
            for sensor in station["monitors"]:
                if not isinstance(sensor, dict) or not sensor.has_key("sensor"):
                    logger.error("no sensor for monitor")
                    continue
                if sensor["sensor"] == "PenmanEvaporation":
                    station["PenmanEvaporation"] = sensor["value"]
                    temp_meteodata.append(station)
        meteodata = temp_meteodata
    except Exception as e:
        print_exc()

def getPenamanEvaporation(latitude, longitude):
    mindist = 999999999
    evaporation = None
    for station in meteodata:
        dist = (station["latitude"] - latitude) **2 + (station["longitude"] - longitude) **2
        if dist < mindist:
            mindist = dist
            evaporation = station["PenmanEvaporation"]
    return evaporation

def meteodatajson():
    return dumps(meteodata)

if __name__ == '__main__':
    loadMeteorologicalData()
    print (meteodatajson())
#     print(len(meteodata))
else:
    loadMeteorologicalData()
    every(10).minutes.do(loadMeteorologicalData)