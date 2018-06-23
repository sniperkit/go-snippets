<?php declare(strict_types=1);

namespace AstExtractor\Test;

use AstExtractor\Formatter\BaseFormatter;
use PHPUnit\Framework\TestCase;

class BaseFormatterTest extends TestCase
{
    /**
     * @expectedException \AstExtractor\Exception\Fatal
     */
    public function testNonStreamingImplementation(): void
    {
        $formatter = new SerialFormatter(tmpfile());
        $formatter->readNext();
    }

    /**
     * @expectedException \AstExtractor\Exception\Fatal
     */
    public function testInvalidReader(): void
    {
        $formatter = new SerialFormatter();
        $formatter->setReader(null);
    }
}

class SerialFormatter extends BaseFormatter
{
    public function encode(array $message){return '';}
    public function decode(string $input){return [];}
}
