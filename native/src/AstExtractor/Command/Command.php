<?php

namespace AstExtractor\Command;

use AstExtractor\Exception\BaseFailure;
use AstExtractor\AstExtractor;
use AstExtractor\Formatter\Json;

class Command
{
    private $extractor;
    private $io;

    public function __construct()
    {
        $formatter = new Json();
        $this->io = new IO($formatter);
        $this->extractor = new AstExtractor(AstExtractor::LEXER_CONF_VERBOSE);
    }

    public static function run(array $argv)
    {
        return (new Command())->init();
    }

    private function init()
    {
        $request = null;
        while ($this->io->isAvailable()) {
            try {
                $request = $this->io->nextRequest();
                $ast = ["nodeType" => "Module", "children" => $this->extractor->getAst($request->content)];
                $response = $request->answer($ast);
                $this->io->write($response);
            } catch (BaseFailure $e) {
                //TODO: separate in error types
                if ($e->getCode() === BaseFailure::EOF) {
                    return true;
                }
                $this->io->writeErr($request, $e);
                continue;
            }
        }

        return true;
    }
}
