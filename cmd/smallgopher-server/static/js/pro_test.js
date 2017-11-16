define([
   
], function () {
    "use strict";
    var Pro = {

        diag: function (app) {
            test(app);
            
        }
    };

    //
    //
    var insertWrap = function (app) {
        return new Promise(function (res, rej) {
            console.log("Promise body. Diagnose app:");
            if (app) {
                res(app);
            } else {
                rej("Error");
            }
        });
    }
    var catchWrap = function (exc) {
        console.log("EXC: ");
        console.log(exc);
        $("#msg").html("ERROR!");
    }
    var insert = function (app) {
        return Promise.resolve(insertWrap(app)).catch(catchWrap);
    }
    var f1 = function (data) {
        console.log("F1. App name:" + data.name);
        return data;
    }
    var f2 = function (data) {
        console.log("F2. App map:" + data.map);
        return data;
    }
    var f3 = function (data) {
        console.log("F3. Config:" + data.config);
        return data;
    }
    var test = function (app) {
        return insert(app).then(f1).then(f2).then(f3);
    }
    //
    //


    return Pro;
});