(function($, w) {
    $('#startArchiveBtn').click(function() {
        var spec = {
            type: "configs",
            devices: [],
            start: 0
        }

        API.startJob(spec, function(data) {
            console.log(data);
        }, function(resp) {
            console.log(resp.responseJSON);
        });
    });
})(jQuery, window);