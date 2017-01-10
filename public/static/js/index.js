(function($, w) {
    var jobID = 0;

    var logLevels = [
        "DEBUG",
        "INFO",
        "NOTICE",
        "WARNING",
        "ERROR",
        "CRITICAL",
        "ALERT",
        "EMERGENCY",
        "FATAL",
    ];

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
            progress.prop('aria-valuemax', 0);
            progress.prop('aria-valuenow', 0);
            progress.css('width', '0%')
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
            statusLabel.text("Idle");
            statusLabel.removeClass();
            statusLabel.addClass('job-idle');

            progress.removeClass('active');
            progress.removeClass('progress-bar-striped');
        }, function(resp) {
            toastr["error"](resp.responseJSON.message);
        });
    }

    function updateUserLogTable(logs) {
        var table = $('#logs');
        var header = "<thead><tr><th>Level</th><th>Timestamp</th><th>Message</th></tr></thead>"
        var html = header + "<tbody>";

        for (var i = 0; i < logs.length; i++) {
            var log = logs[i];
            html += "<tr><td>" + logLevels[log.level] + "</td><td>" + formatDate(log.timestamp) + "</td><td>" + log.message + "</td></tr>";
        }

        html += "</tbody>";
        table.html(html);
    }

    API.getUserLogs(function(data) {
        updateUserLogTable(data.data);
    }, function() {
        console.log("Error getting logs")
    });
})(jQuery, window);