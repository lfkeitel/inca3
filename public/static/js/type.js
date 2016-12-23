(function($, w) {
    var thisType = {};

    function bindUIButtons() {
        $('#cancel-edit-btn').click(function() {
            window.location = "profiles";
        });

        $('#delete-type-btn').click(function() {
            deleteType();
        });

        $('#save-create-btn').click(function() {
            saveType();
        });
    }

    function saveType() {
        var t = {
            slug: thisType.slug,
            name: $('#type-name').val(),
            username: $('#type-username').val(),
            password: $('#type-password').val(),
            enablepw: $('#type-enablepw').val(),
            script: $('#type-script').val()
        };

        API.saveType(t, function(data) {
            flashes.add("success", "Device type saved");
            console.log(data);
            window.location = "profiles";
        }, function(resp) {
            var json = resp.responseJSON;
            toastr["error"](json.message);
        });
    }

    function getThisType() {
        thisType = JSON.parse($('#type-json').html());
    }

    function populateEditForm() {
        $('#type-name').val(thisType.name);
        $('#type-username').val(thisType.username);
        $('#type-password').val(thisType.password);
        $('#type-enablepw').val(thisType.enablepw);
        $('#type-script').val(thisType.script);
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

    function deleteType() {
        var confirm = new jsConfirm();
        confirm.show("Are you sure you want to delete this type?", function() {
            API.deleteType(thisType.slug, function() {
                flashes.add('success', "Type deleted");
                window.location = "profiles";
                return
            }, function(resp) {
                toastr["error"](resp.responseJSON.message);
            })
        });
    }

    bindConfigClickEvents(); // common.js
    bindUIButtons();
    getThisType();
    populateScripts(populateEditForm);
})(jQuery, window);