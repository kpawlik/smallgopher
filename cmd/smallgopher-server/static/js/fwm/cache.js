define([
	"fwm/base"
], function (fwm) {
	"use strict";

	fwm.Cache = fwm.Class.extend({
        initialize: function(){
            fwm.Class.prototype.initialize.call(this);
            this.cache = {};
        },
        collection: function(name){
            var coll = this.cache[name];
            if(!coll){
                coll = this.cache[name] = {};
            }
            return coll;
        },
        record: function(name, id){
            return this.collection(name)[id];
        },
        add: function(name, id){
            this.collection(name)[id] = true;
        },
        addAll: function(name, ids){
            var collection = this.collection(name);
            for(var i=0; i<ids.length; i++){
                collection[ids[i]] = true
            }
        }

	});
	return fwm.Cache;
});