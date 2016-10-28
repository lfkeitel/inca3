var delimiters = ['${', '}'];

var filters = {
    capitalize: function(value) {
        if (!value) return ''
        value = value.toString()
        return value.charAt(0).toUpperCase() + value.slice(1)
    },
    formatDate: function(value) {
        var date = new Date(value * 1000);
        var month = date.getMonth() + 1;
        var day = date.getDate();
        var year = date.getFullYear();

        var hour = "0" + date.getHours();
        var minute = "0" + date.getMinutes();
        var second = "0" + date.getSeconds();

        return month + "/" + day + "/" + year + " " +
            hour.slice(-2) + ":" + minute.slice(-2) + ":" + second.slice(-2);
    },
    formatFileSize: function(value) {
        var sizes = ["B", "KB", "MB", "GB"];
        var i = 0;
        for (i = 0; i < sizes.length; i++) {
            if (value < 1024) {
                break;
            }
            value /= 1024;
        }
        return value + sizes[i];
    }
}