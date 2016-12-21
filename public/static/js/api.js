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

    API.prototype.getAllDevices = function(ok, fail) {
        $.getJSON("/api/devices").done(ok).fail(fail);
    }

    API.prototype.getDevice = function(id, ok, fail) {
        id = encodeURIComponent(id);
        $.getJSON("/api/devices/" + id).done(ok).fail(fail);
    }

    API.prototype.getDeviceConfigs = function(id, ok, fail) {
        id = encodeURIComponent(id);
        $.getJSON("/api/devices/" + id + "/configs").done(ok).fail(fail);
    }

    API.prototype.createDevice = function(device, ok, fail) {
        var json = JSON.stringify(device);
        makeRequest("POST", "/api/devices", json, ok, fail);
    }

    API.prototype.saveDevice = function(device, ok, fail) {
        var json = JSON.stringify(device);
        makeRequest("PUT", "/api/devices/" + device.slug, json, ok, fail);
    }

    API.prototype.deleteDevice = function(slug, ok, fail) {
        slug = encodeURIComponent(slug);
        makeRequest("DELETE", "/api/devices/" + slug, "", ok, fail);
    }

    API.prototype.getAllTypes = function(ok, fail) {
        $.getJSON("/api/types").done(ok).fail(fail);
    }

    API.prototype.startJob = function(spec, ok, fail) {
        var json = JSON.stringify(spec);
        makeRequest("POST", "/api/job/start", json, ok, fail);
    }

    API.prototype.jobStatus = function(id, ok, fail) {
        $.getJSON("/api/job/status/" + id).done(ok).fail(fail);
    }

    w.API = new API();
})(window);