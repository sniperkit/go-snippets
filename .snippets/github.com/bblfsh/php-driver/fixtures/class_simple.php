<?php

class A extends B implements C, D {
    const A = 'B', C = 'D';

    public $a = 'b', $c = 'd';
    protected $e;
    private $f;

    public function a() {}
    public static function b($a) {}
    public final function c() : B {}
    protected function d() {}
    private function e() {}
}
