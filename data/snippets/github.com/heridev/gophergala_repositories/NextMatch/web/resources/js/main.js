$(document).ready(function() {
    'use strict';
    var setsport = function(sport){
        d = new Date();
        $.get('/sport/'+sport, function(data) {
            $('.teams').attr('data-sport', sport);
            $('.teams').attr('data-date', d);
            $('.teams').html(data);
        });
    }

    $('#load-more').click(function(e){
        e.preventDefault();
        var lastd = new Date($('.teams').attr('data-date'));
        lastd.setDate(lastd.getDate()-1);
        var sport = $('.teams').data('sport');
        var datef = ''
                 +lastd.getFullYear()+ ''
                 +lastd.getMonth() + 1 + ''
                 +lastd.getDate();

        var uri = '/sport/'+sport.toLowerCase()+'/'+datef;
        $('.teams').attr('data-date', lastd.toString());
        $.get(uri, function(data) {
            $('.teams').append(data);
        });
    });
});