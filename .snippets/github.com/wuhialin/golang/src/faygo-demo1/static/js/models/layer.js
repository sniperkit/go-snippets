/**
 * Created by wuhailin on 2017/11/24.
 */
define(['jquery', 'layer/layer', 'text!layer/theme/default/layer.css'], ($, layer, css) => {
    'use strict';
    css = css.replace(/url\(\s*(.+?)\s*\)/ig, 'url(static/js/layer/theme/default/$1)');
    $('head').append('<style id="layer-css">' + css + '</style>');
    return layer;
});