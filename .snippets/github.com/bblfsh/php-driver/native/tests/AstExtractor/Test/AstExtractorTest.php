<?php declare(strict_types=1);

namespace AstExtractor\Test;

use AstExtractor\AstExtractor;
use AstExtractor\Exception\Error;
use PHPUnit\Framework\TestCase;

class AstExtractorTest extends TestCase
{
    public function testAstParser(): void
    {
        $parser = new AstExtractor(AstExtractor::LEXER_CONF_SIMPLIFIED);
        $ast = $parser->getAst('<?php echo 1;');
        $this->assertEquals('PhpParser\Node\Stmt\Echo_', get_class($ast[0]));
    }

    public function testAstParserSintaxError(): void
    {
        $parser = new AstExtractor(AstExtractor::LEXER_CONF_SIMPLIFIED);
        try {
            $parser->getAst('<?php echo 1');
        } catch (Error $e) {
            $expectedSemicolonPos = strpos($e->getMessage(), 'expecting \';\' on unknown');
            $this->assertTrue( $expectedSemicolonPos !== false);
        }

        $parser = new AstExtractor(AstExtractor::LEXER_CONF_VERBOSE);
        try {
            $parser->getAst('<?php echo 1');
        } catch (Error $e) {
            $expectedSemicolonPos = strpos($e->getMessage(), 'expecting \';\' on line');
            $this->assertTrue( $expectedSemicolonPos !== false);
        }
    }
}
