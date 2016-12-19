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

toastr.options = {
    "closeButton": true,
    "debug": false,
    "newestOnTop": true,
    "progressBar": false,
    "positionClass": "toast-top-center",
    "preventDuplicates": false,
    "onclick": null,
    "showDuration": "300",
    "hideDuration": "1000",
    "timeOut": "5000",
    "extendedTimeOut": "1000",
    "showEasing": "swing",
    "hideEasing": "linear",
    "showMethod": "fadeIn",
    "hideMethod": "fadeOut"
}

var flashes = {
    lsKey: "flashes",

    add: function(type, message) {
        var lsFlashes = flashes.getFlashes();
        lsFlashes.push({
            type: type,
            message: message
        });
        flashes.save(lsFlashes);
    },

    showAll: function() {
        var lsFlashes = flashes.getFlashes();
        while (lsFlashes.length > 0) {
            var o = lsFlashes.pop();
            toastr[o.type](o.message);
        }
        flashes.clear();
    },

    getFlashes: function() {
        var lsFlashes = window.localStorage.getItem(flashes.lsKey);
        if (lsFlashes === null || lsFlashes === '') {
            return [];
        }
        return JSON.parse(lsFlashes);
    },

    clear: function() {
        window.localStorage.setItem(flashes.lsKey, "[]");
    },

    save: function(f) {
        window.localStorage.setItem(flashes.lsKey, JSON.stringify(f));
    }
}

flashes.showAll();