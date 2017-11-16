define([
	"fwm/base"
], function (fwm) {
	"use strict";

	fwm.Config = fwm.Class.extend({
        initialize: function(config){
            fwm.Class.prototype.initialize.call(this);
            this.config = config;
            this.featuresKeys = {};
        },
        feature: function(name){
            var conf = this.config;
            for(var i=0; i< conf.length; i++){
                if(conf[i].name === name){
                    return conf[i];
                }
            }
        },
        addFeature: function(name, config){
            this.config[name] = config;
        },
        featureKey: function(name){
            var fkey = this.featuresKeys[name];
            if (fkey){
                return fkey;
            }
            var fconfig = this.feature(name),
                fields = fconfig.fields;

            for(var i=0; i < fields.length; i++){
                if(fields[i].key === true){
                    fkey = fields[i].name;
                    this.featuresKeys[name] = fkey;
                    return fkey;
                }
            }
        },
        featureFields:function(name){
            return this.feature(name).fields;
        },
        featureDataset:function(name){
            return this.feature(name).dataset;
        },
        featureDisplayName: function(name){
            return this.feature(name).display_name;
        },
        featureNames: function(){
            return this.config.map(function(f){
                return f.name
                ;})
        },
        fieldDisplayNames: function(name){
            var config = this.feature(name),
                fieldsNames = {};
            for (var i = 0; i < config.fields.length; i++) {
                fieldsNames[config.fields[i].name] = config.fields[i].display_name;
            }
            return fieldsNames;
        },
        fieldNames: function(name){
            var config = this.feature(name),
                fieldsNames = [];
            for (var i = 0; i < config.fields.length; i++) {
                if(config.fields[i].geom){
                    continue;
                }
                fieldsNames.push(config.fields[i].name);
            }
            return fieldsNames;
        },
        searchableFields: function(name){
            var config = this.feature(name),
                fieldsNames = [];
            for (var i = 0; i < config.fields.length; i++) {
                if(config.fields[i].searchable){
                    fieldsNames.push(config.fields[i].name);
                }
            }
            return fieldsNames;
        }

	});
	return fwm.Config;
});