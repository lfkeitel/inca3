(function($, w) {
    var deviceForm = function() { };

    deviceForm.prototype.init = function(callbacks) {
        if (typeof callbacks === "undefined") {
            callbacks = {};
        }
        if (!("types" in callbacks)) {
            callbacks.types = function() { };
        }
        if (!("cancel" in callbacks)) {
            callbacks.cancel = function() { };
        }
        if (!("save" in callbacks)) {
            callbacks.save = function() {
                flashes.add("success", "Device saved");
                w.location.reload();
            };
        }

        populateTypes(callbacks.types);

        $('#cancel-edit-btn').click(function() {
            deviceForm.clear();
            callbacks.cancel();
        });

        $('#save-create-btn').click(function() {
            saveDevice(callbacks.save);
        });
    }

    function saveDevice(callback) {
        var d = {
            name: $('#device-name').val(),
            address: $('#device-addr').val(),
            profile: { id: Number($('#device-type').val()) }
        };

        API.createDevice(d, function(data) {
            callback();
        }, function(resp) {
            var json = resp.responseJSON;
            toastr["error"](json.message);
        });
    }

    deviceForm.prototype.clear = function() {
        $('#device-name').val("");
        $('#device-addr').val("");
        $('#device-type').val(1);
    }

    function populateTypes(callback) {
        API.getAllTypes(function(data) {
            var types = data.data;
            var typeSelect = $('#device-type');
            typeSelect.empty();

            for (var i = 0; i < types.length; i++) {
                var o = types[i];
                typeSelect.append($('<option>', {
                    value: o.id,
                    text: o.name
                }));
            }

            if (callback) { callback(); }
        }, function(j, t, e) {
            console.log(e);
            console.log(j.responseJSON);
        })
    }

    w.deviceForm = new deviceForm();
})(jQuery, window);