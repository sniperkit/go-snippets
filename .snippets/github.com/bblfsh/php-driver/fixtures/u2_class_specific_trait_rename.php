<?php
trait testtrait1 {
    public function testfnc1() {}
}
class testcls1 {
    use testtrait1 { testfnc1 as protected testfnc2; }
}
