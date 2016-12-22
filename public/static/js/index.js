(function($, w) {
    var jobID = 0;

    $('#startArchiveBtn').click(function() {
        var spec = {
            type: "configs",
            devices: [],
            start: 0
        }

        API.startJob(spec, function(data) {
            toastr["success"]("Archive job started");
            jobID = data.data.id;

            var progress = $('#status-progress-bar');
            progress.addClass('active');
            progress.addClass('progress-bar-striped');

            var statusLabel = $('#current-status');
            statusLabel.text("Running");
            statusLabel.removeClass();
            statusLabel.addClass('job-running');

            startStatusChecker();
        }, function(resp) {
            toastr["error"](resp.responseJSON.message);
        });
    });

    function startStatusChecker() {
        API.jobStatus(jobID, function(data) {
            var progress = $('#status-progress-bar');
            progress.prop('aria-valuemax', data.data.total);
            progress.prop('aria-valuenow', data.data.completed);
            progress.css('width', ((data.data.completed / data.data.total) * 100) + '%')

            if (data.data.completed < data.data.total) {
                setTimeout(startStatusChecker, 3000);
                return
            }

            var statusLabel = $('#current-status');
            statusLabel.text("Finished");
            statusLabel.removeClass();
            statusLabel.addClass('job-finished');

            progress.removeClass('active');
            progress.removeClass('progress-bar-striped');
        }, function(resp) {
            toastr["error"](resp.responseJSON.message);
        });
    }
})(jQuery, window);