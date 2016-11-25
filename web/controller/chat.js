var app = angular.module("chat", []);
var rooom = angulaar.module("chat-room", ['ngRoute']).config(
    function($routeProvider, $locationProvider) {
        // configure the routing rules here
        $routeProvider.when('/room/:room', {
            controller: 'PagesCtrl'
        });
        // enable HTML5mode to disable hashbang urls
        $locationProvider.html5Mode(true);
    }).controller('PagesCtrl', function ($routeParams) {
        console.log($room);
    });


genericcontroller = function(isRoom) {
    // we need to get the scope from the html this'll be applied to
    var scope = angular.element(document.getElementById("scope-wrap")).scope();
    scope.message = []; // list of messages

    var host = location.host;
    // check if we're on a secure connection, if so
    // we'll use a secure websocket, if not, we'll
    // use a regular one
    if(location.protocol === "https:") {
        var ws = "wss://";
    } else {
        var ws = "ws://";
    }

    var addr = ws.concat(host); // angular for getting url + port
    if (isRoom) {
        var connection = new WebSocket(addr.concat("/room/:param1"))
    } else {
        var connection = new WebSocket(addr.concat("/chat"));
    }

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
}

app.controller("ChatController", ["$scope", ])
