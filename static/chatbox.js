var chat = (function(){
  var publicAPI = {};
  var token = null;

  // Initiate WS connection
  var socket = new WebSocket("ws://localhost:8080/ws");

  publicAPI.isReady = function _isReady(){
    return socket.readyState == 1 && token != null;
  }

  socket.onmessage = function(msg){
    msg = JSON.parse(msg.data);
    if(msg.Operation == 'token_verified'){
      if(msg.Content == 'true'){
        token = msg.Content;
      }else{
        console.err("ERROR: Token not verified");
      }
    }else if(msg.Operation == 'chat'){
      if(publicAPI.onchat !== undefined && publicAPI.onchat.constructor === Function){
        publicAPI.onchat({
          Content: msg.Content,
          Origin: msg.Origin
        });
      }
    }
  };

  publicAPI.send = function _send(msg) {
    socket.send(JSON.stringify(msg));
  }

  // Give token to server
  publicAPI.giveToken = function _giveToken(tok){

    // Make sure socket ready
    if(socket.readyState != 1){
      setTimeout(function(){
        giveToken(tok);
      }, 100);
      return;
    }

    // Send token to server
    publicAPI.send({
      content: token_id,
      operation: 'giveToken'
    });
  }

  return publicAPI;
})();

var app = angular.module(
  'Chatbox',
  [
    'ngMaterial',
    'Chatbox.directive.message'
  ]
);
