let URL = "/meteodata.json";
let meteoInfo = {};
let stationNames = [];

function createMeteoInfo(meteoInfoArr) {
    for(i = 0; i < meteoInfoArr.length; i++) {
        var name = meteoInfoArr[i].name ? meteoInfoArr[i].name : meteoInfoArr[i].city;
        if (!name || ! meteoInfoArr[i].monitors) {
            continue;
        }
        meteoInfo[name] = {};
        var pe = findSensor(meteoInfoArr[i].monitors, "PenmanEvaporation");
        // Populate meteoInfo dict.
        if (! pe) {
            continue;
        }
        meteoInfo[name]["pe"] = [pe];
        // Populate names array.
        stationNames.push(name);
    }
}

function loadMeteoInfo(callback) {
    var xmlhttp = new XMLHttpRequest();    
    xmlhttp.onreadystatechange = function() {
        if (this.readyState == 4 && this.status == 200) {
            var meteoInfoArr = JSON.parse(this.responseText);
            createMeteoInfo(meteoInfoArr);
            callback();
        }
    };
    xmlhttp.open("GET", URL, true);
    xmlhttp.send();
}

function getStationNames() {
    return stationNames;
}

function getInfoForStation(stationName) {
	return meteoInfo[stationName];
}

function findSensor(monitors, sensorName) {
    for(j = 0; j < monitors.length; j++) {
        if (monitors[j].sensor == sensorName) {
            return monitors[j].value;
        }
    }
    return null;
}
