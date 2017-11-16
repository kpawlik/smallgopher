define([
    "jquery",
    "L",
    "fwm/base",
    "fwm/plugin"
], function ($, L, fwm) {
    "use strict";
    fwm.Layers = fwm.Plugin.extend({
        init: function () {
            this.layers = this.app.layers;
            this.buildLayersPane();
            this.addEventHandlers();
            this.updateLayers();
        },
        /**
         * buildLayersPane builds GUI panel for avaliable map layers
         */
        buildLayersPane: function () {
            var conf = this.config,
                names = conf.featureNames(),
                fs = $("<fieldset>"),
                layers = [];
            for (var i = 0; i < names.length; i++) {
                var layer = conf.feature(names[i]);
                var id = "layer-id-" + layer.name;
                layers.push($("<label>", { "for": id, "title": "Enable/Disable layer" }).html(layer.display_name + " ( " + layer.min_zoom + " - " + layer.max_zoom + " ) "));
                layers.push($("<input>", { "type": "checkbox", "name": layer.name, id: id, "title": "Enable/Disable layer" }));
            }
            fs.append(layers);
            $("#layers-pane").append(fs);
            $("#layers-pane > fieldset > input").button();
            $("#layers-pane > fieldset > label").css("text-align", "left");
            $("#layers-pane-dialog").dialog({
                width: 230,
                minWidth: 200,
                height: 350,
                position: {
                    my: "left top",
                    at: "left top+100"
                },
                title: "Layers",
                autoOpen: false
            });
        },
        addEventHandlers: function () {
            $("#layers-pane").find("input").on('change', function (e) {
                var input = $(e.target),
                    name = input.attr("name");
                if ($(input).is(':checked')) {
                    this.app.displayFeatures(name);
                } else {
                    var current = this.layers[name];
                    if (current) {
                        this.map.removeLayer(current);
                    }
                }
            }.bind(this));
            $("#layers-pane > fieldset > button").on("click", function (e) {
                var bt = e.currentTarget;
                var name = $(bt).attr("name");
                var layer = this.layers[name];
                if (layer && this.map.hasLayer(layer)) {
                    layer.bringToFront();
                }
            }.bind(this));
            this.map.on("moveend", this.updateLayers.bind(this));
            $("#bt-layers").on("click", function () { $("#layers-pane-dialog").dialog("open"); });
        },

        getVisibleLayersNames: function () {
            return $("#layers-pane").find("input:checked").map(function () { return $(this).attr("name") }).toArray();
        },

        updateLayers: function () {
            var zoom = this.map.getZoom();
            $("#control-pane").dialog({ "title": "Control (map zoom: " + zoom + ")" })
            var featureNames = this.config.featureNames();
            $("#layers-pane").find("input").button({ "disabled": true });
            $("#layers-pane").find("button").button({ "disabled": true });

            for (var i = 0; i < featureNames.length; i++) {
                var featureName = featureNames[i];
                var featureConf = this.config.feature(featureName);
                if (zoom <= featureConf.max_zoom && zoom >= featureConf.min_zoom) {
                    $("#layer-id-" + featureName).button({ "disabled": false });
                    $("#layer-id-" + featureName + "btf").button({ "disabled": false });
                }
            }
            var layerNames = this.getVisibleLayersNames();
            for (var i = 0; i < layerNames.length; i++) {
                this.app.displayFeatures(layerNames[i]);
            }
        },
    });
});
