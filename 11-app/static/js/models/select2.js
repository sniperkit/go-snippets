/**
 * Created by wuhailin on 2017/11/24.
 */
define(['jquery'], ($, css) => {
    'use strict';
    const $dom = $('select.select2');
    if ($dom.length) {
        require(['text!select2/css/select2.min.css', 'select2/js/select2.full'], (css) => {
            $('head').append('<style id="select2-css">' + css + '</style>');
            $dom.select2();
        });
    }
});