package tech.sourced.babelfish;

import com.fasterxml.jackson.core.JsonGenerator;
import com.fasterxml.jackson.databind.SerializerProvider;
import com.fasterxml.jackson.databind.ser.std.StdSerializer;
import org.eclipse.cdt.core.dom.ast.*;
import org.eclipse.cdt.core.dom.ast.cpp.CPPASTVisitor;
import org.eclipse.cdt.core.dom.ast.cpp.ICPPASTCompositeTypeSpecifier;
import org.eclipse.cdt.core.dom.ast.cpp.ICPPASTVisitor;
import org.eclipse.cdt.internal.core.dom.rewrite.commenthandler.NodeCommentMap;
import org.eclipse.cdt.internal.core.pdom.indexer.IndexerASTVisitor;
import org.eclipse.jdt.core.dom.CastExpression;

import java.io.IOException;
import java.util.List;

/* TODO: other interesting methods from IASTNode implementers that are not given as children:
 * http://help.eclipse.org/neon/index.jsp
 * CompoundStatement.getScope.getScopeName
 * ElaboratedTypeSpecifier .getKind
 * ForStatement .getScope.getName
 * FunctionDefinition .getScope.getName
 * FunctionStyleMacroParameter .getParameter
 * InitializerList getSize
 * Pointer isConst, isRestrict, isVolatile, setConst, setRestrict, setVolatile
 * Preprocessor MACRO_NAME.toString
 * PreprocessorStatement MACRO_NAME.toString
 * PreprocessorIfStatement getCondition taken
 * PreprocessorElifStatement getCondition taken
 * PreprocessorElseStatement taken
 * PreprocessorIfDefStatement getCondition getMacroReference taken
 * PreprocessorIfDefStatement getCondition getMacroReference taken
 * PreprocessorMacroDefinition getExpansion getExpansionLocation getName isActive
 * PreprocessorMacroExpansion getMacroDefinition getMacroReference getNestedMacroReference
 * PreprocessorFunctionStyleMacroDefinition PARAMETER? getParameters
 * PreprocessorPragmaStatement getMessage isPragmaOperator
 * PreprocessorUndefStatement getMacroName isActive
 * SimpleDeclSpecifier getType (hacer switch) isComplex isImaginary isLong isLongLong isShort
 *      isSigned isUnsigned
 * StandardFunctionDeclarator getFunctionScope.toString takesVarArgs
 * TokenList getTokenType
 *
 * SPECIAL: TranslationUnit
 *      getAllPreprocessorStatements
 *      getBuiltInMacroDefinitions
 *      getComments
 *      getIncludeDirectives
 *      getMacroDefinitions
 *      getMacroExpansions
 *      getNodeSelector => This finds nodes by file offsets, can be useful when branching
 *          the AST by preprocessor branches or unrolling macros.
 *      getOriginalTranslationUnit
 *      getPreprocessorProblems
 *      getScope
 *      hasNodesOmitted
 */

/**
 * Custom Jackson serializer for org.eclipse.cdt.core.dom.IASTTranslationUnit
 */
public class TranslationUnitSerializer extends StdSerializer<TranslationUnit>
{
    // TODO: add the includes and other macro information to the root node
    // in the JSON

    private NodeCommentMap commentMap;
    JsonGenerator json;

    TranslationUnitSerializer() {
        this(null);
    }

    private TranslationUnitSerializer(Class<TranslationUnit> t) {
        super(t);
    }

    @Override
    public void serialize(TranslationUnit unit, JsonGenerator jsonGenerator,
                          SerializerProvider provider) throws IOException {

        JsonASTVisitor visitor = new JsonASTVisitor(jsonGenerator, unit.commentMap);

        this.json = jsonGenerator;
        this.commentMap = unit.commentMap;

        unit.rootNode.accept(visitor);
//      serializeNode(unit.rootNode);

        if (visitor.hasError && visitor.error != null)  {
            throw visitor.error;
        }
        // TODO: close the jsonGenerator? Check that this doesnt close the associated
        // outputstream
    }

