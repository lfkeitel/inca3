(function($, w) {
    function bindUIButtons() {
        $('#add-type-btn').click(function() {
            $('#edit-form').slideToggle();
        });

        $('#cancel-edit-btn').click(function() {
            $('#edit-form').slideUp(400, clearTypeForm);
        });

        $('#save-create-btn').click(function() {
            saveType();
        });
    }

    function saveType() {
        var t = {
            name: $('#type-name').val(),
            brand: $('#type-brand').val(),
            connection: $('#type-connection').val(),
            script: $('#type-script').val()
        };

        API.createType(t, function(data) {
            flashes.add("success", "Device type saved");
            w.location.reload();
        }, function(resp) {
            var json = resp.responseJSON;
            toastr["error"](json.message);
        });
    }

    function clearTypeForm() {
        $('#type-name').val("");
        $('#type-brand').val("");
        $('#type-connection')[0].selectedIndex = 0;
        $('#type-script')[0].selectedIndex = 0;
    }

    function populateScripts(callback) {
        API.getTypeScripts(function(data) {
            var scripts = data.data;
            var scriptSelect = $('#type-script');
            scriptSelect.empty();

            for (var i = 0; i < scripts.length; i++) {
                var o = scripts[i];
                scriptSelect.append($('<option>', {
                    value: o,
                    text: o
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
    populateScripts();
})(jQuery, window);