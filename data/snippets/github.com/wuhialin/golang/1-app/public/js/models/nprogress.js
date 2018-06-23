/*
* @Author: wuhailin
* @Date:   2017-11-15 09:55:17
* @Last Modified by:   wuhailin
* @Last Modified time: 2017-11-15 11:32:33
*/
define([
    'jquery', 'plugins/nprogress',
    'text!/public/css/nprogress.css'
], ($, NProgress, css) => {
    let $body = $('body'),
        cssHtml = [];
    cssHtml.push('<style>');
    cssHtml.push(css);
    cssHtml.push('</style');
    $body.append(cssHtml.join("\n"));

    let $document = $(document),
        app = {
            ajaxStart: function () {
                NProgress.start();
            },
            ajaxEnd: function () {
                NProgress.done();
            }
        };

    $document.ajaxStart(app.ajaxStart).ajaxStop(app.ajaxEnd);
    return NProgress;
});