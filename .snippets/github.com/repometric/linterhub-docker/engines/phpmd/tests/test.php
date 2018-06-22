<?php
    class Something {
    private $q = 15;
    public static function main(array $as, $flag = true) {
        $r = 20 + $this->q;
        for (int $i = 0; $i < 10; $i++) {
            $r += $this->q;
        }
    }
}
?>
