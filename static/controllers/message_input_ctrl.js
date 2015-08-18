app.controller('MessageInputCtrl', ['$scope', 'message_transport', function($scope, message_transport) {
  $scope.message = '';
  $scope.submit = function() {
    message_transport.send($scope.message);
    $scope.message = '';
  }
}]);
