/**
 * Created by wuhailin on 2017/11/24.
 */
define(['jquery', 'models/common/app', 'bootstrap'], ($, app) => {
    $('[data-toggle=tooltip]').tooltip();

    if ($('form').length) {
        require(['models/common/form'], () => {
        });
    }

    if($('ul.pagination').length){
        require(['models/common/pagination'], function(pagination){
            pagination();
        })
    }

    if($('#pjax-container').length){
        require(['models/common/pjax']);
    }

    $('#main-nav').find('a').each((k, v) => {
        const href = $.trim($(v).attr('href'));
        if(app.resolve(href) === app.resolve(location.pathname)){
            $(v).parents('li').addClass('active');
        }
    });
});