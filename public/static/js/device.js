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
                var oldSlug = this.device.slug;
                this.device.type = { id: +$('#deviceType').val() };
                API.saveDevice(this.device, function(data) {
                    if (data.data.slug !== oldSlug) {
                        flashes.add("success", "Device saved");
                        window.location = "/devices/" + data.data.slug;
                        return;
                    }
                    changeState('');
                    toastr["success"]("Device Saved");
                }, function(resp) {
                    var json = resp.responseJSON;
                    toastr["error"](json.message);
                });
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
            jsonKeys: ["created", "slug", "compressed", "size"],
            tableData: [],
            device: { slug: '', type: { id: 0 } },
            section: defaultSection
        },
        methods: {
            editDevice: function() {
                var typeID = this.device.type.id
                populateTypes(function() {
                    $('#deviceType').val(String(typeID));
                });
                changeState('edit', 'deviceEdit');
            },

            cancelEdit: function() {
                this.device = getOriginalDeviceData();
                changeState('');
            },

            deleteDevice: function() {
                var slug = this.device.slug;
                var confirm = new jsConfirm();
                confirm.show("Are you sure you want to delete this device?<br>This will also delete configurations for this device", function() {
                    API.deleteDevice(slug, function() {
                        flashes.add('success', "Device deleted");
                        window.location = "/devices";
                        return
                    }, function(resp) {
                        toastr["error"](resp.responseJSON.message);
                    })
                });
            }
        }
    });

    var devID = w.location.pathname.split('/');
    devID = devID[devID.length - 1];

    if (w.location.hash === '#edit') {
        populateTypes();
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

    function populateTypes(callback) {
        API.getAllTypes(function(data) {
            var types = data.data;
            var typeSelect = $('#deviceType');

            for (var i = 0; i < types.length; i++) {
                var o = types[i];
                typeSelect.append($('<option>', {
                    value: o.id,
                    text: o.name
                }));
            }

            if (callback) { callback(); }
        }, function(j, t, e) {
            console.log(e);
            console.log(j.responseJSON);
        })
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