/*
* @Author: wuhailin
* @Date:   2017-11-15 11:33:30
* @Last Modified by:   wuhailin
* @Last Modified time: 2017-11-15 11:39:26
*/
define(['jquery', 'models/nprogress', 'plugins/pjax'], ($, NProgress) => {
    let pjaxContainer = '#pjax-container',
        app = {};
    if($(pjaxContainer).length){
        require(['plugins/pjax'], () => {
            $document.pjax('[data-pjax] a, a[data-pjax]', pjaxContainer)
            .on('submit', 'form[data-pjax]', e => $.pjax.submit(e, pjaxContainer))
            .on('pjax:send', () => NProgress.start());
        });
    }
});