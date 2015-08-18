angular.module('Chatbox.directive.message', [])
.directive('message', function(){
  return {
    scope: {
      message: '=data'
    },
    templateUrl: 'static/templates/message.html',
    controller: function($scope){
      $scope.show_name = false;
      $scope.showName = function(b){
        $scope.show_name = b;
      };
      $scope.toggleName = function(){
        $scope.show_name = !$scope.show_name;
      };
    }
  };
});
