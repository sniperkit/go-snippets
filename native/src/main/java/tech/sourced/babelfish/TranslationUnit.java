package tech.sourced.babelfish;

import org.eclipse.cdt.core.dom.ast.IASTTranslationUnit;
import org.eclipse.cdt.internal.core.dom.rewrite.commenthandler.NodeCommentMap;

public class TranslationUnit {
    IASTTranslationUnit rootNode;
    NodeCommentMap commentMap;

    public TranslationUnit(IASTTranslationUnit rootNode, NodeCommentMap commentMap) {
        this.rootNode = rootNode;
        this.commentMap = commentMap;
    }
}
