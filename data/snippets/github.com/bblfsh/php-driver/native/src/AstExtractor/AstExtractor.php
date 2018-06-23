<?php

/*
 * Part of bblfsh project
 */
namespace AstExtractor;

use AstExtractor\Exception\Error;
use PhpParser\ParserFactory;
use PhpParser\Lexer;
use PhpParser\Node;

/**
 * Extractor returns the AST
 *
 * @category Extractor
 * @package  AstExtractor
 * @author   bblfish <hello@bblfish.com>
 * @license  GPL https://github.com/bblfsh
 * @link     https://github.com/bblfsh
 */
class AstExtractor
{
    public const LEXER_CONF_SIMPLIFIED = 1;
    public const LEXER_CONF_VERBOSE = 2;
    private const LEXER_CONFS = [
        self::LEXER_CONF_SIMPLIFIED => ['usedAttributes' => []],
        self::LEXER_CONF_VERBOSE => ['usedAttributes' => [
            'comments',
            'startLine', 'endLine',
            'startTokenPos', 'endTokenPos',
            'startFilePos', 'endFilePos'
        ]],
    ];
    public const PARSER_VERSION = ParserFactory::PREFER_PHP7;

    public $lexer;
    public $parser;

    /**
     * Returns a new Extractor
     */
    public function __construct(int $conf)
    {
        $this->lexer = new Lexer(self::LEXER_CONFS[$conf]);
        $this->parser = (new ParserFactory)->create(self::PARSER_VERSION, $this->lexer);
    }

    /**
     * Returns the AST. Node kinds can be read here:
     * https://github.com/nikic/php-ast#ast-node-kinds
     *
     * @param string $code Source code to analyze
     *
     * @throws \Exception in case it happens a parsing exception
     * @return Node[]|null Array of statements (or null if the 'throwOnError'
     *      option is disabled and the parser was unable to recover from an error).
     */
    public function getAst($code)
    {
        try {
            return $this->parser->parse($code);
        } catch (\Exception $e) {
            throw new Error($e->getMessage());
        }
    }
}

