define([
    'jquery',
    'text!plugins/nprogress/nprogress.css',
    'plugins/nprogress/nprogress'
], ($, css, NProgress) => {
    'use strict';
    $('head').append('<style id="nprogress-css">' + css + '</style>');
    return NProgress;
});