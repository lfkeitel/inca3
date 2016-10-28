(function($) {
    Vue.component("device-list", {
        template: "#device-list-template",
        delimiters: delimiters,
        props: {
            tableData: {
                type: Array,
                required: false
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
    });
})(jQuery);