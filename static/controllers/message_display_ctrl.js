app.controller('MessageDisplayCtrl', ['$scope', 'message_transport', function($scope, message_transport) {
  $scope.messages = [];
  message_transport.addMessageListener(function(msg) {
    $scope.$apply(function(){
      $scope.messages.push(msg);
    });
  });
}]);
