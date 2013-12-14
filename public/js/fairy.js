"use strict";

function MessageCtrl($scope, $http, $timeout) {
  $scope.messages = [];

  (function tick() {
    $http.get("/topic").
      success(function(data, status, headers, config) {
        $scope.messages.push({body: data});
        $timeout(tick, 0);
      }).
      error(function(data, status, headers, config) {
        alert(data);
      });
  })();
}
