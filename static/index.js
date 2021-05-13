let map;
let pos;
let markers = [];
let filter = [];
let singapore;
let URL = 'https://localhost/api/v1/'

const getSingaporeLocation = async () => {
    const response = await fetch(URL + 'location/singapore', {
        method: 'GET',
        headers: new Headers({ 'Content-Type': 'application/json' }),
        insecure: true,
    })

    return response.json()
}

async function initMap() {
    //use singapore location for first load
    //if no singapore location returned, use current location
    singapore = await getSingaporeLocation()
        .catch(err => console.log(err));
    if (singapore.lat > 0) {
        map = new google.maps.Map(document.getElementById("map"), {
            zoom: 12,
            disableDefaultUI: true,
            center: singapore,
        });
        pos = singapore;
    } else {
        //use geolocation
        infoWindow = new google.maps.InfoWindow();
        if (navigator.geolocation) {
            navigator.geolocation.getCurrentPosition(
                (position) => {
                    pos = {
                        lat: position.coords.latitude,
                        lng: position.coords.longitude,
                    };
                    map = new google.maps.Map(document.getElementById("map"), {
                        zoom: 14,
                        disableDefaultUI: true,
                        center: pos,
                    });

                    addMarker(map);
                },
                () => {
                    handleLocationError(true, infoWindow, map.getCenter());
                }
            );
        } else {
            // Browser doesn't support Geolocation
            handleLocationError(false, infoWindow, map.getCenter());
        }

    };
    addMarker(map);

    document.getElementById("submit").addEventListener("click", () => {
        searchAddress(map);
    });

    document.getElementById("providerList").addEventListener("change", () => {
        getProviderLocation(map);
    });

}

function handleLocationError(browserHasGeolocation, infoWindow, pos) {
    infoWindow.setPosition(pos);
    infoWindow.setContent(
        browserHasGeolocation
            ? "Error: The Geolocation service failed."
            : "Error: Your browser doesn't support geolocation."
    );
    infoWindow.open(map);
}


const getEVLocations = async () => {
    const response = await fetch(URL + 'location/points', {
        method: 'GET',
        cache: 'no-cache',
        headers: new Headers({ 'Content-Type': 'application/json' }),
        redirect: 'follow',
        insecure: true,
        credentials: "same-origin",
    });
    return response.json()
}

var InfoObj = [];

async function addMarker(map) {
    const evLocation = await getEVLocations()
        .catch(err => console.log(err));

    let url = "http://maps.google.com/mapfiles/ms/icons/";

    deleteMarkers();
    deleteFilter();
    for (var i = 0; i < evLocation.length; i++) {
        var contentString = '<h6>' + evLocation[i].provider + '</h6>' +
            '<h6>' + evLocation[i].address + '</h6>';
        var color = "green";
        for (var j = 0; j < evLocation[i].charger.length; j++) {
            contentString += '<h6>' + evLocation[i].charger[j].type + ' - ' + evLocation[i].charger[j].station + ' station </h6>';
            if (evLocation[i].charger[j].match === true) {
                color = "blue";
            };
        };

        const marker = new google.maps.Marker({
            position: { lat: evLocation[i].location[0], lng: evLocation[i].location[1] },
            map: map,
            icon: { url: url + color + ".png" }
        });

        markers.push(marker);

        setMapOnAll(map)
        const infowindow = new google.maps.InfoWindow({
            content: contentString,
        });
        marker.addListener("mouseover", () => {
            infowindow.open(marker.get('map'), marker);
            InfoObj[0] = infowindow;
        });
        marker.addListener("mouseout", () => {
            closeOtherInfo();
        });

        marker.addListener("click", () => {
            const position = marker.getPosition();
            locationDetailInfo(position.lat(), position.lng());
        });
    };

    getFilter();
}

function closeOtherInfo() {
    if (InfoObj.length > 0) {
        InfoObj[0].set("marker", null);
        InfoObj[0].close();
        InfoObj[0].length = 0;
    }
}

function home() {
    initMap()
}


//get location info
const getEVLocationInfo = async (lat, lng) => {
    const response = await fetch(URL + 'location/info?lat=' + lat + '&lng=' + lng, {
        method: 'GET',
        cache: 'no-cache',
        headers: { 'Content-Type': 'application/json' },
        redirect: 'follow',
        insecure: true,
    });

    return response.json()
}
async function locationDetailInfo(lat, lng) {
    const evLocationInfo = await getEVLocationInfo(lat, lng);
    displayInfo(evLocationInfo);
}

function displayInfo(evLocationInfo) {
    document.getElementById('displayInfo').innerHTML = "";
    document.getElementById('displayInfo').innerHTML += "<p><b>Information:</b></p>";
    document.getElementById('displayInfo').innerHTML += "<h6>" + evLocationInfo.provider + "<h6>";
    document.getElementById('displayInfo').innerHTML += "<h6>" + evLocationInfo.address + "<h6>";
    if (evLocationInfo.operator) {
        document.getElementById('displayInfo').innerHTML += "<h6>Operated by: " + evLocationInfo.operator + "<h6>";
    };

    if (evLocationInfo.operator) {
        document.getElementById('displayInfo').innerHTML += "<h6>" + evLocationInfo.requirement + "<h6>";
    };

    for (i = 0; i < evLocationInfo.charger.length; i++) {
        document.getElementById('displayInfo').innerHTML += "<h6><b><u>Charger Type: </u></b></h6>";
        if (evLocationInfo.charger[i].match === true) {
            var color = "blue";
        } else {
            var color = "none";
        };

        document.getElementById('displayInfo').innerHTML += "<h6 style=color:" + color + ">" +
            evLocationInfo.charger[i].type + " - " + evLocationInfo.charger[i].detail + "</h6>";
        document.getElementById('displayInfo').innerHTML += "<h6 style=color:" + color + ">Price:" +
            evLocationInfo.charger[i].price + "</h6>";
        document.getElementById('displayInfo').innerHTML += "<h6 style=color:" + color + ">Number of station: " +
            evLocationInfo.charger[i].station + "</h6>";
    }

    document.getElementById('displayInfo').innerHTML += "<h6>Parking: " + evLocationInfo.parking + "<h6>";
    document.getElementById('displayInfo').innerHTML += "<h6>Hour: " + evLocationInfo.hour + "<h6>";
    document.getElementById('displayInfo').innerHTML += "<h6>Facility: " + evLocationInfo.facility + "<h6>";
    document.getElementById('displayInfo').innerHTML += "<h6>Website: <a href=" + evLocationInfo.website + ">" + evLocationInfo.website + "</a><h6>";
}

