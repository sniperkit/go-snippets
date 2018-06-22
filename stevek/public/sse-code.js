stream = new EventSource("/events");
stream.onopen = function() {
    // plan to send start trigger
    //$.post("/users", {
    //    user: username
    //});
};
stream.onmessage = function(e) {
    payload = JSON.parse(e.data);
    //if (payload.type === "message") {
    //    addMessage(payload.data);
    //}
    //if (payload.type === "users") {
    //    updateUsers(payload.data);
    //}
    el = $("<p>").html(JSON.stringify(payload.data));
    $('#data').append(el);
    push(payload.data);
};
window.onbeforeunload = function() {
    //$.ajax({
    //    url: "/users?user=" + username,
    //    type: "DELETE",
    //});
    stream.close();
};
