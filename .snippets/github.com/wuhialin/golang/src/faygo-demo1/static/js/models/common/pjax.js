require(['jquery', 'models/common/nprogress', 'models/common/pagination', 'plugins/jquery.pjax'], ($, NProgress, pagination) => {
    if ($('#pjax-container').length) {
        $(document).pjax('a[data-pjax]', '#pjax-container');
        $(document).on('submit', 'form[data-pjax]', function (e) {
            $.pjax.submit(e, '#pjax-container');
        });
        $(document).on('pjax:start', function () {
            NProgress.start();
        });
        $(document).on('pjax:end', function () {
            NProgress.done();
            pagination();
        });
    }
});