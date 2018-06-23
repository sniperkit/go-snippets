(function(){
    var questions = [];
    var currentQuestion = -1;

    var btnStart;
    var btnSubmit;
    var btnNext;
    var btnRestart;
    var introDiv;
    var questionDiv;
    var finishDiv;
    var errorDiv;
    var questionBlock;
    var snippetBlock;
    var answerBlock;
    var statusBlock;
    var timerBlock;
    var scoreBlock;
    var helpLink;
    var aboutLink;
    var boardLink;
    var twitterLink;
    var facebookLink;
    var score = 0;
    var time;

    var timerInterval;

    var answerBlockInit = "//Write your codes here";

    function init(){
      btnStart = $('#btnStartGame');
      btnStart.click(function(){
        startTimer();
        nextQuestion();
      });

      btnSubmit = $('#btnSubmitAnswer');
      btnSubmit.click(function(){
        submitQuestion();
      });
      btnNext = $('#btnNextQuestion');
      btnNext.click(function(){
        nextQuestion();
      });

      btnRestart = $('#btnRestartGame');
      btnRestart.click(function(){
        startTimer();
        currentQuestion = -1;
        score = 0;
        nextQuestion();
      });

      helpLink = $('#helpLink');
      helpLink.click(function(){
        showModal("How to Play", helpHtml);
      });
      aboutLink = $('#aboutLink');
      aboutLink.click(function(){
        showModal("About Gopher Quiz", aboutHtml);
      });
      boardLink = $("#boardLink");
      boardLink.click(function(){
        showHighScores();
      });
      twitterLink = $('#twitterLink');
      twitterLink.click(function(evt){
        evt.preventDefault();
        var text = encodeURIComponent("I scored "+score+" in Gopher Quiz! Beat that. http://gopherquiz.com #golang #gopherquiz");
        var href = "http://twitter.com/intent/tweet?text="+text;
        window.open(href);
      })
      facebookLink = $("#facebookLink");
      facebookLink.click(function(evt){
        evt.preventDefault();
        var href = "http://www.facebook.com/sharer.php?s=100&p[title]=Gopher Quizt&p[summary]=I scored "+score+" in Gopher Quiz! Beat that.&p[url]=http://gopherquiz.com&p[images][0]=http://gopherquiz.com/static/images/gophergame.png";
        window.open(href);
      });

      introDiv = $("#introdiv");
      questionDiv = $('#questiondiv');
      finishDiv = $('#finishdiv');
      errorDiv = $('#errordiv');

      errorDiv.find(".btn").click(function(){
        window.location.reload();
      });

      snippetBlock = $("#snippetBlock");
      questionBlock = $("#questionBlock");
      timerBlock = $("#timerBlock");
      scoreBlock = $("#scoreBlock");

      answerBlock = $("#answerBlock");
      answerBlock.bind('paste', function(e){
        e.preventDefault();
      });
      answerBlock.keydown(function(e) {
          if(e.keyCode === 9) { // tab was pressed
              // get caret position/selection
              var start = this.selectionStart;
              var end = this.selectionEnd;

              var $this = $(this);
              var value = $this.val();

              // set textarea value to: text before caret + tab + text after caret
              $this.val(value.substring(0, start)
                          + "\t"
                          + value.substring(end));

              // put caret at right position again (add one for the tab)
              this.selectionStart = this.selectionEnd = start + 1;

              // prevent the focus lose
              e.preventDefault();
          }
      });

      statusBlock = $('#statusBlock');
    }

    $(document).ready(function(){
        init();
        $.get("/questions", function(data){
            if (data.status !== "success"){
                showInitializeError();
                return
            }
            questions = data.data;
        }).error(function(err){
            showInitializeError();
        });
    });

    function reset(){
        currentQuestion = -1;
        score = 0;
        introDiv.show();
        questionDiv.hide();
        finishDiv.hide();
        stopTimer();
    }

    function nextQuestion(){
        currentQuestion++;

        questionBlock.html(questions[currentQuestion].question)
        snippetBlock.html(questions[currentQuestion].snippet)
        answerBlock.val(answerBlockInit);

        introDiv.hide();
        questionDiv.show();
        answerBlock.focus();
        statusBlock.hide();
        finishDiv.hide();
        btnSubmit.show();
        btnNext.hide();
    }

    function submitQuestion(){
        var answer = answerBlock.val();
        setStatus("Processing...");
        $.post(
            '/submitQuestion',
            {answer: answer, id: questions[currentQuestion].id},
            function(data){
                if (data.status !== "success"){
                    var msg = data.message || "An error occured, please try again."
                    if (data.output){
                        msg += "\n\nOutput:\n" + data.output
                    }
                     setStatus(msg, true);
                     return;
                }
                var output = data.data.output || "success"
                output = data.data.message + "\n\nOutput:\n" + output;
                setStatus(output)
                btnSubmit.hide();
                btnNext.show();
                score += 10;
                if (currentQuestion >= questions.length - 1){
                    finish();
                }
        }).error(function (err){
            setStatus("Error occurred while submitting question, please try again.", true);
        });
    }

    function finish(){
        var tagLine = "Game over! Well done."
        if (time > 0){
            tagLine = "You rock! You finished before time runs out and earned extra "+time+" score.";
            score += time;
        }
        submitScore(score);
        stopTimer();
        finishDiv.find(".finaltagline").html(tagLine);
        finishDiv.find(".finalscore").html(score);

        introDiv.hide();
        questionDiv.hide();
        finishDiv.show();
    }

    function setStatus(msg, err){
        statusBlock.removeClass("error");
        statusBlock.html(msg);
        if (err) {
            statusBlock.addClass("error")
        }
        statusBlock.show();
    }

    function showInitializeError(){
        introDiv.hide();
        questionDiv.hide();
        finishDiv.hide();
        errorDiv.show();
    }

    function startTimer(){
        time =  60 * 5;
        timerInterval = setInterval(function(){
            time--;
            showTimeAndScore();
            if (time <= 0){
                stopTimer();
                finish();
            }
        }, 1000);
    }

    function stopTimer(){
        clearInterval(timerInterval);
    }

    function showTimeAndScore(){
        var m = Math.floor(time / 60);
        var s = time % 60;
        var t = m + "m " + s + "s";
        timerBlock.html(t);
        scoreBlock.html('<span style="font-size: 20px;">Score:</span> '+score);
    }

    function showModal(title, content){
        var div = $('#modaldiv');
        div.find(".modal-title").html(title);
        div.find(".modal-body").html(content)
    }

    function submitScore(score){
        $.post("/submitScore", {score: score}, function(data){});
    }

    function showHighScores(){
        showModal("LeaderBoard", "<p style='text-align: center;'>Loading...</p>");
        $.get("/highscores", function(data){
            if (data.status !== "success"){
                showModal('LeaderBoard', '<span class="error">Error occured while loading LeaderBoard</span>');
                return;
            }
            var contents = '<table class="table table-hover"><thead><td>Pos</td><td>Username</td><td>Highscore</td></thead><tbody>';
            for(var i in data.data){
              var user = data.data[i]
              if (user.Highscore > 0 ){
                 contents += "<tr><td>"+((i/1)+1)+"</td><td>"+user.Username+"</td><td>"+user.Highscore+"</td></tr>";
              }
            }
            contents += "</tbody></table>";
            showModal("LeaderBoard", contents);
        }).error(function(err){
            showModal('LeaderBoard', '<span class="error">Error occured while loading LeaderBoard</span>');
        });
    }

    var aboutHtml = '<div class="text-center"><img src="/static/images/gophergame.png"  style="height: 100px;"/><p>Gopher Quiz is a simple programming quiz for <a href="http://golang.org" target="_blank">Go programming language</a>.</p><p>Feedback are welcomed at  <a href="mailto:abiola89@gmail.com" target="_blank">abiola89@gmail.com</a></p></div>';
    var helpHtml = '<p>You are required to input the snippet that will complement the source code to make the program work as stated in the questions.</p><p>Every correct answer gives a score of 10 and finishing ahead of time gives an additional score of number of seconds remaining.</p><p>Remember, you have 5 minutes to show the SuperGopher you really are.</p><p>Have fun.</p>';
})()