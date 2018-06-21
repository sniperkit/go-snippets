require.config({
    baseUrl: '/public/js',
    waitSeconds: 1,
    paths: {
        'jquery': 'plugins/jquery-3.2.1.min',
        'bootstrap': 'plugins/bootstrap.min'
    },
    shim: {
        'layer': {
            deps: ['jquery'],
            exports: "layer"
        },
        'bootstrap': {
            deps: ['jquery']
        }
    }
});

require([
    'jquery', 'models/nprogress'
], function ($) {
    'use strict';
    let curPage = $('#current-page').data('page') || location.pathname.replace(/(.+?)\.html/i, '$1');
    curPage = curPage.replace(/^\//, '').replace(/$\//, '') || 'index';
    require(['models/' + curPage], () => {
    }, (err) => {
        console.log('loading ' + err.requireModules.toString() + ' failure!');
    })
});