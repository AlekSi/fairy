"use strict";

function MessageCtrl($scope, $http, $timeout) {
  $scope.messages = [];
  $scope.count = 0;

  (function tick() {
    $http.get("/topic").
      success(function(data, status, headers, config) {
        $scope.messages.push(data);
        if ($scope.count > 20) {
          $scope.messages.splice(0, 1);
        }
        $scope.count = $scope.count + 1;
        $timeout(tick, 0);
      }).
      error(function(data, status, headers, config) {
        alert(data);
      });
  })();
}
