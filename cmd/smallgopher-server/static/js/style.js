define([
    "jquery",
    "GeoJSONstyles/pni",
    "GeoJSONstyles/cambridge"
],
    function ($, pni, cambridge) {
        "use strict";
        var defaultStyle = function (feature) {
            return {
                weight: 3,
                color: 'black',
                dashArray: '3',
                fillColor: '#ff0000'
            };
        }
        
        var highlightLine = {
            weight: 10,
            color: 'orange',
            opacity: 0.5
        }
        var highlightPoint = {
            radius: 12,
            color: 'black',
            fillColor: "orange",
            opacity: 0.6,
        }
        
        var style = {
            "highlightLine": highlightLine,
            "highlightPoint": highlightPoint
        };
        $.extend(true, style, pni, cambridge);
        return style;
    }

)