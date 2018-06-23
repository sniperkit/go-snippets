<?php declare(strict_types=1);

namespace AstExtractor\Test\Formatter;

use AstExtractor\Exception\BaseFailure;
use AstExtractor\Exception\Fatal;
use AstExtractor\Formatter\Json;
use PHPUnit\Framework\TestCase;

class JsonTest extends TestCase
{
    private const JSON_OK = '{"a":1,"b":[1,2],"c":{"a":1}}';
    private const JSON_FAIL_SYNTAX = '{"a":1,';
    private const DECODED_OK = ['a'=>1, 'b'=>[1,2], 'c'=>['a'=>1]];

    private $jsonFormatter;

    public function __construct()
    {
        parent::__construct();
        $this->jsonFormatter = self::getJsonStreamFormatter('');
    }

    public function testDecode(): void
    {
        $this->assertEquals(self::DECODED_OK, $this->jsonFormatter->decode(self::JSON_OK));
    }

    /**
     * @expectedException \AstExtractor\Exception\Fatal
     */
    public function testDecodeFail(): void
    {
        $this->jsonFormatter->decode(self::JSON_FAIL_SYNTAX);
    }

    public function testStreamingDecoding(): void
    {
        $string = self::JSON_OK . PHP_EOL .     //A: should succeed
            PHP_EOL .                           //B: should be ignored
            self::JSON_FAIL_SYNTAX . PHP_EOL .  //C: should throw a Fatal
            '   {"b":99} ' . PHP_EOL .          //D: should succeed
            self::JSON_OK . self::JSON_OK . PHP_EOL . //E: should throw a Fatal
            '   {"b":99} ' . PHP_EOL .          //F: should succeed
            PHP_EOL;                            //G: should be ignored

        $jsonFormatter = self::getJsonStreamFormatter($string);

        //A: should succeed
        $this->assertEquals(self::DECODED_OK, $jsonFormatter->readNext());

        //C: should throw a Fatal
        try {
            $this->assertNull($jsonFormatter->readNext(), '$jsonFormatter->readNext() should throw an exception');
        } catch(Fatal $e) {
            $this->assertEquals(BaseFailure::FATAL, $e->getCode());
        }

        //D: should succeed
        $this->assertEquals(["b"=>99], $jsonFormatter->readNext());

        //E: should throw a Fatal
        try {
            $this->assertNull($jsonFormatter->readNext(), '$jsonFormatter->readNext() should throw an exception');
        } catch(Fatal $e) {
            $this->assertEquals(BaseFailure::FATAL, $e->getCode());
        }

        //F: should succeed
        $this->assertEquals(["b"=>99], $jsonFormatter->readNext());

        //END
        try {
            $this->assertNull($jsonFormatter->readNext(), '$jsonFormatter->readNext() should throw EOF');
        } catch (BaseFailure $e) {
            $this->assertEquals(BaseFailure::EOF, $e->getCode());
        }
    }

    public function testStreamingDecodingOnClosedStream(): void
    {
        $stream = tmpfile();
        $jsonFormatter = new Json($stream);
        fclose($stream);

        //END
        try {
            $this->assertNull($jsonFormatter->readNext(), '$jsonFormatter->readNext() should throw EOF');
        } catch (BaseFailure $e) {
            $this->assertEquals(BaseFailure::EOF, $e->getCode());
        }
    }

    public function testEncode(): void
    {
        $arr = ['a'=>1, 'b'=>[1,2], 'c'=>['a'=>1]];
        $this->assertEquals(self::JSON_OK, $this->jsonFormatter->encode($arr));
    }

    public function testEncodeStrangeChars(): void
    {
        $chars_in = "\x80\x83\x8a\x8e\x9a\x9e"
            ."\x9f\xa2\xa5\xb5\xc0\xc1\xc2"
            ."\xc3\xc4\xc5\xc7\xc8\xc9\xca"
            ."\xcb\xcc\xcd\xce\xcf\xd1\xd2"
            ."\xd3\xd4\xd5\xd6\xd8\xd9\xda"
            ."\xdb\xdc\xdd\xe0\xe1\xe2\xe3"
            ."\xe4\xe5\xe7\xe8\xe9\xea\xeb"
            ."\xec\xed\xee\xef\xf1\xf2\xf3"
            ."\xf4\xf5\xf6\xf8\xf9\xfa\xfb"
            ."\xfc\xfd\xff";
        $chars_out = utf8_encode($chars_in);

        $stdObject = new \stdClass();
        $stdObject->prop1 = 1;
        $stdObject->prop2 = $chars_in;

        $arr_strange_chars = [
            'a' => $chars_in,
            'b' => [$chars_in,2],
            'c' => ['a'=>$chars_in],
            'd' => $stdObject
        ];
        $expected = sprintf(
            '{"a":"%s","b":["%s",2],"c":{"a":"%s"},"d":{"prop1":1,"prop2":"%s"}}',
            $chars_out, $chars_out, $chars_out, $chars_out
        );

        $this->assertEquals($expected, $this->jsonFormatter->encode($arr_strange_chars));
    }

    public function testEncodeNotTooDeepStructure(): void
    {
        $deepth = 512;
        $deepArray = self::deepArray($deepth);
        $this->assertNotEmpty($this->jsonFormatter->encode($deepArray));
    }

    /**
     * @expectedException \AstExtractor\Exception\Fatal
     */
    public function testEncodeTooDeepStructure(): void
    {
        $deepth = 512 + 1;
        $deepArray = self::deepArray($deepth);
        $this->jsonFormatter->encode($deepArray);
    }

    private static function deepArray(int $deepness): array
    {
        $arr = ['a'=>'END'];
        while (--$deepness > 0) {
            $arr = ['a' => $arr];
        }

        return $arr;
    }

    private static function getJsonStreamFormatter(string $string): Json {
        if ($string !== '') {
            $stream = fopen('php://memory','r+');
            fwrite($stream, $string);
            rewind($stream);
        } else {
            $stream = tmpfile();
        }

        $jsonFormatter = new Json();
        $jsonFormatter->setReader($stream);
        return $jsonFormatter;
    }
}
