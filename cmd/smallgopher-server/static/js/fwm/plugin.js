var fwm = {};

define([
	"L",
    "bluebird",
    "fwm/base",
], function (L, bluebird, fwm) {
	 "use strict";
    fwm.Plugin = L.Class.extend({
        initialize: function(app, config, options){
            fwm.Class.prototype.initialize.call(this);
            this.app = app;
            this.map = app.map;
            this.config = config;
            this.options = options;
		}

    });
});