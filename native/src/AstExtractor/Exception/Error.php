<?php

namespace AstExtractor\Exception;

class Error extends BaseFailure
{
    public function __construct(string $msg, ...$values)
    {
        parent::__construct(BaseFailure::ERROR, $msg, ...$values);
    }
}
