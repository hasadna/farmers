var PE = "PenmanEvaporation";
var RAIN = "RAIN";
var URL = "http://http://www.israelmeteo.mobi/Ajax/getStations";
var meteoInfo = {};

function getMeteoInfo(stationName, city) {
    var xmlhttp = new XMLHttpRequest();    
    xmlhttp.onreadystatechange = function() {
        if (this.readyState == 4 && this.status == 200) {
            var meteoInfoArr = JSON.parse(this.responseText);
            //createMeteoInfo(meteoInfoArr, stationName, city);
        }
    };
    xmlhttp.open("GET", url, false);
    xmlhttp.send();
}

function createMeteoInfo(meteoInfoArr, stationName, city) {
    for(i = 0; i < meteoInfoArr.length; i++) {
        var rain = findRain(meteoInfoArr[i]);
        var pe = findPE(meteoInfoArr[i]);
        meteoInfo[meteoInfoArr[i].name] = {
            rain: rain,
            pe: pe
        }
    }
}

function getStationNames() {
    return ["נווה שאנן", "ק. אתא", "אשדוד"];
}