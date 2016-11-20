var app = angular.module("chat", []);

app.controller("ChatController", ["$scope", function() {
    // we need to get the scope from the html this'll be applied to
    var scope = angular.element(document.getElementById("scope-wrap")).scope();
    scope.message = []; // list of messages
    var host = "ws://".concat(location.host); // angular for getting url + port
    var connection = new WebSocket(host.concat("/chat"));

    connection.onclose = function(e) {
        scope.$apply(function() {
            scope.message.push("DESCONECTADO");
        })
    }

    connection.onopen = function(e) {
        scope.$apply(function() {
            scope.message.push("CONECTADO");
        })
    }

    connection.onmessage = function(e) {
        scope.$apply(function() {
            scope.message.push(e.data);
        })
    }

    scope.send = function() {
        connection.send(scope.msg);
        scope.msg = '';
    }
}])
