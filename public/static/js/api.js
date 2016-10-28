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

    w.API = new API();
})(window);