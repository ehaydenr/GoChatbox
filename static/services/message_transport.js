app.factory('message_transport', function() {

  chat.onchat = function(msg) {
    _publishMessage(msg);
  }

  var messageListeners = [];

  function _send(content) {
    if(!chat.isReady()){
      setTimeout(function(){
        _send(content);
      }, 300);
      return;
    }

    chat.send({
      content: content,
      operation: 'chat'
    });
  }

  function _addMessageListener(fn) {
    messageListeners.push(fn);
  }

  function _publishMessage(msg) {
    for(var i = 0; i < messageListeners.length; ++i){
      messageListeners[i](msg);
    }

    var display = document.querySelector('.message-display');
    display.scrollTop = display.scrollHeight;
  }

  return {
    send: _send,
    addMessageListener: _addMessageListener,
    publishMessage: _publishMessage
  };
});
