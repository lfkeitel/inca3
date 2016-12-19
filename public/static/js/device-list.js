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
                    type: { id: 1 }
                }
            }
        },
        filters: filters,
        methods: {
            saveDevice: function() {
                console.log("Saving device");
                // Get type ID and convert to int
                this.device.type = { id: +$('#deviceType').val() };
                API.saveDevice(this.device, function(data) {
                    loadDeviceList();
                    changeState('');
                    toastr["success"]("Device Saved");
                }, function(resp) {
                    var json = resp.responseJSON;
                    toastr["error"](json.message);
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
                populateTypes();
                changeState('add', 'deviceAdd');
            },
        }
    });

    if (w.location.hash === '#add') {
        populateTypes();
        vm.section = "deviceAdd";
    }

    function loadDeviceList() {
        API.getAllDevices(function(data) {
            vm.tableData = data.data;
        }, function(j, t, e) {
            console.error(e);
        });
    }

    function populateTypes() {
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
        w.location.hash = '#' + hash;
        vm.section = section;
    }

    loadDeviceList();
})(jQuery, window);