(function(w) {
    function makeRequest(method, url, data, ok, fail) {
        $.ajax({
            contentType: "application/json",
            data: data,
            dataType: "json",
            method: method,
            url: url
        }).done(ok).fail(fail);
    }

    var API = function() { };

    // Device APIs
    API.prototype.getAllDevices = function(ok, fail) {
        $.getJSON("api/devices").done(ok).fail(fail);
    }

    API.prototype.getDevice = function(id, ok, fail) {
        id = encodeURIComponent(id);
        $.getJSON("api/devices/" + id).done(ok).fail(fail);
    }

    API.prototype.getDeviceConfigs = function(id, ok, fail) {
        id = encodeURIComponent(id);
        $.getJSON("api/devices/" + id + "/configs").done(ok).fail(fail);
    }

    API.prototype.createDevice = function(device, ok, fail) {
        var json = JSON.stringify(device);
        makeRequest("POST", "api/devices", json, ok, fail);
    }

    API.prototype.saveDevice = function(device, ok, fail) {
        var json = JSON.stringify(device);
        makeRequest("PUT", "api/devices/" + device.slug, json, ok, fail);
    }

    API.prototype.deleteDevice = function(slug, ok, fail) {
        slug = encodeURIComponent(slug);
        makeRequest("DELETE", "api/devices/" + slug, "", ok, fail);
    }

    // Type APIs
    API.prototype.getAllTypes = function(ok, fail) {
        $.getJSON("api/profiles").done(ok).fail(fail);
    }

    API.prototype.getTypeScripts = function(ok, fail) {
        $.getJSON("api/profiles/_scripts").done(ok).fail(fail);
    }

    API.prototype.createType = function(type, ok, fail) {
        var json = JSON.stringify(type);
        makeRequest("POST", "api/profiles", json, ok, fail);
    }

    API.prototype.saveType = function(type, ok, fail) {
        var json = JSON.stringify(type);
        makeRequest("PUT", "api/profiles/" + type.slug, json, ok, fail);
    }

    API.prototype.deleteType = function(slug, ok, fail) {
        slug = encodeURIComponent(slug);
        makeRequest("DELETE", "api/profiles/" + slug, "", ok, fail);
    }

    // Job APIs
    API.prototype.startJob = function(spec, ok, fail) {
        var json = JSON.stringify(spec);
        makeRequest("POST", "api/job/start", json, ok, fail);
    }

    API.prototype.jobStatus = function(id, ok, fail) {
        $.getJSON("api/job/status/" + id).done(ok).fail(fail);
    }

    w.API = new API();
})(window);