(function($, w) {
    var thisDevice = {};

    function bindUIButtons() {
        $('#edit-device-btn').click(function() {
            $('#edit-form').slideToggle();
        });

        $('#cancel-edit-btn').click(function() {
            $('#edit-form').slideUp(400, populateEditForm);
        });

        $('#delete-device-btn').click(function() {
            deleteDevice();
        });

        $('#save-create-btn').click(function() {
            saveDevice();
        });
    }

    function saveDevice() {
        var d = {
            name: $('#device-name').val(),
            address: $('#device-addr').val(),
            slug: thisDevice.slug,
            profile: { id: Number($('#device-type').val()) }
        };

        API.saveDevice(d, function(data) {
            if (data.data.slug !== thisDevice.slug) {
                flashes.add("success", "Device saved");
                window.location = "devices/" + data.data.slug;
                return;
            }
            thisDevice = data.data;
            toastr["success"]("Device Saved");
            $('#edit-form').slideUp();
        }, function(resp) {
            var json = resp.responseJSON;
            toastr["error"](json.message);
        });
    }

    function getThisDevice() {
        thisDevice = JSON.parse($('#device-json').html());
    }

    function populateEditForm() {
        $('#device-name').val(thisDevice.name);
        $('#device-addr').val(thisDevice.address);
        $('#device-type').val(String(thisDevice.profile.id))
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

    function deleteDevice() {
        var confirm = new jsConfirm();
        confirm.show("Are you sure you want to delete this device?<br>This will also delete configurations for this device", function() {
            API.deleteDevice(thisDevice.slug, function() {
                flashes.add('success', "Device deleted");
                window.location = "devices";
                return
            }, function(resp) {
                toastr["error"](resp.responseJSON.message);
            })
        });
    }

    bindConfigClickEvents(); // common.js
    bindUIButtons();
    getThisDevice();
    populateTypes(populateEditForm);
})(jQuery, window);