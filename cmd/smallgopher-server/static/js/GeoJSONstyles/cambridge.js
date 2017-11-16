define([],
    function () {
        "use strict";
        var pl_of_interest = function (feature, latlng) {
            return L.marker(latlng, {
                icon: L.icon({
                    iconUrl: "/img/poi.png",
                    iconSize: [24, 24],
                    opacity: 1.0
                })
            });
        }
        var hotel = function (feature, latlng) {
            return L.marker(latlng, {
                icon: L.icon({
                    iconUrl: "/img/hotel.png",
                    iconSize: [24, 24],
                    opacity: 1.0
                })
            });
        }
        var car_park = function (feature, latlng) {
            return L.marker(latlng, {
                icon: L.icon({
                    iconUrl: "/img/carpark.png",
                    iconSize: [24, 24],
                    opacity: 0.8
                })
            });
        }
        var rail_line = function () {
            return {
                stroke: true,
                weight: 4,
                color: 'black',
                dashArray: '15',
                fillColor: '#FFFFFF',

            };
        }
        var office = function (feature, latlng) {
            return L.marker(latlng, {
                icon: L.icon({
                    iconUrl: "/img/office.png",
                    iconSize: [24, 24],
                    opacity: 0.8
                })
            });
        }
        var town = function () {
            return {
                fillColor: "#555555",
                opacity: 0.4,
                stroke: false
            }
        }
        var park = function () {
            return {
                fillColor: "#0e5e2b",
                fillOpacity: 0.7,
                stroke: true,
                color: "#064f20",
            }
        }
        var footpath = function(){
            return {
                weight: 2,
                color: "#667788"
            }
        }
        var min_road = function(){
            return {
                weight: 4,
                color: "#443388"
            }
        }
        var slip_road = function(){
            return {
                weight: 3,
                color: "#4C4C4C"
            }
        }
        var trunk_road = function(){
            return {
                weight: 4,
                color: "#565654"
            }
        }
        return {
            "pl_of_interest": pl_of_interest,
            "hotel": hotel,
            "car_park": car_park,
            "rail_line": rail_line,
            "office": office,
            "town": town,
            "park": park,
            "footpath": footpath,
            "min_road": min_road,
            "slip_road": slip_road,
            "trunk_road": trunk_road,
        }
    }

)