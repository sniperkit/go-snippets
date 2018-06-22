package tech.sourced.babelfish;

import org.eclipse.cdt.core.dom.ast.*;
import org.eclipse.cdt.core.dom.ast.cpp.ICPPASTVisibilityLabel;
import org.eclipse.cdt.core.dom.parser.IScannerExtensionConfiguration;
import org.eclipse.cdt.core.dom.parser.c.GCCParserExtensionConfiguration;
import org.eclipse.cdt.core.dom.parser.c.GCCScannerExtensionConfiguration;
import org.eclipse.cdt.core.dom.parser.c.ICParserExtensionConfiguration;
import org.eclipse.cdt.core.dom.parser.cpp.GPPParserExtensionConfiguration;
import org.eclipse.cdt.core.dom.parser.cpp.GPPScannerExtensionConfiguration;
import org.eclipse.cdt.core.dom.parser.cpp.ICPPParserExtensionConfiguration;
import org.eclipse.cdt.core.parser.*;
import org.eclipse.cdt.internal.core.dom.parser.AbstractGNUSourceCodeParser;
import org.eclipse.cdt.internal.core.dom.parser.c.GNUCSourceParser;
import org.eclipse.cdt.internal.core.dom.parser.cpp.GNUCPPSourceParser;
import org.eclipse.cdt.internal.core.dom.rewrite.commenthandler.ASTCommenter;
import org.eclipse.cdt.internal.core.dom.rewrite.commenthandler.NodeCommentMap;
import org.eclipse.cdt.internal.core.parser.scanner.CPreprocessor;

import java.util.HashMap;
import java.util.List;
import java.util.Map;

class EclipseCPPParser {
    NodeCommentMap commentMap;
    private static final ScannerInfo GNU_SCANNER_INFO = new ScannerInfo(getGnuMap());

    private static Map<String, String> getGnuMap() {
        Map<String, String> map = new HashMap<String, String>();
        map.put("__GNUC__", Integer.toString(99));
        map.put("__GNUC_MINOR__", Integer.toString(99));
        map.put("__SIZEOF_SHORT__", "2");
        map.put("__SIZEOF_INT__", "4");
        map.put("__SIZEOF_LONG__", "8");
        map.put("__SIZEOF_POINTER__", "8");
        return map;
    }

    private static IScanner createScanner(String code, ParserLanguage lang, IScannerInfo scannerInfo) {

        IScannerExtensionConfiguration configuration;
        if (lang == ParserLanguage.C) {
            configuration = GCCScannerExtensionConfiguration.getInstance(scannerInfo);
        } else {
            configuration = GPPScannerExtensionConfiguration.getInstance(scannerInfo);
        }
        IScanner scanner;
        FileContent fileContent = FileContent.create(code, code.toCharArray());
        scanner = new CPreprocessor(fileContent, scannerInfo, lang, new NullLogService(), configuration,
                IncludeFileContentProvider.getSavedFilesProvider());
        return scanner;
    }

    TranslationUnit parseCPP(String code) {
        AbstractGNUSourceCodeParser parser;
        ICPPParserExtensionConfiguration config;
        IScanner scanner = createScanner(code, ParserLanguage.CPP, GNU_SCANNER_INFO);
        config = new GPPParserExtensionConfiguration();
        parser = new GNUCPPSourceParser(scanner, ParserMode.COMPLETE_PARSE, new NullLogService(),
                config, null);
        parser.setMaximumTrivialExpressionsInAggregateInitializers(Integer.MAX_VALUE);

        IASTTranslationUnit parsed = parser.parse();
        commentMap = ASTCommenter.getCommentedNodeMap(parsed);
        return new TranslationUnit(parsed, commentMap);
    }

    IASTTranslationUnit parseC(String code) {
        AbstractGNUSourceCodeParser parser;
        ICParserExtensionConfiguration config;
        IScanner scanner = createScanner(code, ParserLanguage.C, GNU_SCANNER_INFO);
        config = new GCCParserExtensionConfiguration();
        parser = new GNUCSourceParser(scanner, ParserMode.COMPLETE_PARSE, new NullLogService(),
                config, null);
        parser.setMaximumTrivialExpressionsInAggregateInitializers(Integer.MAX_VALUE);
        // FIXME: add commentMap
        return parser.parse();
    }

