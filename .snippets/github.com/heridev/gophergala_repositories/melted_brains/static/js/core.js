window.preventType = function(){
  window.onkeypress = function(){
    return false;
  }
}

window.allowType = function(){
  window.onkeypress = window.GO_RACER_NS.onKeyPress;
}
window.countDown = function(){
  $('.countdown .counter h1').html("3");
  $('.countdown').show();
  setTimeout(function(){
    $('.countdown .counter h1').html("2");
  }, 1000)
  setTimeout(function(){
    $('.countdown .counter h1').html("1");
  }, 2000)
  setTimeout(function(){
    $('.countdown .counter h1').html("GO!");
    allowType();
  }, 3000)
  setTimeout(function(){
    $('.countdown').hide();
  }, 4000)
}
window.GO_RACER_NS = {}

window.GO_RACER_NS['parseTokens'] = function (elements){
  // skip comment & new lines
  var collected = [];
  for(var i=0; i<elements.length; i++){
    if(["com"].indexOf(elements[i].className)==-1 && elements[i].innerHTML.replace(/\n|\t/g, '') != ''){
      elements[i].classList.add('arc');
      collected.push(elements[i]);
    }
  }
  return collected;
}
function showWinners(){
  var board = $('.board ol');
  for(var u in window.winner_list){
    var div = $("<li></li>").html(window.users_ids[u]);
    board.append(div);
  }
  $('.board').css({left: $('pre').width() + 100, top: 100});
  $('.board').show();
}
function MacOsXCharCode(charCode){
  if(charCode == 13){
    return 10;
  }
  return charCode;
}
function numberOFDupInStr(str){
  var res = str.match(/(\t{1,})/g);

  if(str[0] == "\t" && res){
    return res[0].length;
  }
  else {
    return 0;
  }
}

window.GO_RACER_NS['carretMoveLogic'] = function(charCode, user){
  var newOffset = null,
      isNextTab,
      i;
  // debugger;
  var tokenElem = window.users[user].tokenElem;

  if(MacOsXCharCode(charCode) == tokenElem.value.charCodeAt(tokenElem.offSet)){
    isNextTabs = numberOFDupInStr(tokenElem.value.substr(tokenElem.offSet+1));
    newOffset = tokenElem.offSet++;

    if(isNextTabs > 0){
      tokenElem.offSet = tokenElem.offSet + isNextTabs;
      newOffset = tokenElem.offSet - 1;
    }

    tokenElem.el.innerHTML = setCarretAndSkipTab(tokenElem.value, newOffset);
    if(user != window.current_user_id){
      var p = $(tokenElem.el).find('.carret').offset();
      window.users[user].carret.offset({top: p.top, left: p.left}).show();
      window.users[user].flag.offset({top: p.top-13, left: p.left}).show();
    }

    // next span element
    if( (newOffset+1) >= tokenElem.value.length || tokenElem.value == "&amp;"){
      tokenElem.el.innerHTML = tokenElem.value;
      if(user == window.current_user_id){
        tokenElem.el.classList.remove('arc');
      }
      if(tokenElem.index+1 == window.users[user].tokens.length){
        window.winner_list.push(user);
        if(window.winner_list.length == Object.keys(window.users_ids).length){
          showWinners();
        }
        return false;
      }
      window.users[user].tokenElem = window.GO_RACER_NS.chooseNextAt(user, tokenElem.index+1);
      var str = window.users[user].tokenElem.value;
      window.users[user].tokenElem.el.innerHTML = setCarret(str);
    }
  }
}

window.GO_RACER_NS['onKeyPress'] = function (pEvent){
  window.ws.send(pEvent.charCode);

  if(pEvent.keyCode == 32 && pEvent.target == document.body) {
    pEvent.preventDefault();
    return false;
  }
}

function setCarretAndSkipTab(val, offSet){
  // case01
  if(val == "&amp;"){
    return (val + window.carret);
  }
  return val.substr(0, offSet+1) + window.carret + val.substr(offSet+1);
}
function setCarret(str){
  return window.carret + str;
}
window.GO_RACER_NS['chooseNextAt'] = function(user, idx){
  return {
    index: idx,
    el: window.users[user].tokens[idx],
    value: window.users[user].tokens[idx].innerHTML,
    offSet: 0
  }
}

window.GO_RACER_NS['prepareGameField'] = function(){
  var codeBlock;
  // we need global
  window.users = {};
  window.winner_list = [];
  window.carret = "<span class='carret blink'></span>";
  window.colors = ['#FFA500', '#FF4500', '#DA70D6', '#DB7093'];

  var tokens = null;
  for (var userName in window.users_ids) {
    console.log('visit', userName);
    if (window.users_ids.hasOwnProperty(userName)) {
      codeBlock = $('#code_'+userName + ' span');
      tokens = window.GO_RACER_NS.parseTokens(codeBlock);

      window.users[userName] = {
        tokens: tokens,
        status: 'on_start',
        tokenElem: null
      }
      if(userName != window.current_user_id){
        var carret = $("<span></span>");
        carret.addClass('carret blink');
        carret.offset({top: 0 + 100, left: 0}).hide();
        $(document.body).append(carret);

        var flag = $("<span></span>");

        flag.addClass('flag');
        flag.css('background-color', window.colors.pop());
        flag.html(window.users_ids[userName]).hide();
        $(document.body).append(flag);
        window.users[userName]['carret'] = carret;
        window.users[userName]['flag'] = flag;
      }
      window.users[userName]['tokenElem'] = window.GO_RACER_NS.chooseNextAt(userName, 0);
    }
  }
  // window.tokenElem = window.GO_RACER_NS.chooseNextAt(0);

}

window.GO_RACER_NS['user1'] = function(){
  var chaCode;
  setInterval(function(){
  // for(i=0; i<100; i++){
    charCode = window.users['user1'].tokenElem.value.charCodeAt(window.users['user1'].tokenElem.offSet);
    window.GO_RACER_NS.carretMoveLogic(charCode, 'user1');
  // }
  }, 500)
}
window.GO_RACER_NS['user2'] = function(){
  var chaCode;
  for(i=0; i<100; i++){
    setTimeout(function(){
      charCode = window.users['user2'].tokenElem.value.charCodeAt(window.users['user2'].tokenElem.offSet);
      console.log('user2 typed', charCode);
      window.GO_RACER_NS.carretMoveLogic(charCode, 'user2');
    }, 5000)
  }
}
