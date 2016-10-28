(function($, w) {
    Vue.component("config-list", {
        template: "#config-list-template",
        delimiters: delimiters,
        props: {
            tableData: {
                type: Array,
                required: false
            },
            tableCols: {
                type: Array,
                required: true
            },
            jsonKeys: {
                type: Array,
                required: true
            },
            deviceName: {
                type: String,
                required: true
            }
        },
        filters: filters,
        methods: {
            gotoConfig: function(id) {
                window.location = "/devices/" + this.deviceName + "/" + id;
            }
        }
    });

    Vue.component("device-edit-form", {
        template: "#edit-form-template",
        delimiters: delimiters,
        props: {
            device: {
                type: Object,
                required: true
            }
        },
        filters: filters,
        methods: {
            saveDevice: function() {
                console.log("Saving device");
                changeState('');
            },
            cancelEdit: function() {
                this.$emit('cancel-edit');
            }
        }
    });

    var defaultSection = "configs";
    var vm = new Vue({
        el: "#app",
        delimiters: delimiters,
        data: {
            tableColumns: ["date", "name", "compressed", "size"],
            jsonKeys: ["created", "id", "compressed", "size"],
            tableData: [],
            device: { "id": '' },
            section: defaultSection
        },
        methods: {
            editDevice: function() {
                changeState('edit', 'deviceEdit');
            },
            cancelEdit: function() {
                this.device = getOriginalDeviceData();
                changeState('');
            }
        }
    });

    var devID = w.location.pathname.split('/');
    devID = devID[devID.length - 1];

    if (w.location.hash === '#edit') {
        vm.section = "deviceEdit";
    }

    var originalDevice = "";

    API.getDevice(devID, function(data) {
        originalDevice = JSON.stringify(data.data);
        vm.device = data.data;
        getDeviceConfigs();
    }, function(j, t, e) {
        console.log(e);
    });

    function getDeviceConfigs() {
        API.getDeviceConfigs(devID, function(data) {
            vm.tableData = data.data;
        }, function(j, t, e) {
            console.error(e);
        });
    }

    function getOriginalDeviceData() {
        return JSON.parse(originalDevice);
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
})(jQuery, window);