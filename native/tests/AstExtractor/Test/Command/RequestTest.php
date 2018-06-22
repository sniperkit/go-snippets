<?php declare(strict_types=1);

namespace AstExtractor\Test\Command;

use PHPUnit\Framework\TestCase;
use AstExtractor\Command\Request;
use AstExtractor\Command\Response;

class RequestTest extends TestCase
{
    private $fixtureRequest;
    private const FIXTURE_CONTENT_OK = '<?php echo 1;';
    private const FIXTURE_CONTENT_ERR = '<?php echo 1';
    private const FIXTURE_NAME = 'fixture_name';

    public function __construct()
    {
        parent::__construct();
        $this->fixtureRequest = new Request(self::FIXTURE_CONTENT_OK, self::FIXTURE_NAME);
    }

    public function testNewRequest(): void
    {
        $this->assertEquals(self::FIXTURE_NAME, $this->fixtureRequest->name);
        $this->assertEquals(self::FIXTURE_CONTENT_OK, $this->fixtureRequest->content);
    }

    public function testFromArrayRequest(): void
    {
        $req = Request::fromArray([
            'metadata' => ['name' => self::FIXTURE_NAME,],
            'content' => self::FIXTURE_CONTENT_OK,
        ]);

        $this->assertEquals($this->fixtureRequest->name, $req->name);
        $this->assertEquals($this->fixtureRequest->content, $req->content);
    }

    /**
     * @expectedException \AstExtractor\Exception\Fatal
     */
    public function testNewRequestFails(): void
    {
        Request::fromArray([]);
    }

    public function testAnswer(): void
    {
        $ast = ['level1'=>'ast'];
        $response = $this->fixtureRequest->answer($ast);

        $this->assertEquals($this->fixtureRequest->name, $response->name);
        $this->assertEquals($ast, $response->ast);
        $this->assertEquals(Response::STATUS_OK, $response->status);
        $this->assertCount(0, $response->errors);
    }

    /**
     * @expectedException \AstExtractor\Exception\Fatal
     */
    public function testAnswerFails(): void
    {
        $this->fixtureRequest->answer(null);
    }

    public function testToArrayRequest(): void
    {
        $arr = $this->fixtureRequest->toArray();
        $expected = [
            'metadata' => ['name' => self::FIXTURE_NAME,],
            'content' => self::FIXTURE_CONTENT_OK,
        ];

        $this->assertEquals($expected, $arr);
    }
}