    private void serializeNode(IASTNode node) throws IOException {
        // FIXME: divide this into several methods by node type
        json.writeFieldName("IASTClass");
        json.writeString(node.getClass().getSimpleName());

        json.writeFieldName("Snippet");
        // FIXME: move getSnippet here
        json.writeString(EclipseCPPParser.getSnippet(node));

        serializeNodeLocation(node);

        ASTNodeProperty propInParent = node.getPropertyInParent();
        if (propInParent != null) {
            json.writeFieldName("Role");
            json.writeString(propInParent.getName());
        }

        if (node instanceof IASTLiteralExpression) {
            IASTLiteralExpression lit = (IASTLiteralExpression) node;

            json.writeFieldName("LiteralKind");
            json.writeString(lit.getExpressionType().toString());

            json.writeFieldName("LiteralValue");
            json.writeString(lit.toString());
        }

        if (node instanceof IASTName) {
            IASTName name = (IASTName) node;
            json.writeFieldName("SymbolName");
            json.writeString(name.toString());
        }

        if (node instanceof IASTExpression) {
            IASTExpression expr = (IASTExpression) node;
            json.writeFieldName("ExpressionType");
            json.writeString(expr.getExpressionType().toString());
            json.writeFieldName("ExpressionValueCategory");
            json.writeString(expr.getValueCategory().toString());
            // FIXME: add isLValue
        }

        if (node instanceof IASTASMDeclaration) {
            IASTASMDeclaration asm = (IASTASMDeclaration) node;
            json.writeFieldName("ASMContent");
            json.writeString(asm.getAssembly());
        }

        if (node instanceof IASTBinaryExpression) {
            serializeBinaryOperator(node);
        }

        serializeComments(node);

        IASTNode[] children = node.getChildren();
        if (children != null && children.length > 0) {
            json.writeFieldName("childs");
            json.writeStartArray();
            try {
                for (IASTNode child : children) {
                    json.writeStartObject();
                    try {
                        serializeNode(child);
                    } finally {
                        json.writeEndObject();
                    }
                }
            } finally {
                json.writeEndArray();
            }
        }
    }

    // FIXME: remove
    void serializeNodeLocation(IASTNode node) throws IOException {
        IASTFileLocation loc = node.getFileLocation();
        int lineStart = -1;
        int lineEnd = -1;
        int offsetStart = -1;
        int offsetLength = -1;

        if (loc != null) {
            lineStart = loc.getStartingLineNumber();
            lineEnd = loc.getEndingLineNumber();
            offsetStart = loc.getNodeOffset();
            offsetLength = loc.getNodeLength();
        }
        json.writeFieldName("LocLineStart");
        json.writeNumber(lineStart);
        json.writeFieldName("LocLineEnd");
        json.writeNumber(lineEnd);
        json.writeFieldName("LocOffsetStart");
        json.writeNumber(offsetStart);
        json.writeFieldName("LocOffsetLength");
        json.writeNumber(offsetLength);
    }

