/**
 * Created by wuhailin on 2017/11/24.
 */
require.config({
    //urlArgs: 'ver=1.0_' + (new Date()).getTime(),
    //指定本地模块位置的基准目, 如果引入require.js的脚本里指定了data-main,没配置baseUrl时，会自动使用data-main的前缀路径
    //<script defer async src="js/require.js" data-main="js/main" id="current-page" data-page="index"></script> baseUrl=js
    baseUrl: '/static/js',
    waitSeconds: 1,//等待时间，超时指定时间则请求下一个
    paths: {
        jquery: ['jquery-3.2.1.min', 'https://cdn.bootcss.com/jquery/3.2.1/jquery.min.js'],//如果本地文件没有加载成功，则加载远程文件
        bootstrap: 'bootstrap.min',
        underscore: ['underscore-min', 'http://www.bootcss.com/p/underscore/underscore-min.js'],
        list: ['plugins/list.min']
    },
    shim: {
        // bootstrap: {
        //     deps: ['jquery']
        // },
        bootstrap: ['jquery'],
        underscore: {
            exports: '_'
        },
        jquery: {
            exports: '$' //把jQuery的变量值传到回调区域
        },
        list: {
            exports: 'List'
        }
    }
});

//初始化加载文件，多页加载
require(['jquery', 'underscore', 'models/bootstrap', 'models/common/animate'], ($, _) => {
    let curPage = $('#current-page').attr('data-page') || location.pathname.replace(/(.+)\.html\??.*/, '$1');
    curPage = curPage.replace(/^\//, '').replace(/\/$/, '') || 'index';
    if (curPage) {
        require(['models/' + curPage], () => {
        }, (err) => {
            const failedId = err.requireModules;
            if (failedId) {
                _.each(failedId, require.undef);//文件加载失败时，依赖该文件的取消加载
                console.warn('loading ' + failedId.toString() + ' error');
            }
        });
    }
});
