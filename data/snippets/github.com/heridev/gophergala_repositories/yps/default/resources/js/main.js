$(document).ready(function () {
    $("#addToQueue").click(function (e) {
        e.preventDefault();
        $.ajax({
            type: "POST",
            url: "/addToQueue",
            data: {
                url: $("#videoLink").val()
                //access_token: $("#access_token").val()
            },
            success: function (result) {
                console.log(result)
                //$("#sharelink").html(result);
            }
        });
    });
});