    private void serializeBinaryOperator(IASTNode node) throws IOException {
        IASTBinaryExpression binExpr = (IASTBinaryExpression) node;
        int op = binExpr.getOperator();
        String opStr = "UNKOWN OPERATOR WITH CODE " + Integer.toString(op);
        switch (op) {
            case IASTBinaryExpression.op_assign:
                opStr = "assignement =";
                break;
            case IASTBinaryExpression.op_binaryAnd:
                opStr = "binary And &";
                break;
            case IASTBinaryExpression.op_binaryAndAssign:
                opStr = "binary And assign &=";
                break;
            case IASTBinaryExpression.op_binaryOr:
                opStr = "binary or |";
                break;
            case IASTBinaryExpression.op_binaryOrAssign:
                opStr = "Or assign |=";
                break;
            case IASTBinaryExpression.op_binaryXor:
                opStr = "Xor ^";
                break;
            case IASTBinaryExpression.op_binaryXorAssign:
                opStr = "Xor assign ^=";
                break;
            case IASTBinaryExpression.op_divide:
                opStr = "/";
                break;
            case IASTBinaryExpression.op_divideAssign:
                opStr = "assignemnt /=";
                break;
            case IASTBinaryExpression.op_ellipses:
                opStr = "gcc compilers, only.";
                break;
            case IASTBinaryExpression.op_equals:
                opStr = "==";
                break;
            case IASTBinaryExpression.op_greaterEqual:
                opStr = "than or equals >=";
                break;
            case IASTBinaryExpression.op_greaterThan:
                opStr = "than >";
                break;
            case IASTBinaryExpression.op_lessEqual:
                opStr = "than or equals <=";
                break;
            case IASTBinaryExpression.op_lessThan:
                opStr = "than <";
                break;
            case IASTBinaryExpression.op_logicalAnd:
                opStr = "and &&";
                break;
            case IASTBinaryExpression.op_logicalOr:
                opStr = "or ||";
                break;
            case IASTBinaryExpression.op_max:
                opStr = "g++, only.";
                break;
            case IASTBinaryExpression.op_min:
                opStr = "g++, only.";
                break;
            case IASTBinaryExpression.op_minus:
                opStr = "-";
                break;
            case IASTBinaryExpression.op_minusAssign:
                opStr = "assignment -=";
                break;
            case IASTBinaryExpression.op_modulo:
                opStr = "%";
                break;
            case IASTBinaryExpression.op_moduloAssign:
                opStr = "assignment %=";
                break;
            case IASTBinaryExpression.op_multiply:
                opStr = "*";
                break;
            case IASTBinaryExpression.op_multiplyAssign:
                opStr = "assignment *=";
                break;
            case IASTBinaryExpression.op_notequals:
                opStr = "equals !";
                break;
            case IASTBinaryExpression.op_plus:
                opStr = "+";
                break;
            case IASTBinaryExpression.op_plusAssign:
                opStr = "assignment +=";
                break;
            case IASTBinaryExpression.op_pmarrow:
                opStr = "c++, only.";
                break;
            case IASTBinaryExpression.op_pmdot:
                opStr = "c==, only.";
                break;
            case IASTBinaryExpression.op_shiftLeft:
                opStr = "left <<";
                break;
            case IASTBinaryExpression.op_shiftLeftAssign:
                opStr = "left assignment <<=";
                break;
            case IASTBinaryExpression.op_shiftRight:
                opStr = "right >>";
                break;
            case IASTBinaryExpression.op_shiftRightAssign:
                opStr = "right assign >>=";
                break;
        }
        json.writeFieldName("Operator");
        json.writeString(opStr);
    }

    private void serializeCommentList(List<IASTComment> comments, String commentType) throws IOException {
        if (comments != null && comments.size() > 0) {
            json.writeFieldName(commentType + "Comments");
            json.writeStartArray();
            try {
                for (IASTComment comment : comments) {
                    json.writeStartObject();
                    try {
                        json.writeFieldName(commentType + "RelatedComments");
                        json.writeString(comment.toString());

                        json.writeFieldName("IsBlockComment");
                        json.writeBoolean(comment.isBlockComment());

                        serializeNodeLocation(comment);
                    } finally {
                        json.writeEndObject();
                    }
                }
            } finally {
                json.writeEndArray();
            }
        }
    }

    private void serializeComments(IASTNode node) throws IOException {
        serializeCommentList(commentMap.getLeadingCommentsForNode(node), "Leading");
        serializeCommentList(commentMap.getFreestandingCommentsForNode(node), "Freestading");
        serializeCommentList(commentMap.getTrailingCommentsForNode(node), "Trailing");
    }
}
