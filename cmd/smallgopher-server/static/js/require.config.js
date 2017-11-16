require.config({
	baseUrl: "../js",
	paths: {
		bluebird: "./lib/bluebird",
		L: "./lib/leaflet/leaflet",
		jquery: "./lib/jquery/jquery-3.2.1.min",
		"jquery-ui": "./lib/jquery/jquery-ui.min",
		"jquery-buttonsetv": "./lib/jquery/plugins/jquery.buttonsetv",
		"fwm-client": "./fwm/fwm",
	},
	shim: {
		"jquery-buttonsetv": { deps: ["jquery", "jquery-ui"] },
	}
});