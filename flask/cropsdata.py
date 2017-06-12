'''
Created on Jun 10, 2017

@author: ohad
'''
from os import listdir, path
import pandas
from logging import getLogger
from datetime import datetime


# import sys
# reload(sys)
# sys.setdefaultencoding('UTF-8')

CROPSDATA_DIR = 'cropsdata'

crops_types = []
crop_coeffiecients = {}

def loadCrops():
    global crops_types, crop_coeffiecients
    logger = getLogger("__main__")
    logger.debug("reloading crops data")
    now = datetime.now()
    for filename in listdir(CROPSDATA_DIR):
        if filename.endswith(".csv"):
            data = pandas.read_csv(path.join(CROPSDATA_DIR, filename))
            months = data['month'].asobject
            start_days = data['start_day'].asobject
            end_days = data['end_day'].asobject
            
            is_int = lambda d:isinstance(d, int) or isinstance(d, long) or (isinstance(d, float) and d.is_integer())
            is_int_arr = lambda a: all(map(is_int, a))
            if not (is_int_arr(months) and is_int_arr(start_days) and is_int_arr(end_days)):
                logger.error('failed to read '+ filename)
                continue
            
            today = None
            for line in range(len(months)):
                if months[line] == now.month and start_days[line] <= now.day and end_days[line] <= now.day:
                    today = line
                    continue
            if not today:
                logger.error('no data for today' + filename)
                continue
            crops = [d for d in data.columns if d and d not in [u'month', u'start_day', u'end_day']]
            crops_types += crops
            for crop in crops:
                crop_coeffiecients[crop] = data[crop][today]

loadCrops()
if __name__ == '__main__':
    print(crop_coeffiecients)