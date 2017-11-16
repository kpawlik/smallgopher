define([
    "jquery",
    "L",
    "fwm/base",
    "style",
    "jquery-ui",
    "jquery-buttonsetv",
    "fwm/cache",
    "fwm/config",
    "fwm/search",
    "fwm/layers",
], function ($, L, fwm, gstyle) {
    "use strict";
    fwm.Application = L.Class.extend({
        pluginsDefs: {
            "search": fwm.Search,
            "layers": fwm.Layers
        },
        initialize: function (name, center) {
            this.plugins = {};
            this.name = name;
            this.center = center;
            this.config = null;
            this.layers = {};
            this.cache = new fwm.Cache();

        },
        run: function () {
            var mbAttr = 'Fun with Maps!',
                mbUrl = 'https://api.tiles.mapbox.com/v4/{id}/{z}/{x}/{y}.png?access_token=pk.eyJ1IjoibWFwYm94IiwiYSI6ImNpejY4NXVycTA2emYycXBndHRqcmZ3N3gifQ.rJcFIG214AriISLbB6B5aw';

            var streets = L.tileLayer(mbUrl, {
                id: 'mapbox.streets',
                attribution: mbAttr,
                maxZoom: 20
            });

            this.map = L.map('mapid', {
                //center: [32.72435148682413, -117.16269111601927],
                //center: [52.22184111652178, 0.10957786174227713],
                center: this.center,
                zoom: 17,
                layers: [streets],
                maxZoom: 20
            });
            this.get('/config').then(
                function (config) {
                    if (config.error) {
                        alert(config.error);
                        return;
                    }
                    this.config = new fwm.Config(config);
                    this.initGui();

                }.bind(this)
            );
            $("#mapid").css("width", "100%");
        },
        /**
         * getFeatureConfig returns config object for feature with the 'name'
         */
        getFeatureConfig: function (name) {
            var conf = this.config.feature(name);
            if (conf) {
                return new Promise(function (res, rej) { res(conf); });
            }
            return this.get("/config/" + name).then(function (data) {
                this.config.addFeature(name, data);
                return data;
            }.bind(this));
        },

        initPlugins: function () {
            var defs = this.pluginsDefs,
                defKeys = Object.keys(defs);
            for (var i = 0; i < defKeys.length; i++) {
                var pluginName = defKeys[i];
                var pluginClass = defs[pluginName];
                var plugin = new pluginClass(this, this.config, {});
                plugin.init();
                this.plugins[pluginName] = plugin;
            }
        },



        buildControlPane: function () {
            $("#control-pane").dialog({
                closeOnEscape: false,
                open: function (event, ui) {
                    $(".ui-dialog-titlebar-close", ui.dialog | ui).hide();
                },
                title: "Control",
                width: 220,
                height: 70,
                position: {
                    my: "right top",
                    at: "right top"
                },
            });
            $("#bt-layers").button();
            $("#bt-viewer").button();
            $("#bt-search").button();
        },

        /**
         * addEventHandlers add event handlers to GUI elements
         */
        addEventHandlers: function () {

            $("#bt-viewer").on("click", function () { $("#feature-viewer-dialog").dialog("open"); });


        },
        /**
         * Initialize GUI
         */
        initGui: function () {
            this.buildControlPane();
            $("#feature-viewer-dialog").dialog({
                width: 230,
                minWidth: 200,
                height: 200,
                position: {
                    my: "left bottom",
                    at: "left bottom-10px"
                },
                title: "Details",
                autoOpen: false
            });
            this.addEventHandlers();
            this.initPlugins();
        },

        get: function (url) {
            return $.get(url).catch(function (resp) {
                if (resp.responseText) {
                    try {
                        var json = JSON.parse(resp.responseText);
                        alert(json.error);
                    } catch (e) {
                        alert(resp.responseText);
                    }

                }
            });
        },



        displayFeatureDetails: function (feature, name) {
            this.clearFeatureDetails();
            var config = this.config.feature(name),
                featureTitle = config.display_name,
                props = feature.properties,
                fieldNames = this.config.fieldDisplayNames(name),
                fields = this.config.fieldNames(name),
                arr = ['<table class="feature-table">'];
            arr.push("<th colspan=2>", featureTitle, "</th>");
            for (var i = 0; i < fields.length; i++) {
                var field = fields[i];
                arr.push('<tr class="feature-table-row">');
                arr.push('<td>');
                arr.push(fieldNames[field])
                arr.push('</td>');
                arr.push('<td>');
                arr.push(props[field])
                arr.push('</td>');
                arr.push('</tr>');
            }
            arr.push('</table>');
            $("#feature-viewer").html(arr.join(""));
            $("#feature-viewer-dialog").dialog("open");
        },
        setCurrentFeature(feature, name) {
            this.displayFeatureDetails(feature, name);
            this.highlight(feature);
            this.currentFeature = feature;
        },
        removeCurrentFeature() {
            this.clearFeatureDetails();
            this.unHighlight();
            this.currentFeature = null;
        },
        clearFeatureDetails: function () {
            $("#feature-viewer").html("");
        },
        highlight: function (feature) {
            this.unHighlight();
            var geom = null;
            if (feature.geometry.type == "LineString") {
                var latLngs = feature.geometry.coordinates.map(function (latLng) { return L.latLng(latLng[1], latLng[0]) });
                geom = L.polyline(latLngs, gstyle.highlightLine);
            }
            if (feature.geometry.type == "Point") {
                var coord = feature.geometry.coordinates;
                var latLng = L.latLng(coord[1], coord[0]);
                geom = L.circleMarker(latLng, gstyle.highlightPoint);
            }
            if (geom) {
                geom.addTo(this.map);
                this.highlightGeom = geom;
            }
        },
        unHighlight: function () {
            if (!this.highlightGeom) {
                return;
            }
            this.map.removeLayer(this.highlightGeom);
        },

        /**
        * Get JSON from server and add it to map as layer with styles and popups
        */
        displayFeatures: function (name) {
            this.getFeatureConfig(name).then(function (config) {
                var bbox = this.map.getBounds().toBBoxString();
                this._displayFeatures(name, bbox, 0, config);
            }.bind(this))
        },
        /**
         * Get JSON from server and add it to map as layer with styles and popups
         */
        _displayFeatures: function (name, bbox, limit, config) {
            if (!this.layerIsVisible(name)) {
                var layer = this.layers[name];
                if (layer) {
                    this.addLayerIfNeeded(layer, name);
                }
                return;
            }
            var key = this.config.featureKey(name);
            this.get("/features/" + name + "/" + bbox + "/" + limit).then(function (data) {
                if (data.error) {
                    alert(data.error);
                    return;
                }
                if (!data.features) {
                    return;
                }
                var layer = this.layers[name];
                if (!layer) {
                    // add all features to cache
                    this.cache.addAll(name, data.features.map(function (feature) { return feature.properties[key] }));
                    layer = this.layers[name] = L.geoJSON(data, {
                        style: gstyle[name],
                        pointToLayer: gstyle[name],
                        onEachFeature: function (feature, layer) {
                            layer.on("click", function (e) {
                                var featureKey = feature.properties[key];
                                this.get("/feature/" + name + "/" + featureKey).then(function (data) {
                                    if (data.error) {
                                        alert(data.error);
                                        return;
                                    }
                                    this.setCurrentFeature(data.features[0], name);
                                }.bind(this));
                            }.bind(this));
                        }.bind(this)
                    });
                    this.addLayerIfNeeded(layer, name);
                } else {
                    this.addLayerIfNeeded(layer, name);
                    var cache = this.cache;
                    for (var i = 0; i < data.features.length; i++) {
                        var feature = data.features[i];
                        var featureKey = feature.properties[key];
                        if (!cache.record(name, featureKey)) {
                            layer.addData(feature);
                            this.cache.add(name, featureKey);
                        }
                    }

                }
            }.bind(this));
        },
        layerIsChecked: function (name) {
            return $("#layer-id-" + name).is(':checked')
        },
        layerIsVisible: function (name) {
            var zoom = this.map.getZoom();
            var featureConf = this.config.feature(name);
            return zoom <= featureConf.max_zoom && zoom >= featureConf.min_zoom;
        },
        addLayerIfNeeded: function (layer, name) {
            if (this.layerIsChecked(name) && this.layerIsVisible(name)) {
                layer.addTo(this.map);
            } else {
                this.map.removeLayer(layer);
            }
        },



    });

    return fwm.Application;
});