<?php

namespace AstExtractor\Formatter;

use AstExtractor\Exception\BaseFailure;
use AstExtractor\Exception\Fatal;

class Json extends BaseFormatter
{
    /**
     * ENCODING_OPTS forces encoding to "encode multibyte Unicode characters literally"
     *  and "substitute some unencodable values instead of failing"
     */
    private const ENCODING_FALLBACK_OPTS = JSON_UNESCAPED_UNICODE | JSON_PARTIAL_OUTPUT_ON_ERROR;

    /**
     * DECODE_USING_ASSOC_ARRAY is true if the decode value will be represented by an
     *  associative array instead of by an object
     */
    private const DECODE_USING_ASSOC_ARRAY = true;

    /**
     * MAX_DEPTH is the max allowed depth for encoding json structures
     */
    private const MAX_DEPTH = 512;

    /**
     * @inheritdoc
     */
    public function encode(array $message)
    {
        $encode = json_encode($message, 0, self::MAX_DEPTH);
        if (!$encode && json_last_error() === JSON_ERROR_UTF8) {
            self::utf8_encode_recursive($message);
            $encode = json_encode($message, self::ENCODING_FALLBACK_OPTS, self::MAX_DEPTH);
        }

        if (!$encode) {
            throw new Fatal('Error#%s, %s', json_last_error(), json_last_error_msg());
        }

        return $encode;
    }

    /**
     * @inheritdoc
     */
    public function decode(string $input)
    {
        $decoded = json_decode($input, self::DECODE_USING_ASSOC_ARRAY);
        if (!$decoded) {
            throw new Fatal('Error#%s, %s', json_last_error(), json_last_error_msg());
        }

        return $decoded;
    }

    /**
     * @inheritdoc
     */
    public function readNext()
    {
        if (!$this->isReaderOpened()) {
            throw new BaseFailure(BaseFailure::EOF, 'End of reader reached');
        }

        while ($this->isReaderOpened() && $read = fgets($this->reader)) {
            if ($read === false) {
                throw new Fatal('Error reading from passed stream');
            }

            if (trim($read) === "") {
                continue;
            }

            return self::decode($read);
        }

        throw new BaseFailure(BaseFailure::EOF, 'End of reader reached');
    }

    /**
     * utf8_encode_recursive encodes the string contents of the passed $input as valid utf8
     * All inner string contents are scanned, and all array and public object values are converted to utf8
     * @param $input Input data to convert
     * @throws \Exception If something went wrong
     */
    private static function utf8_encode_recursive(&$input)
    {
        if (is_string($input)) {
            $input = utf8_encode($input);
        }

        if (is_object($input)) {
            $ovs = get_object_vars($input);
            foreach ($ovs as $k => $v)    {
                self::utf8_encode_recursive($input->$k);
            }
        }

        if (is_array($input)) {
            $success = array_walk_recursive($input, function(&$v){self::utf8_encode_recursive($v);});
            if (!$success) {
                throw new Fatal('Unexpected error during recursive array utf8 encoding');
            }
        }
    }
}
