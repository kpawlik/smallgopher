define([],
    function () {
        "use strict";
        var defaultStyle = function (feature) {
            return {
                weight: 3,
                color: 'black',
                dashArray: '3',
                fillColor: '#ff0000'
            };
        }
        var sheatheStyle = function (feature) {
            return {
                weight: 3,
                color: 'black',
            };
        }
        var coaxialCableStyle = function (feature) {
            return {
                weight: 3,
                color: 'blue',
            };
        }
        var mit_hub = function (feature, latlng) {
            return L.marker(latlng, {
                icon: L.icon({
                    iconUrl: "/img/hub.png",
                    iconSize: [24, 24],
                    opacity: 0.8
                })
            });
        }
        var amplifier = function (feature, latlng) {
            return L.marker(latlng, {
                icon: L.icon({
                    iconUrl: "/img/amplifier.png",
                    iconSize: [24, 24]
                })
            });
        }
        var tap = function (feature, latlng) {
            return L.circleMarker(latlng, {
                radius: 8,
                color: "#555",
                opacity: 1.0,
                fillOpacity: 0
            });
        }
        var te = function (feature, latlng) {
            return L.circleMarker(latlng, {
                radius: 5,
                color: "#A22",
                opacity: 1.0,
                fillOpacity: 0
            });
        }
        return {
            "sheath": sheatheStyle,
            "mit_hub": mit_hub,
            "tap": tap,
            "amplifier": amplifier,
            "te": te,
            "coaxial_cable": coaxialCableStyle
        }
    });