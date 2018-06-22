var testcls1 = new Object(); testcls1.foo = 'A'
var testcls2 = Object.create({'foo': 'A'});
var testcls3 = {'foo': 'A'};
var testcls4 = function() {
    this.foo = 'A';
};
class testcls5 {
    constructor() {
        this.foo = 'A';
    }
}