//get filter list
const getFilterList = async () => {
    const response = await fetch(URL + 'filter', {
        method: 'GET',
        cache: 'no-cache',
        headers: { 'Content-Type': 'application/json' },
        redirect: 'follow',
        insecure: true,
    });
    return response.json()
}
async function getFilter() {
    const filter = await getFilterList();
    displayFilter(filter);
}

function displayFilter(filter) {
    document.getElementById('displayInfo').innerHTML = "";
    document.getElementById("providerList").innerHTML = "";
    var select = document.getElementById("providerList");
    var option = document.createElement("option");
    var val = 0
    option.value = val.toString;
    option.text = "select provider";
    select.add(option, null);
    for (i = 0; i < filter.provider.length; i++) {
        option = document.createElement("option");
        val = i + 1
        option.value = val.toString;
        option.text = filter.provider[i];
        select.add(option, null);
    }
    option = document.createElement("option");
    var len = document.getElementById("providerList").options.length
    val = len + 1
    option.value = val.toString;
    option.text = "select all";
    select.add(option, null);
}

const getEVLocationAddressInfo = async (addr) => {
    const response = await fetch(URL + 'location/address?addr=' + addr, {
        method: 'GET',
        cache: 'no-cache',
        headers: { 'Content-Type': 'application/json' },
        redirect: 'follow',
        insecure: true,
    });

    return response.json()
}

async function searchAddress(map) {
    const address = document.getElementById("searchAddress").value;
    const evLocations = await getEVLocationAddressInfo(address);

    const newCenter = {
        lat: evLocations[0].location[0],
        lng: evLocations[0].location[1],
    };

    map.panTo(newCenter);
    map.setZoom(14);
    document.getElementById("searchAddress").value = "";
}


const getProviderLocationInfo = async (provider) => {
    const response = await fetch(URL + 'provider?provider=' + provider, {
        method: 'GET',
        cache: 'no-cache',
        headers: { 'Content-Type': 'application/json' },
        redirect: 'follow',
        insecure: true,
    });

    return response.json()
}

async function getProviderLocation(map) {
    const select = document.getElementById("providerList");
    const provider = select.options[select.selectedIndex].text;

    deleteFilter();

    if (provider === "select all") {
        showMarkers(map);
    } else {
        const evLocations = await getProviderLocationInfo(provider);
        clearMarkers();

        let url = "http://maps.google.com/mapfiles/ms/icons/";

        for (var i = 0; i < evLocations.length; i++) {
            var contentString = '<h6>' + evLocations[i].provider + '</h6>' +
                '<h6>' + evLocations[i].address + '</h6>';
            var color = "green";
            for (var j = 0; j < evLocations[i].charger.length; j++) {
                contentString += '<h6>' + evLocations[i].charger[j].type + ' - ' + evLocations[i].charger[j].station + ' station </h6>';
                if (evLocations[i].charger[j].match === true) {
                    color = "blue";
                };
            };

            const marker = new google.maps.Marker({
                position: { lat: evLocations[i].location[0], lng: evLocations[i].location[1] },
                map: map,
                icon: { url: url + color + ".png" }
            });

            filter.push(marker);
            setFilterOn(map);
            const infowindow = new google.maps.InfoWindow({
                content: contentString,
            });
            marker.addListener("mouseover", () => {
                infowindow.open(marker.get('map'), marker);
                InfoObj[0] = infowindow;
            });
            marker.addListener("mouseout", () => {
                closeOtherInfo();
            });

            marker.addListener("click", () => {
                const position = marker.getPosition();
                locationDetailInfo(position.lat(), position.lng());
            });
            map.panTo(pos);
            map.setZoom(12);
        }
    }
    getFilter();
}


// Sets the map on all markers in the array.
function setMapOnAll(map) {
    for (let i = 0; i < markers.length; i++) {
        markers[i].setMap(map);
    }
}

// Removes the markers from the map, but keeps them in the array.
function clearMarkers() {
    setMapOnAll(null);
}

// Shows any markers currently in the array.
function showMarkers(map) {
    setMapOnAll(map);
}

// Deletes all markers in the array by removing references to them.
function deleteMarkers() {
    clearMarkers();
    markers = [];
}


// Sets the map on all markers in the array.
function setFilterOn(map) {
    for (let i = 0; i < filter.length; i++) {
        filter[i].setMap(map);
    }
}
function showFilter(map) {
    setFilterOn(map);
}
function clearFilter() {
    setFilterOn(null);
}
function deleteFilter() {
    clearFilter();
    filter = [];
}