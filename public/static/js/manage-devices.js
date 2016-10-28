(function($) {
    Vue.component("device-item", {
        template: "#device-item-template",
        delimiters: delimiters,
        props: {
            device: {
                type: Object,
                required: true
            }
        },
        computed: {
            itemHref: function() {
                return '#device-' + this.device.id;
            },
            itemID: function() {
                return 'device-' + this.device.id;
            }
        }
    });

    var vm = new Vue({
        el: "#app",
        delimiters: delimiters,
        data: {
            devices: []
        }
    });

    API.getAllDevices(function(data) {
        vm.devices = data.data;
    }, function(j, t, e) {
        console.error(e);
    });
})(jQuery);