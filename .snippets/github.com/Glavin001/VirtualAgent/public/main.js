function showTab(tabId) {
  $('.tab-content .tab-pane').removeClass('active in');
  $('.tab-content .tab-pane#' + tabId).addClass('active in');
};

$(document).ready(function() {

  console.log('Simply Sift!');

  let data = window.agent_data = {

  };

  const tabPanes = {
    github: $('.tab-pane#github'),
    resume: $('.tab-pane#json-resume'),
  };


  // Update progress bar
  socket.on('status', function(status) {
    console.log('status', status);
    $('.status-message').text(status.message);
    $('.progress-bar').width(status.progress * 100 + '%');
  });

  $('.sync-github-btn').click(function(event) {
    // Show tab
    showTab('syncing');
    let username = $('#github-username').val().trim();

    socket.emit('sync-github-user', username, function(error, results) {
      console.log(error, results);
      var projects = results.repos;
      // Save projects
      data.projects = projects;

      var langMap = projectsByLanguage(projects);
      data.langMap = langMap;
      console.log('langMap', langMap);
      var langScores = _.mapValues(langMap, function(v, k) {
        return _.sumBy(v, 'languages.' + k);
      });
      data.langScores = langScores;
      console.log('langScores', langScores);

      // Generate chart of language stats
      // var chart = c3.generate({
      //   bindto: '#language-chart',
      //   data: {
      //     columns: [
      //       ['data1', 30, 200, 100, 400, 150, 250],
      //       ['data2', 130, 100, 140, 200, 150, 50]
      //     ],
      //     type: 'bar'
      //   },
      //   bar: {
      //     width: {
      //       ratio: 0.5 // this makes bar width 50% of length between ticks
      //     }
      //     // or
      //     //width: 100 // this makes bar width 100px
      //   }
      // });
      // setTimeout(function() {
      //   chart.load({
      //     columns: [
      //       ['data3', 130, -150, 200, 300, -200, 100]
      //     ]
      //   });
      // }, 1000);

      showTab('json-resume');
    });

  });

  $('input#github-username').keyup(_.debounce(function(event) {
    tabPanes.github.removeClass('invalid-username valid-username').addClass('checking-username');
    let username = $(this).val().trim();
    // console.log(username);
    socket.emit('check-github-username', username, function(error, user) {
      console.log(error, user);
      if (error) {
        tabPanes.github.removeClass('checking-username valid-username').addClass('invalid-username');
      } else {
        $('.github-user-full-name').text(user.name);
        $('img.github-user-avatar').attr('src', user.avatar_url);
        $('a.github-user-url').attr('href', user.html_url);
        tabPanes.github.removeClass('checking-username invalid-username').addClass('valid-username');
      }
    });
  },100));

  $('#import-json-resume-file').change((event) => {
    let $el = $(event.currentTarget);
    // Source: http://stackoverflow.com/a/13747921/2578205
    let file = $el.prop('files')[0];
    let reader = new FileReader();
    reader.onload = () => {
      let data = reader.result;
      let json = JSON.parse(data);
      console.log('JSON Resume', json);
      data.json_resume = json;
      tabPanes.resume.addClass('has-resume');
    };
    reader.readAsText(file);
  });

  const $jobSearchInput = $('#job-search-query');
  $('.use-resume-btn').click((event) => {
    console.log(data);
    let keywords = Object.keys(data.langMap || {});
    let query = keywords.join(', ');
    $jobSearchInput.val(query);
    $jobSearchInput.keyup();
    console.log(query);
  });

  $jobSearchInput.keyup(_.debounce(function(event) {
    let query = $jobSearchInput.val();
    socket.emit('job-search', query, function(error, results) {
      console.log('job-search', error, results);
      let source = $('#jobs-template').html();
      let template = Handlebars.compile(source);
      let context = {
        jobs: results.hits
      };
      let html = template(context);
      // console.log('html', html);
      $('.job-results').html(html);
    });
  },100));


  function projectsByLanguage(projects) {
    var langMap = {};
    for (var project of projects) {
      var langs = Object.keys(project.languages);
      for (var lang of langs) {
        if (langMap.hasOwnProperty(lang)) {
          // Has Language already
          // Add project
          langMap[lang].push(project)
        } else {
          langMap[lang] = [project]
        }
      }
    }
    return langMap;
  }


});