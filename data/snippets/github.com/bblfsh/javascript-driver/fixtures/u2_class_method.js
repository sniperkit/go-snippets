var testcls1 = new Object(); 
testcls1.testfnc1 = function() {};

var testcls2 = Object.create(null);
testcls2.testfnc2 = function() {};

function testcls1() {
    this.testfnc1 = function() {};
}
function testcls2() {}
testcls2.prototype.testfnc2 = function() {};

var testcls1 = new Object(); testcls1.foo = function() {};
var testcls2 = Object.create({'testfnc2': function() {}});
var testcls3 = {'testfnc3': function() {}};
var testcls4 = function() {
    this.testfnc4 = function() {};
};
class testcls5 {
    testfnc5() {}
}
