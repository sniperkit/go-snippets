$(function(){
    $('.dir').dblclick(function(){
        var path = $(this).data('path');
        $.post("/", {path: path}).done(function(data){
            $('body').html(data)
        })
    })
    
    $('.dir').click(function(){
        var path = $(this).data('path');
        $('div#pathfield').html(path);
    })
    
    $('#sortbutton').click(function(){
        var path = $('#pathfield').text();
//        console.log(path)
        $.post("/sort", {path: path}).done(function(data){
            $('body').html(data)
        })
    })
})