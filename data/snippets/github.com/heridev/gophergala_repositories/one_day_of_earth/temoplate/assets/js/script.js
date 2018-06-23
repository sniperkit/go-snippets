$(function(){

    var $container = $('#posts');
    $container.isotope({
        itemSelector: '.post_item',
        layoutMode: 'fitRows'
    });

});


function refresh_layout() {
    var $container = $('#posts');
    $container.isotope( 'reLayout', function(){
        alert(1);
    } );
}