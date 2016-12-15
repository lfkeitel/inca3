(function(w) {
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

    API.prototype.saveDevice = function(device, ok, fail) {
        // Make a copy of the object so we don't affect a vue data member
        var device = JSON.parse(JSON.stringify(device));
        if (device.id === 0) {
            this.createDevice(device, ok, fail);
        } else {
            this.editDevice(device, ok, fail);
        }
    }

    API.prototype.createDevice = function(device, ok, fail) {
        delete device.configs;
        var json = JSON.stringify(device);
        $.post("/api/devices", json, null, "json").done(ok).fail(fail);
    }

    API.prototype.editDevice = function(device, ok, fail) {
        delete device.configs;
        var json = JSON.stringify(device);
        $.ajax({
            contentType: "application/json",
            data: json,
            dataType: "json",
            method: "PUT",
            url: "/api/devices/" + device.slug
        }).done(ok).fail(fail);
    }

    w.API = new API();
})(window);