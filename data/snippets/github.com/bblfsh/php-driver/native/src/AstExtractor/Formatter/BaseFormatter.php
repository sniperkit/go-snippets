<?php

namespace AstExtractor\Formatter;

use AstExtractor\Command\IO;
use AstExtractor\Exception\Fatal;

/**
 * Class BaseFormatter modelates an encoder/decoder with streaming capabilities
 * @package AstExtractor\Formatter
 */
abstract class BaseFormatter
{
    public const MSGPACK = 'msgpack';
    public const JSON = 'json';

    protected $reader;

    /**
     * encode returns a string representation of the passed $input.
     * @param array $input
     * @throws \Exception if the encoding could not be made
     * @return string
     */
    public abstract function encode(array $message);

    /**
     * decode returns the decode value given a string representation.
     * @param string $input
     * @throws \Exception if the decoding could not be made
     * @return array
     */
    public abstract function decode(string $input);

    /**
     * readNext reads the next value from the reader, and returns it.
     * Each read action will block the process till a new value is read
     * @throws \Exception if the reading of the next value could not be made
     * @return array
     */
    public function readNext()
    {
        // TODO: Implement next() method.
        throw new Fatal('No streaming reader available');
    }

    /**
     * BaseFormatter constructor.
     */
    public function __construct(){}

    /**
     * readFrom configure the BaseFormater to use the passed reader
     * The passed $reader will be set as "blocker", so each read will block the
     *   process till a new value is read.
     * @param $reader
     * @throws \Exception
     */
    public function setReader($reader)
    {
        if (!is_resource($reader) || !stream_set_blocking($reader, true)) {
            throw new Fatal('The formatter needs a valid reader, but wrong passed.');
        }

        $this->reader = $reader;
    }

    /**
     * isReaderOpened returns true if the underlying reader is readable
     * @return bool
     */
    protected function isReaderOpened() {
        return IO::isOpened($this->reader);
    }
}
