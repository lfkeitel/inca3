(function($, w) {
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

    Vue.component("device-edit-form", {
        template: "#edit-form-template",
        delimiters: delimiters,
        data: function() {
            return {
                device: {
                    id: 0,
                    name: '',
                    address: '',
                    brand: '',
                    connection: ''
                }
            }
        },
        filters: filters,
        methods: {
            saveDevice: function() {
                console.log("Saving device");
                API.saveDevice(this.device, function(data) {
                    loadDeviceList();
                    changeState('');
                });
            },
            cancel: function() {
                changeState('');
            }
        }
    });

    var defaultSection = "devList";
    var vm = new Vue({
        el: "#app",
        delimiters: delimiters,
        data: {
            tableColumns: ["name", "address", "connection", "brand"],
            tableData: [],
            section: defaultSection
        },
        methods: {
            addDevice: function() {
                changeState('add', 'deviceAdd');
            },
        }
    });

    if (w.location.hash === '#add') {
        vm.section = "deviceAdd";
    }

    function loadDeviceList() {
        API.getAllDevices(function(data) {
            vm.tableData = data.data;
        }, function(j, t, e) {
            console.error(e);
        });
    }

    function changeState(hash, section) {
        if (hash === '' || vm.section === section) {
            w.location.hash = '';
            vm.section = defaultSection;
            return;
        }
        w.location.hash = hash;
        vm.section = section;
    }

    loadDeviceList();
})(jQuery, window);