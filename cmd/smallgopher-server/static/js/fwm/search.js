define([
    "jquery",
    "L",
    "fwm/base",
    "fwm/plugin"
], function ($, L, fwm) {
    "use strict";
    fwm.Search = fwm.Plugin.extend({
        init: function(){
            this.buildSearchPane();
            this.addEventHandlers();
        },
        buildSearchPane: function () {
            this.searchDialog = $("#search-dialog").dialog({
                autoOpen: false,
                title: "Search",
                width: 800,
                height: 70,
                position: {
                    my: "middle top",
                    at: "middle top"
                }
            });
            $("#search-button").button();
            var options = [],
                conf = this.config,
                names = conf.featureNames();
            for (var i = 0; i < names.length; i++) {
                var name = names[i],
                    feature = conf.feature(name);
                options.push($("<option>", {
                    value: feature.name,
                    html: feature.display_name
                }));
            }
            $("#search-feature").append(options);
            this.updateSearchFields(names[0]);
        },
        updateSearchFields: function (value) {
            var conf = this.config,
                fields = conf.featureFields(value),
                options = [];
            for (var i = 0; i < fields.length; i++) {
                var field = fields[i];
                if (!field.searchable) {
                    continue;
                }
                options.push($("<option>", {
                    value: field.name,
                    html: field.display_name
                }))
            }
            var enable = options.length === 0;
            $("#search-field").prop("disabled", enable);
            $("#search-value").prop("disabled", enable);
            $("#search-button").button("option", "disabled", enable);

            $("#search-field").html("").append(options);
        },
        addEventHandlers: function(){
            $("#bt-search").on("click", this.openSearchDialog.bind(this));
            $("#search-feature").on("change", function (e) {
                this.updateSearchFields(e.target.value);
            }.bind(this));
            $("#search-button").on("click", this.searchFeature.bind(this));
        },
        openSearchDialog: function () {
            this.searchDialog.dialog("open");
        },
        searchFeature: function () {
            var name = $("#search-feature").val(),
                field = $("#search-field").val(),
                value = $("#search-value").val(),
                data = JSON.stringify({
                    "name": name,
                    "fieldName": field,
                    "values": [value]
                });
            $.post("/search", data).then(function (data) {
                $("#search-result").html("");
                if (data.error) {
                    alert(data.error)
                    return;
                }
                var results = [];
                var features = data.features;
                for (var i = 0; i < features.length; i++) {
                    var feature = features[i];
                    var resultRow = $("<div>", { html: feature.properties[field], style: "cursor:pointer;" });
                    resultRow.on("click", function () {
                        var keyField = this.config.featureKey(name),
                            featureKey = feature.properties[keyField];
                        this.app.get("/feature/" + name + "/" + featureKey).then(function (data) {
                            if (data.error) {
                                alert(data.error);
                                return;
                            }
                            this.app.setCurrentFeature(data.features[0], name);
                        }.bind(this))
                    }.bind(this));
                    results.push(resultRow);
                }
                $("#search-result").append(results);

            }.bind(this))
        }

    })
});