define(['jquery'], ($) => {
    let used = false;
    return function () {
        $('ul.pagination:not(.parsed_page)').each((_, elem) => {
            let data = $(elem).parent('nav').data(),
                url = location.href.replace(/page=\d+&?/, '');
            const page = parseInt(data['page']) || 1,
                size = parseInt(data['size']) || 10,
                total = parseInt(data['total']) || 0,
                maxPage = Math.ceil(total / size),
                range = 5;
            if (total > 0) {
                let pagination = [], max, i;
                i = page > range ? page - range + 1 : 1;
                max = i + range * 2;
                if (maxPage <= max) {
                    max = maxPage + 1;
                    i = max - range * 2 + 1;
                    if(i < 1){
                        i = 1;
                    }
                }
                if (url.search(/\?/) === -1) {
                    url = url + "?";
                }
                if (url.substr(url.length - 1) === '&') {
                    url = url.substring(0, url.length - 1);
                }
                if (url.substr(url.length - 1) !== '?') {
                    url += "&"
                }
                url += "page=";
                pagination.push(`
                <li class="${page <= 1 ? 'disabled' : ''}">
                    <a href="${page > 1 ? url + (page - 1) : '#'}" aria-label="上一页" data-pjax>
                        <span aria-label="true">&larr;</span>
                    </a>
                </li>
            `);
                for (; i < max; i++) {
                    let activeClass = i === page ? 'active' : '';
                    pagination.push(`
                    <li class="${activeClass}">
                        <a href="${url + i}" data-pjax>${i}</a>
                    </li>
                `);
                }
                pagination.push(`
                <li class="${page >= maxPage ? 'disabled' : ''}">
                    <a href="${page < maxPage ? url + (page + 1) : '#'}" aria-label="下一页" data-pjax>
                        <span aria-label="true">&rarr;</span>
                    </a>
                </li>
            `);
                $(elem).html(pagination.join("\n")).addClass('parse_page');
                if(!used){
                    $(elem).animateCss('fadeInDown');
                }
            }
        });
        used = true;
    };
});