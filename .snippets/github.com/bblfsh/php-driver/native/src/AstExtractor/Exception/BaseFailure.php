<?php

namespace AstExtractor\Exception;

class BaseFailure extends \Exception
{
    public const ERROR = 1;
    public const FATAL = 2;

    public const EOF = -1;

    public function __construct($type, string $msg, ...$values)
    {
        $msg = count($values) == 0 ? $msg : sprintf($msg, ...$values);
        parent::__construct($msg, $type, null);
    }
}
