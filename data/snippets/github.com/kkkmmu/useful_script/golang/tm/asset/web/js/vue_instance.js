var vm = new Vue({
    el: "#vue_det",
    data: {
        firstname: "Ria",
        lastname: "Signh",
        address : "Mumbai"
    },

    methods: {
        mydetails: function() {
            return "I am " + this.firstname + " " + this.lastname;
        }
    }
})
