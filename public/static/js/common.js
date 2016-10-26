(function($) {
    var delimiters = ['${', '}'];

    var filters = {
        capitalize: function(value) {
            if (!value) return ''
            value = value.toString()
            return value.charAt(0).toUpperCase() + value.slice(1)
        }
    }

    Vue.component("device-list", {
        template: "#device-list-template",
        delimiters: delimiters,
        props: {
            tableData: {
                type: Array,
                required: true
            },
            tableCols: {
                type: Array,
                required: true
            }
        },
        filters: filters,
        methods: {
            gotoDevice: function(id) {
                window.location = "/devices/" + id;
            }
        }
    });

    Vue.component("config-list", {
        template: "#config-list-template",
        delimiters: delimiters,
        props: {
            tableData: {
                type: Array,
                required: true
            },
            tableCols: {
                type: Array,
                required: true
            }
        },
        filters: filters,
        methods: {
            gotoConfig: function(id) {
                window.location = "/devices/1/" + id;
            }
        }
    });

    var vm = new Vue({
        el: "#app",
        delimiters: delimiters,
        data: {
            tableColumns: ["name", "address", "connection", "brand"],
            tableData: []
        }
    });

    API.getAllDevices(function(data) {
        vm.tableData = data.data;
    }, function(j, t, e) {
        console.error(e);
    })
})(jQuery);