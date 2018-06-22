<?php declare(strict_types=1);

namespace AstExtractor\Test\Command;

use PHPUnit\Framework\TestCase;
use AstExtractor\Exception\Error;
use AstExtractor\Exception\Fatal;
use AstExtractor\Command\Request;
use AstExtractor\Command\Response;

class ResponseTest extends TestCase
{
    private $fixtureResponse;
    private const AST = ['level1' => 'ast'];
    private const FIXTURE_NAME = 'fixture_name';

    public function __construct()
    {
        parent::__construct();
        $this->fixtureResponse = new Response(self::AST, self::FIXTURE_NAME);
    }

    public function testNewResponse(): void
    {
        $this->assertEquals(self::FIXTURE_NAME, $this->fixtureResponse->name);
        $this->assertEquals(self::AST, $this->fixtureResponse->ast);
        $this->assertEquals(Response::STATUS_PENDING, $this->fixtureResponse->status);
        $this->assertCount(0, $this->fixtureResponse->errors);
    }

    public function testResponseFromRequestWithError(): void
    {
        $req = new Request('content', self::FIXTURE_NAME);
        $response = Response::fromRequest($req, self::AST, new Error('error'));

        $this->assertEquals(self::FIXTURE_NAME, $response->name);
        $this->assertEquals(self::AST, $response->ast);
        $this->assertEquals(Response::STATUS_ERROR, $response->status);
        $this->assertCount(1, $response->errors);
    }

    public function testResponseFromFailure(): void
    {
        $req = new Request('content', self::FIXTURE_NAME);

        $responseErr = Response::fromError(new Error('error'));
        $this->assertEquals(Response::STATUS_ERROR, $responseErr->status);
        $this->assertCount(1, $responseErr->errors);

        $responseFatal = Response::fromError(new Fatal('fatal'));
        $this->assertEquals(Response::STATUS_FATAL, $responseFatal->status);
        $this->assertCount(1, $responseFatal->errors);
    }

    public function testToArray(): void
    {
        $arr = $this->fixtureResponse->toArray();
        $expected = [
            'ast' => self::AST,
            'metadata' => ['name' => self::FIXTURE_NAME,],
            'status' => Response::STATUS_PENDING,
            'errors' => [],
        ];

        $this->assertEquals($expected, $arr);
    }
}
