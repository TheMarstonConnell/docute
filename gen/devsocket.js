window.addEventListener("load", start)


function start() {
    console.log("hello devs!")

    const u = "ws://" + document.location.host + "/ws"
    console.log(u)

    const ws = new WebSocket(u)
    ws.onopen = function(event) {
        console.log("WebSocket is open now.");
    };

    ws.onmessage = function(event) {
        if (event.data == "refresh") {
            location.reload()
        }
    };

    ws.onclose = function(event) {
        console.log("WebSocket is closed now.");
    };

    ws.onerror = function(error) {
        console.log("WebSocket error observed:", error);
    };


}