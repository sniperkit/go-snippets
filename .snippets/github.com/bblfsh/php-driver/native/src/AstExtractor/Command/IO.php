<?php

namespace AstExtractor\Command;

use AstExtractor\Exception\BaseFailure;
use AstExtractor\Exception\Fatal;
use AstExtractor\Formatter\BaseFormatter;

class IO
{
    private $reader;
    private $writer;
    private $errorWriter;
    private $formatter;

    public function __construct(BaseFormatter $formatter, $reader = null, $writer = null)
    {
        if ($reader === null) {
            $this->reader = fopen('php://stdin', 'rb');
        } elseif (self::isOpened($reader)) {
            $this->reader = $reader;
        } else {
            throw new Fatal('No proper reader provided');
        }

        if ($writer === null) {
            $this->writer = fopen('php://stdout', 'ab');
        } elseif (self::isOpened($reader)) {
            $this->writer = $reader;
        } else {
            throw new Fatal('No proper writer provided');
        }

        $formatter->setReader($this->reader);
        $this->formatter = $formatter;
    }

    public function nextRequest()
    {
        while ($this->isAvailable()) {
            try {
                $rawReq = $this->formatter->readNext();
                return Request::fromArray($rawReq);
            } catch (BaseFailure $e) {
                //TODO: separate in error types
                if ($e->getCode() === BaseFailure::EOF) {
                    throw $e;
                }
                $this->writeErr(null, $e);
            }
        }

        throw new BaseFailure(BaseFailure::EOF, 'End of reader');
    }

    public function write(Response $response)
    {
        $output = $this->formatter->encode($response->toArray()) . PHP_EOL /*. PHP_EOL*/;
        fwrite($this->writer, $output);
    }

    public function writeErr(Request $request = null, \Exception $e) {
        if ($request === null) {
            $response = Response::fromError($e);
        } else {
            $response = $request->answer([]);
            $response->errors = [$e];
            $response->status = Response::getStatus($e->getCode());
        }

        self::write($response);
    }

    public function isAvailable() {
        return self::isOpened($this->reader);
    }

    public static function isOpened($resource)
    {
        return is_resource($resource) && !feof($resource);
    }
}
