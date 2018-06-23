/**
 * Created by wuhailin on 2017/11/27.
 */
define([
    'text!data/warehouse.json'
], (json) => {
    try {
        let result = JSON.parse(json);
        return (result && result['RECORDS']) || [];
    }
    catch (e) {
        //
    }
});