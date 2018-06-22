$(function () {
    /**
     * Mark menu elements from the navbar as active when clicked.
     */
    $('ul.nav > li > a[href="' + document.location.pathname + '"]').parent().addClass('active');
});
