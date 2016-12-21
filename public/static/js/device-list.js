(function($, w) {
    function bindUIButtons() {
        $('#add-device-btn').click(function() {
            $('#edit-form').slideToggle();
        });

        $('#cancel-edit-btn').click(function() {
            $('#edit-form').slideUp(400, clearDeviceForm);
        });

        $('#save-create-btn').click(function() {
            saveDevice();
        });
    }

    function saveDevice() {
        var d = {
            name: $('#device-name').val(),
            address: $('#device-addr').val(),
            type: { id: Number($('#device-type').val()) }
        };

        API.createDevice(d, function(data) {
            flashes.add("success", "Device saved");
            w.location.reload();
        }, function(resp) {
            var json = resp.responseJSON;
            toastr["error"](json.message);
        });
    }

    function clearDeviceForm() {
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

    bindConfigClickEvents(); // common.js
    bindUIButtons();
    populateTypes();
})(jQuery, window);