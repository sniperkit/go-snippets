require([
    'jquery',
    'underscore',
], ($) => {
    'use strict';
    const $body = $('body');
    $body.on('submit', 'form', e => {
        const $this = $(e.currentTarget),
            attr = 'skip-trim';
        if (!$this.attr(attr)) {
            e.preventDefault();
            $this.find(':enabled:not([' + attr + '])').each((k, d) => {
                const $this = $(d);
                $this.val($.trim($this.val()));
            });
            $this.attr(attr, 1);
            $this.trigger('submit');
        }
        else{
            $this.removeAttr(attr);
        }
    });
});