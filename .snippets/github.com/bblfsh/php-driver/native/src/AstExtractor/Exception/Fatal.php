<?php

namespace AstExtractor\Exception;

class Fatal extends BaseFailure
{
    public function __construct(string $msg, ...$values)
    {
        parent::__construct(BaseFailure::FATAL, $msg, ...$values);
    }
}
