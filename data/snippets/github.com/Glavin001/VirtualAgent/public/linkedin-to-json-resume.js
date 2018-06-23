(function() {

  /* Check if LinkedIn */
  var redirectProfileUrl = 'https://www.linkedin.com/profile/view';
  if (location.hostname !== 'www.linkedin.com') {
    window.location.href = redirectProfileUrl;
    return;
  }


  var basics = {};
  basics.name = $('.profile-overview-content .full-name')[0].textContent;
  basics.label = $('#headline-container .title.field-text')[0].textContent;
  basics.picture = $('.profile-picture img')[0].getAttribute('src');
  basics.location = $('#location .locality')[0].textContent;
  var profileUrl = $('.public-profile-url')[0].textContent;
  basics.website = profileUrl;
  basics.summary = $('#summary-item .summary .field-text')[0].textContent;
  basics.profiles = [{
    network: "LinkedIn",
    "url": profileUrl
  }];

  /* Work */
  var $experience = Array.prototype.slice.call($('#background-experience-container .entity-container'));
  var work = $experience.map(function($job) {
    var job = {
      position: $('header .main-header-field .field-text', $job)[0].innerHTML,
      company: $('header .sub-header-field .field-text', $job)[0].innerHTML,
      summary: $('.body-field .field-text', $job)[0].innerHTML,
    };
    job.highlights = job.summary.split('<br>\n');
    return job;
  });

  /* Education */
  var $education = Array.prototype.slice.call($('#background-education .entity-container'));
  var education = $education.map(function($edu) {
    var edu = {
      institution: $('header .main-header-field .field-text', $edu)[0].innerHTML,
      area: $('header .sub-header-field .field-text .major', $edu)[0].innerHTML,
      studyType: $('header .sub-header-field .field-text .degree', $edu)[0].innerHTML,
      courses: []
    };
    return edu;
  });

  var resume = {
    basics: basics,
    work: work,
    education: education
  };
  var resumeStr = JSON.stringify(resume, undefined, 2);
  downloadData('linkedin-json-resume.json', resumeStr, "application/json");

  /**
    @desc Download the given data to the users computer as a file
    @param {string} name - Name of the file
    @param {string} data - Contents of the file
    @param {string} type - MIME type of the file
    */
  function downloadData(name, data, type) {
    /* Browser support */
    window.URL = window.URL || window.webkitURL;
    /* Arg defaults */
    type = type || "text/plain";
    name = name || "download";
    data = data || "";
    /* Create Blob */
    var blob;
    if (data instanceof Blob) {
      blob = data;
    } else {
      blob = new Blob([data], {
        type: type
      });
    }
    var url = window.URL.createObjectURL(blob);
    /* Create link */
    var link = document.createElement("a");
    link.download = name;
    link.href = url;
      /* Download! */
      /* See http://stackoverflow.com/a/25047811/2578205 for more details */
    var event = document.createEvent("MouseEvents");
    event.initMouseEvent("click", true, false, window, 0, 0, 0, 0, 0, false, false, false, false, 0, null);
    link.dispatchEvent(event);
  }


})();