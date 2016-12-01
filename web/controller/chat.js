var app = angular.module("chat", []);

app.controller("ChatController", ["$scope", function() {
    // we need to get the scope from the html this'll be applied to
    var scope = angular.element(document.getElementById("scope-wrap")).scope();
    scope.message = []; // list of messages

    var host = location.host;
    // check if we're on a secure connection, if so
    // we'll use a secure websocket, if not, we'll
    // use a regular one
    if (host.indexOf("heroku") !== -1) {
        var ws = "wss://";
    } else {
        var ws = "ws://";
    }

    var addr = ws.concat(host); // angular for getting url + port
    var connection = new WebSocket(addr.concat("/chat"));

    connection.onclose = function(event) {
        scope.$apply(function() {
            scope.message.push("SYSTEM")
            scope.message.push("DESCONECTADO")
        })
    }

    connection.onopen = function(event) {
        scope.$apply(function() {
            scope.message.push("SYSTEM")
            scope.message.push("CONECTADO")
        })
    }

    connection.onmessage = function(event) {
        scope.$apply(function() {
            scope.message.push(event.data);
        })
    }

    scope.send = function() {
        connection.send([
            scope.name,
            scope.msg
        ]);
        // scope.name = "";
        scope.msg = "";
    }
}])