    void printAST(String code)
            throws Exception {
        // TODO: use C or CPP depending on the specified language
        TranslationUnit tuWrapper = parseCPP(code);
        IASTTranslationUnit translationUnit = tuWrapper.rootNode;
//        System.out.println("XXX comments:");
//        for (IASTComment comment : translationUnit.getComments()) {
//            System.out.print("toStr: ");
//            System.out.println(comment.toString());
//            System.out.println("location: ");
//            System.out.println(getNodeLocationStr(comment));
//            System.out.print("parent: ");
//            System.out.println(comment.getParent().toString());
//            System.out.print("propertyInParent: ");
//            System.out.println(comment.getPropertyInParent().toString());
//            System.out.print("isBlockComment: ");
//            System.out.println(comment.isBlockComment());
//            System.out.println("");
//        }

        System.err.println("XXX Includes:");
        IASTPreprocessorIncludeStatement[] includes = translationUnit.getIncludeDirectives();
        for (IASTPreprocessorIncludeStatement include : includes) {
            System.err.println("include - " + include.getName());
        }

//        System.out.println("XXX preprocStatements: ");
//		IASTPreprocessorStatement[] preprocessorStatements = translationUnit.getAllPreprocessorStatements();
//        for (IASTPreprocessorStatement ps : preprocessorStatements){
//            System.out.println(ps.toString());
//        }

//        System.out.println("XXX macroDefinitions: ");
//        IASTPreprocessorMacroDefinition[] macroDefinitions = translationUnit.getMacroDefinitions();
//        for (IASTPreprocessorMacroDefinition macroDefinition : macroDefinitions){
//            System.out.println(macroDefinition);
//            if(macroDefinition instanceof IASTPreprocessorFunctionStyleMacroDefinition){
//                IASTPreprocessorFunctionStyleMacroDefinition styleDef = (IASTPreprocessorFunctionStyleMacroDefinition) macroDefinition;
//
//                IASTFunctionStyleMacroParameter[] macroParameters = styleDef.getParameters();
//                for(IASTFunctionStyleMacroParameter mp : macroParameters){
//                    System.out.println("macro parameter: "+mp.getParameter());
//                }
//            }
//            System.out.println(macroDefinition.getName() );
//            System.out.println(macroDefinition.getExpansion());
//        }

//        System.out.println("XXX macroExpansions: ");
//        IASTPreprocessorMacroExpansion[] macroExpansions = translationUnit.getMacroExpansions();
//        for(IASTPreprocessorMacroExpansion macroExpansion : macroExpansions){
//            System.out.println(macroExpansion.getRawSignature());
//            System.out.println(macroExpansion.getMacroReference());
//            System.out.println(macroExpansion.getMacroDefinition());
//        }
        System.out.println("\nAST Tree:");
        printTree(translationUnit, 1);

    }

    static String getNodeLocationStr(IASTNode node) {
        String location;
        try {
            IASTFileLocation locationObj = node.getFileLocation();
            if (locationObj != null) {
                location = " (lines: " + node.getFileLocation().getStartingLineNumber() + "," +
                        node.getFileLocation().getEndingLineNumber() +
                        " offset: " + node.getFileLocation().getNodeOffset() + "," +
                        node.getFileLocation().getNodeLength() + ")";
            } else {
                location = " (no location) ";
            }
        } catch (UnsupportedOperationException e) {
            location = "UnsupportedOperationException";
        }
        return location;
    }

    public static String getSnippet(IASTNode node) {
        String snippet = node.getRawSignature().replaceAll("\n", " \\ ");
        int maxlen = node.getRawSignature().length();
        if (maxlen > 20) maxlen = 20;
        return snippet.substring(0, maxlen) + (maxlen == 20 ? "..." : "");
    }

    private void printTree(IASTNode node, int index) {
        // FIXME: move the comment extraction and this doc to another method
        // FIXME: comments at the start of the line doesn't seem to be extracted by the
        // commentMap, add it ourselves to the IASTTranslationUnit

        // Leading comments are before the line of the node, both in previous lines (without
        // any other node between them because then the comments would be from that node)
        // or in the same line just before the node without any other node from the start
        // of the line or the comment (because in that case it would be a trailing comment
        // of the node that is before in the same line)
        // Examples:
        // /* leading of the full IASTDeclarationStatement */ int x = 25;

        // // this comment leading of the ASTReturnStatement
        // // this one too
        // return something;

        // Hint for the AST -> Code generator: these cases can be distinguished by the IASTComment
        // location line (same line as the associated node = leading in the same line)

        List<IASTComment> leadingComments = commentMap.getLeadingCommentsForNode(node);
        for (IASTComment comment : leadingComments) {
            System.out.println(comment.toString());
        }

        // Freestanding comments are comments at the end of IASTCompoundStatements before the closing "}"
        // Example:
        // int someFunc() {
        //   auto foo = 42;
        //   // this is a freestandingComment associated with the IASTFunctionDeclaration
        // }
        List<IASTComment> freestandingComments = commentMap.getFreestandingCommentsForNode(node);
        for (IASTComment comment : freestandingComments) {
            System.out.println(comment.toString());
        }

        // Trailing comment are just after the node in the same line, examples:
        // int a = 3; // this comment will be a trailingComment of the root node of the full IASTDeclaration (int a = 3)
        // int a /* trailing comment of the ASTName (a) node */ = 3;
        List<IASTComment> trailingComments = commentMap.getTrailingCommentsForNode(node);
        for (IASTComment comment : trailingComments) {
            System.out.println(comment.toString());
        }

        String snippet = getSnippet(node);
        String prefix = String.format("%1$" + index * 2 + "s", "-");
        System.out.println(prefix + "[[ " + node.getClass().getSimpleName() + " ]]");
        System.out.println(prefix + "Location: " + getNodeLocationStr(node) + " -> " + snippet);
        ASTNodeProperty propInParent = node.getPropertyInParent();
        if (propInParent != null) {
            System.out.println(prefix + "PropertyInParent: " + propInParent.getName());
        }
        // FIXME: get this into its own method
        if (node instanceof IASTLiteralExpression) {
            IASTLiteralExpression lit = (IASTLiteralExpression) node;
            System.out.println(prefix + "LiteralValue: " + lit.toString());
            System.out.println(prefix + "LiteralKind: " + lit.getExpressionType().toString());
        }

        if (node instanceof IASTName) {
            IASTName name = (IASTName) node;
            System.out.println(prefix + "Name: " + name.toString());
        }
        IASTNode[] children = node.getChildren();
        System.out.println("");

        for (IASTNode iastNode : children)
            printTree(iastNode, index + 1);
    }

    static boolean isVisible(IASTNode current) {
        IASTNode declator = current.getParent().getParent();
        IASTNode[] children = declator.getChildren();

        for (IASTNode iastNode : children) {
            if ((iastNode instanceof ICPPASTVisibilityLabel)) {
                return 1 == ((ICPPASTVisibilityLabel) iastNode).getVisibility();
            }
        }
        return false;
    }
}
