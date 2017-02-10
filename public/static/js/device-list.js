(function($, w) {
    function bindUIButtons() {
        $('#add-device-btn').click(function() {
            $('#edit-form').slideToggle();
        });
    }

    deviceForm.init({
        cancel: function() {
            $('#edit-form').slideUp(400);
        }
    });

    bindConfigClickEvents(); // common.js
    bindUIButtons();
})(jQuery, window);