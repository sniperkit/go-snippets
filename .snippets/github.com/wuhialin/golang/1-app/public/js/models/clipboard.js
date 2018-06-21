/*
* @Author: wuhailin
* @Date:   2017-11-15 11:41:00
* @Last Modified by:   wuhailin
* @Last Modified time: 2017-11-15 11:46:36
*/
define(['jquery', 'clipboard'], ($, Clipboard) => {
    'use strict';
    new Clipboard('[data-clipboard]');
    return Clipboard
});