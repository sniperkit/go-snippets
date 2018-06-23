package tech.sourced.babelfish;

import com.fasterxml.jackson.core.JsonGenerator;
import org.eclipse.cdt.core.dom.ast.*;
import org.eclipse.cdt.core.dom.ast.c.ICASTArrayModifier;
import org.eclipse.cdt.core.dom.ast.c.ICASTDesignator;
import org.eclipse.cdt.core.dom.ast.cpp.*;
import org.eclipse.cdt.internal.core.dom.rewrite.commenthandler.NodeCommentMap;

import java.io.IOException;
import java.util.Arrays;
import java.util.List;

// FIXME: Review all nodes implemented on the ASTVisitor interface
// and do the isinstance of their subclasses inside all call the
// right methods on them, renaming them from visit to visit_something

// TODO: abstract this visitor from the JSON object, use an injected proxy
// class to communicate the node data. This would allow to remove the clumsy
// try-catch-finallys around every visit method. In fact, in order to implement
// the #if preprocessor-branching this should instead write the nodes to a tree
// data structure (with a quickaccess HashMap) that we can modify to alter
// and relocate the branches (under the preprocessor #if nodes).

/// Visitor pattern implementation for the CPP AST. This will write every
/// node in the Json output. Since CDT unfortunately doesnt have something like JDT
/// "structuralPropertiesForType", we must visit ALL the base nodes and check
/// for every derived node because while IASTNode.getChildren will return all the
/// children that derive for IASTNode (including all the ones returned by any method
/// of that node returning an IASTNode-derived object), a lot of nodes also have
/// other methods that return something not derived from IASTNode and this not returned
/// by getChildren. So for all these IASTNode subinterfaces we must call and store the
/// value of those non-IASTNode-returning methods. The most infamous is probably the
/// getOperator of some nodes that return an int than then you have to match in a switch
/// because the possible values are not even declarated in an enum but as final int
/// class members (thank god for [Idea]Vim macros!).

// Many Bothans died to bring us this class
public class JsonASTVisitor extends ASTVisitor {

    private JsonGenerator json;
    private NodeCommentMap commentMap;
    IOException error;
    boolean hasError = false;

    JsonASTVisitor(JsonGenerator json, NodeCommentMap commentMap) {
        super();
        this.json = json;
        this.commentMap = commentMap;

        shouldVisitArrayModifiers = true;
        shouldVisitBaseSpecifiers = true;
        shouldVisitDeclarations = true;
        shouldVisitDeclarators = true;
        shouldVisitDeclSpecifiers = true;
        shouldVisitExpressions = true;
        shouldVisitInitializers = true;
        shouldVisitNames = true;
        shouldVisitNamespaces = true;
        shouldVisitParameterDeclarations = true;
        shouldVisitPointerOperators = true;
        shouldVisitStatements = true;
        shouldVisitTemplateParameters = true;
        shouldVisitTranslationUnit = true;
        shouldVisitTypeIds = true;
    }

    private void enableErrorState(IOException e) {
        error = e;
        hasError = true;
    }

    private void serializeNodeLocation(IASTNode node) throws IOException {
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
        json.writeNumberField("LocLineStart", lineStart);
        json.writeNumberField("LocLineEnd", lineEnd);
        json.writeNumberField("LocOffsetStart", offsetStart);
        json.writeNumberField("LocOffsetLength", offsetLength);
    }

    private void serializeCommonData(IASTNode node) throws IOException {
        json.writeStringField("IASTClass", node.getClass().getSimpleName());
        json.writeStringField("Snippet", EclipseCPPParser.getSnippet(node));

        ASTNodeProperty propInParent = node.getPropertyInParent();
        if (propInParent != null) {
            json.writeStringField("Role", propInParent.getName());
        }

        serializeNodeLocation(node);
    }

    private void serializeCommentList(List<IASTComment> comments, String commentType) throws IOException {
        if (comments != null && comments.size() > 0) {
            json.writeFieldName(commentType + "Comments");
            json.writeStartArray();
            try {
                for (IASTComment comment : comments) {
                    json.writeStartObject();
                    try {
                        json.writeStringField("Comment", comment.toString());
                        json.writeBooleanField("IsBlockComment", comment.isBlockComment());
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

    // We need to call this on each visitor instead of just retuning
    // PROCESS_SKIP and let the base clase do it automatically because
    // we need to start the child JSON array
    public void visitChilds(IASTNode node) throws IOException {
        IASTNode[] children = node.getChildren();

        if (children != null && children.length > 0) {
            json.writeFieldName("childs");
            json.writeStartArray();
            try {
                for (IASTNode child : children) {
                    child.accept(this);
                }
            } finally {
                json.writeEndArray();
            }
        }
    }

    public int visit(IASTNode node) {
        try {
            json.writeStartObject();
            try {
                serializeCommonData(node);
                serializeComments(node);
                visitChilds(node);
            } finally {
                json.writeEndObject();
            }
        } catch (IOException e) {
            enableErrorState(e);
            return PROCESS_ABORT;
        }
        return PROCESS_SKIP;

    }

    public int XXX(IASTASMDeclaration node) {
        try {
            serializeCommonData(node);
            json.writeStringField("ASMContent", node.getAssembly());
            serializeComments(node);
            visitChilds(node);
        } catch (IOException e) {
            enableErrorState(e);
            return PROCESS_ABORT;
        }
        return PROCESS_SKIP;
    }

    private void serializeUnaryExpression(IASTUnaryExpression node) throws IOException {
        int operator = node.getOperator();
        String opStr;

        switch (operator) {
            case IASTUnaryExpression.op_alignOf:
                opStr = "op_alignOf";
                break;
            case IASTUnaryExpression.op_amper:
                opStr = "op_amper";
                break;
            case IASTUnaryExpression.op_bracketedPrimary:
                opStr = "op_bracketedPrimary";
                break;
            case IASTUnaryExpression.op_labelReference:
                opStr = "op_labelReference";
                break;
            case IASTUnaryExpression.op_minus:
                opStr = "op_minus";
                break;
            case IASTUnaryExpression.op_noexcept:
                opStr = "op_noexcept";
                break;
            case IASTUnaryExpression.op_not:
                opStr = "op_not";
                break;
            case IASTUnaryExpression.op_plus:
                opStr = "op_plus";
                break;
            case IASTUnaryExpression.op_postFixDecr:
                opStr = "op_postFixDecr";
                break;
            case IASTUnaryExpression.op_postFixIncr:
                opStr = "op_postFixIncr";
                break;
            case IASTUnaryExpression.op_prefixDecr:
                opStr = "op_prefixDecr";
                break;
            case IASTUnaryExpression.op_prefixIncr:
                opStr = "op_prefixIncr";
                break;
            case IASTUnaryExpression.op_sizeof:
                opStr = "op_sizeof";
                break;
            case IASTUnaryExpression.op_sizeofParameterPack:
                opStr = "op_sizeofParameterPack";
                break;
            case IASTUnaryExpression.op_star:
                opStr = "op_star";
                break;
            case IASTUnaryExpression.op_throw:
                opStr = "op_throw";
                break;
            case IASTUnaryExpression.op_tilde:
                opStr = "op_tilde";
                break;
            case IASTUnaryExpression.op_typeid:
                opStr = "op_typeid";
                break;
            case IASTUnaryExpression.op_typeof:
                opStr = "op_typeof";
                break;
            default:
                opStr = "op_unkown";
                break;
        }
        json.writeStringField("operator", opStr);
    }

    private void serializeTypeIdExpression(IASTTypeIdExpression node) throws IOException {
        int operator = node.getOperator();
        String opStr = "op_unkown";

        switch (operator) {
            case IASTTypeIdExpression.op_alignof:
                opStr = "op_alignof";
                break;
            case IASTTypeIdExpression.op_has_nothrow_assign:
                opStr = "op_has_nothrow_assign";
                break;
            case IASTTypeIdExpression.op_has_nothrow_constructor:
                opStr = "op_has_nothrow_constructor";
                break;
            case IASTTypeIdExpression.op_has_nothrow_copy:
                opStr = "op_has_nothrow_copy";
                break;
            case IASTTypeIdExpression.op_has_trivial_assign:
                opStr = "op_has_trivial_assign";
                break;
            case IASTTypeIdExpression.op_has_trivial_constructor:
                opStr = "op_has_trivial_constructor";
                break;
            case IASTTypeIdExpression.op_has_trivial_copy:
                opStr = "op_has_trivial_copy";
                break;
            case IASTTypeIdExpression.op_has_trivial_destructor:
                opStr = "op_has_trivial_destructor";
                break;
            case IASTTypeIdExpression.op_has_virtual_destructor:
                opStr = "op_has_virtual_destructor";
                break;
            case IASTTypeIdExpression.op_is_abstract:
                opStr = "op_is_abstract";
                break;
            case IASTTypeIdExpression.op_is_class:
                opStr = "op_is_class";
                break;
            case IASTTypeIdExpression.op_is_empty:
                opStr = "op_is_empty";
                break;
            case IASTTypeIdExpression.op_is_enum:
                opStr = "op_is_enum";
                break;
            case IASTTypeIdExpression.op_is_final:
                opStr = "op_is_final";
                break;
            case IASTTypeIdExpression.op_is_literal_type:
                opStr = "op_is_literal_type";
                break;
            case IASTTypeIdExpression.op_is_pod:
                opStr = "op_is_pod";
                break;
            case IASTTypeIdExpression.op_is_polymorphic:
                opStr = "op_is_polymorphic";
                break;
            case IASTTypeIdExpression.op_is_standard_layout:
                opStr = "op_is_standard_layout";
                break;
            case IASTTypeIdExpression.op_is_trivial:
                opStr = "op_is_trivial";
                break;
            case IASTTypeIdExpression.op_is_trivially_copyable:
                json.writeString("op_is_trivially_copyable");
                break;
            case IASTTypeIdExpression.op_is_union:
                opStr = "op_is_union";
                break;
            case IASTTypeIdExpression.op_sizeof:
                opStr = "op_sizeof";
                break;
            case IASTTypeIdExpression.op_sizeofParameterPack:
                opStr = "op_sizeofParameterPack";
                break;
            case IASTTypeIdExpression.op_typeid:
                opStr = "op_typeid";
                break;
            case IASTTypeIdExpression.op_typeof:
                opStr = "op_typeof";
                break;
            default:
        }
        json.writeStringField("operator", opStr);
    }

    private void serializeBinaryExpression(IASTBinaryExpression node) throws IOException {
        int op = node.getOperator();
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
        json.writeStringField("Operator", opStr);
    }

    @Override
    public int visit(IASTName node) {
        try {
            json.writeStartObject();
            try {
                serializeCommonData(node);
                json.writeStringField("Name", node.toString());
                json.writeStringField("FullName", Arrays.toString(node.toCharArray()));
                json.writeBooleanField("IsQualified", node.isQualified());
                json.writeStringField("Binding", node.getBinding().toString());
                json.writeStringField("ResolvedBinding", node.resolveBinding().toString());
                json.writeStringField("PreBinding", node.getPreBinding().toString());
                json.writeStringField("ResolvedPreBinding", node.resolvePreBinding().toString());

                int roleOfName = node.getRoleOfName(true);
                json.writeFieldName("RoleOfName");
                switch (roleOfName) {
                    case IASTNameOwner.r_declaration:
                        json.writeString("declaration");
                        break;
                    case IASTNameOwner.r_definition:
                        json.writeString("definition");
                        break;
                    case IASTNameOwner.r_reference:
                        json.writeString("reference");
                        break;
                    default:
                        json.writeString("unclear");
                }

                if (node instanceof IASTImplicitName) {
                    IASTImplicitName impl = (IASTImplicitName) node;
                    json.writeStringField("IsAlternate", String.valueOf(impl.isAlternate()));
                    json.writeStringField("IsOverloadedOperator", String.valueOf(impl.isOperator()));
                }

                serializeComments(node);
                visitChilds(node);
            } finally {
                json.writeEndObject();
            }
        } catch (IOException e) {
            enableErrorState(e);
            return PROCESS_ABORT;
        }
        return PROCESS_SKIP;
    }

    @Override
    public int visit(IASTExpression node) {
        try {
            json.writeStartObject();
            try {
                serializeCommonData(node);
                serializeComments(node);
                json.writeFieldName("ExpressionType");
                json.writeString(node.getExpressionType().toString());
                json.writeFieldName("ExpressionValueCategory");
                json.writeString(node.getValueCategory().toString());
                json.writeBooleanField("IsLValue", node.isLValue());

                if (node instanceof IASTBinaryExpression) {
                    serializeBinaryExpression((IASTBinaryExpression) node);
                }

                if (node instanceof IASTLiteralExpression) {
                    IASTLiteralExpression lit = (IASTLiteralExpression) node;
                    json.writeStringField("LiteralKind", lit.getExpressionType().toString());
                    json.writeStringField("LiteralValue", lit.toString());
                }

                if (node instanceof IASTTypeIdExpression) {
                    serializeTypeIdExpression((IASTTypeIdExpression) node);
                }

                if (node instanceof IASTBinaryTypeIdExpression) {
                    IASTBinaryTypeIdExpression.Operator operator = ((IASTBinaryTypeIdExpression) node).getOperator();
                    json.writeStringField("BinaryTypeIdOperator", operator.toString());
                }

                if (node instanceof IASTFieldReference) {
                    json.writeBooleanField("IsPointerDereference",
                            ((IASTFieldReference) node).isPointerDereference());
                }

                if (node instanceof ICPPASTDeleteExpression) {
                    json.writeBooleanField("IsGlobal", ((ICPPASTDeleteExpression) node).isGlobal());
                    json.writeBooleanField("IsVectored", ((ICPPASTDeleteExpression) node).isVectored());
                }

                if (node instanceof ICPPASTNaryTypeIdExpression) {
                    ICPPASTNaryTypeIdExpression.Operator operator = ((ICPPASTNaryTypeIdExpression) node).getOperator();
                    json.writeStringField("NaryTypeIdOperator", operator.toString());
                }

                if (node instanceof ICPPASTNewExpression) {
                    json.writeBooleanField("IsArrayAllocation", ((ICPPASTNewExpression) node).isArrayAllocation());
                    json.writeBooleanField("IsGlobal", ((ICPPASTNewExpression) node).isGlobal());
                    json.writeBooleanField("IsNewTypeId", ((ICPPASTNewExpression) node).isNewTypeId());
                }

                if (node instanceof IASTUnaryExpression) {
                    serializeUnaryExpression((IASTUnaryExpression) node);
                }

                if (node instanceof ICPPASTUnaryExpression) {
                    ICPPASTUnaryExpression uexp = (ICPPASTUnaryExpression) node;

                    ICPPFunction overload = uexp.getOverload();
                    if (overload != null) {
                        json.writeStringField("OverloadedBy", overload.toString());
                    }
                }

                visitChilds(node);
            } finally {
                json.writeEndObject();
            }
        } catch (IOException e) {
            enableErrorState(e);
            return PROCESS_ABORT;
        }
        return PROCESS_SKIP;
    }

    @Override
    public int visit(IASTArrayModifier node) {
        try {
            json.writeStartObject();
            try {
                serializeCommonData(node);
                if (node instanceof ICASTArrayModifier) {
                    ICASTArrayModifier carr = (ICASTArrayModifier) node;
                    json.writeBooleanField("IsConst", carr.isConst());
                    json.writeBooleanField("IsRestrict", carr.isRestrict());
                    json.writeBooleanField("IsStatic", carr.isStatic());
                    json.writeBooleanField("IsVolatile", carr.isVolatile());
                }
                serializeComments(node);
                visitChilds(node);
            } finally {
                json.writeEndObject();
            }
        } catch (IOException e) {
            enableErrorState(e);
            return PROCESS_ABORT;
        }
        return PROCESS_SKIP;
    }

    @Override
    public int visit(IASTAttribute node) {
        try {
            json.writeStartObject();
            try {
                serializeCommonData(node);
                json.writeStringField("AttrName", String.valueOf(node.getName()));

                if (node instanceof ICPPASTAttribute) {
                    ICPPASTAttribute attr = (ICPPASTAttribute) node;
                    json.writeBooleanField("HasPackExpansion", attr.hasPackExpansion());
                }

                serializeComments(node);
                visitChilds(node);
            } finally {
                json.writeEndObject();
            }
        } catch (IOException e) {
            enableErrorState(e);
            return PROCESS_ABORT;
        }
        return PROCESS_SKIP;
    }

    @Override
    public int visit(IASTDeclaration node) {
        try {
            json.writeStartObject();
            try {
                serializeCommonData(node);
                // FIXME:
                // IASTASMDeclaration
                // ICPPASTExplicitTemplateInstantiation .getModifier
                // ICPPASTFunctionDefinition .isDefaulted .isDeleted
                // ICPPASTLinkageSpecification .getLiteral
                // ICPPASTNamespaceDefinition .isInline
                // ICPPASTTemplateDeclaration .isExported
                // ICPPASTUsingDeclaration .isTypeName
                // ICPPASTVisibilityLabel .getVisibility
                serializeComments(node);
                visitChilds(node);
            } finally {
                json.writeEndObject();
            }
        } catch (IOException e) {
            enableErrorState(e);
            return PROCESS_ABORT;
        }
        return PROCESS_SKIP;
    }

    @Override
    public int visit(IASTDeclarator node) {
        try {
            json.writeStartObject();
            try {
                serializeCommonData(node);
                // FIXME:
                // IASTStandardFunctionDeclarator .takesVarArgs
                // ICPPASTDeclarator .declaresParametersPack
                // ICPPASTFunctionDeclarator .isConst .isFinal .isMutable .isOverride .isPureVirtual
                //      .isVolatile
                serializeComments(node);
                visitChilds(node);
            } finally {
                json.writeEndObject();
            }
        } catch (IOException e) {
            enableErrorState(e);
            return PROCESS_ABORT;
        }
        return PROCESS_SKIP;
    }

    @Override
    public int visit(IASTDeclSpecifier node) {
        try {
            json.writeStartObject();
            try {
                serializeCommonData(node);
                // FIXME:
                // this: .getStorageClass .isConst .isInline .isRestrict .isVolatile .
                // IASTCompositeTypeSpecifier .getKey
                // IASTElaboratedTypeSpecifier .getKind
                // IASTSimpleDeclSpecifier .getType .isComplex .isImaginary .isLong .isLongLong
                //      .isShort .isSigned .isUnsigned .
                // ICPPASTCompositeTypeSpecifier .isFinal
                // ICPPASTDeclSpecifier .isConstexp .isExplicit .isFriend .isThreadLocal .isVirtual
                // ICPPASTEnumerationSpecifier .isOpaque .isScoped
                // ICPPASTNamedTypeSpecifier .isTypeName
                // ICPPASTTypeTransformationSpecifier .getOperator
                //
                serializeComments(node);
                visitChilds(node);
            } finally {
                json.writeEndObject();
            }
        } catch (IOException e) {
            enableErrorState(e);
            return PROCESS_ABORT;
        }
        return PROCESS_SKIP;
    }

    @Override
    public int visit(IASTInitializer node) {
        try {
            json.writeStartObject();
            try {
                serializeCommonData(node);
                // FIXME:
                // IASTInitializerList .getSize
                serializeComments(node);
                visitChilds(node);
            } finally {
                json.writeEndObject();
            }
        } catch (IOException e) {
            enableErrorState(e);
            return PROCESS_ABORT;
        }
        return PROCESS_SKIP;
    }

    @Override
    public int visit(IASTParameterDeclaration node) {
        try {
            json.writeStartObject();
            try {
                serializeCommonData(node);
                serializeComments(node);
                visitChilds(node);
            } finally {
                json.writeEndObject();
            }
        } catch (IOException e) {
            enableErrorState(e);
            return PROCESS_ABORT;
        }
        return PROCESS_SKIP;
    }

    @Override
    public int visit(IASTPointerOperator node) {
        try {
            json.writeStartObject();
            try {
                serializeCommonData(node);
                // FIXME:
                // IASTPointer .isConst .isRestrict .isVolatile
                // ICPPASTReferenceOperator .isRValueReference
                serializeComments(node);
                visitChilds(node);
            } finally {
                json.writeEndObject();
            }
        } catch (IOException e) {
            enableErrorState(e);
            return PROCESS_ABORT;
        }
        return PROCESS_SKIP;
    }

    @Override
    public int visit(IASTProblem node) {
        try {
            json.writeStartObject();
            try {
                serializeCommonData(node);
                serializeComments(node);
                visitChilds(node);
            } finally {
                json.writeEndObject();
            }
        } catch (IOException e) {
            enableErrorState(e);
            return PROCESS_ABORT;
        }
        return PROCESS_SKIP;
    }

    @Override
    public int visit(IASTStatement node) {
        try {
            json.writeStartObject();
            try {
                serializeCommonData(node);
                // FIXME:
                // ICPPASTCatchHandler .isCatchAll
                serializeComments(node);
                visitChilds(node);
            } finally {
                json.writeEndObject();
            }
        } catch (IOException e) {
            enableErrorState(e);
            return PROCESS_ABORT;
        }
        return PROCESS_SKIP;
    }

    @Override
    public int visit(IASTToken node) {
        try {
            json.writeStartObject();
            try {
                serializeCommonData(node);
                // FIXME:
                // this.getTokenType .getTokenCharImage?
                serializeComments(node);
                visitChilds(node);
            } finally {
                json.writeEndObject();
            }
        } catch (IOException e) {
            enableErrorState(e);
            return PROCESS_ABORT;
        }
        return PROCESS_SKIP;
    }

    @Override
    public int visit(IASTTranslationUnit node) {
        try {
            json.writeStartObject();
            try {
                serializeCommonData(node);
                // FIXME:
                // This must be explored to get the preprocessor directives with:
                // THIS
                // getMacroDefinitions (non built in macros defined in this file)
                // getMacroExpansions (where are they expanded)
                // getBuildtinMacroDefinitions
                // getIncludeDirectives
                // getAllPreprocessorStatements
                // isHeaderUnit
                // flattenLocationsToFile (try it)
                // copy (to modify it)
                // setSignificantMacros
                //
                // ICPPASTTranslationUnit getNamespaceScope, getMemberBindings, isInline
                serializeComments(node);
                visitChilds(node);
            } finally {
                json.writeEndObject();
            }
        } catch (IOException e) {
            enableErrorState(e);
            return PROCESS_ABORT;
        }
        return PROCESS_SKIP;
    }

    // XXX continue here
    @Override
    public int visit(IASTTypeId node) {
        try {
            json.writeStartObject();
            try {
                serializeCommonData(node);
                serializeComments(node);
                visitChilds(node);
            } finally {
                json.writeEndObject();
            }
        } catch (IOException e) {
            enableErrorState(e);
            return PROCESS_ABORT;
        }
        return PROCESS_SKIP;
    }

    @Override
    public int visit(ICASTDesignator node) {
        try {
            json.writeStartObject();
            try {
                serializeCommonData(node);
                serializeComments(node);
                visitChilds(node);
            } finally {
                json.writeEndObject();
            }
        } catch (IOException e) {
            enableErrorState(e);
            return PROCESS_ABORT;
        }
        return PROCESS_SKIP;
    }

    @Override
    public int visit(ICPPASTCapture node) {
        try {
            json.writeStartObject();
            try {
                serializeCommonData(node);
                serializeComments(node);
                visitChilds(node);
            } finally {
                json.writeEndObject();
            }
        } catch (IOException e) {
            enableErrorState(e);
            return PROCESS_ABORT;
        }
        return PROCESS_SKIP;
    }

    @Override
    public int visit(ICPPASTCompositeTypeSpecifier.ICPPASTBaseSpecifier node) {
        try {
            json.writeStartObject();
            try {
                serializeCommonData(node);
                serializeComments(node);
                visitChilds(node);
            } finally {
                json.writeEndObject();
            }
        } catch (IOException e) {
            enableErrorState(e);
            return PROCESS_ABORT;
        }
        return PROCESS_SKIP;
    }

    //    @Override
    public int visit(ICPPASTDecltypeSpecifier node) {
        try {
            json.writeStartObject();
            try {
                serializeCommonData(node);
                serializeComments(node);
                visitChilds(node);
            } finally {
                json.writeEndObject();
            }
        } catch (IOException e) {
            enableErrorState(e);
            return PROCESS_ABORT;
        }
        return PROCESS_SKIP;
    }

    @Override
    public int visit(ICPPASTNamespaceDefinition node) {
        try {
            json.writeStartObject();
            try {
                serializeCommonData(node);
                serializeComments(node);
                visitChilds(node);
            } finally {
                json.writeEndObject();
            }
        } catch (IOException e) {
            enableErrorState(e);
            return PROCESS_ABORT;
        }
        return PROCESS_SKIP;
    }

    @Override
    public int visit(ICPPASTTemplateParameter node) {
        try {
            json.writeStartObject();
            try {
                serializeCommonData(node);
                serializeComments(node);
                visitChilds(node);
            } finally {
                json.writeEndObject();
            }
        } catch (IOException e) {
            enableErrorState(e);
            return PROCESS_ABORT;
        }
        return PROCESS_SKIP;
    }

    @Override
    public int visit(IASTAttributeSpecifier node) {
        try {
            json.writeStartObject();
            try {
                serializeCommonData(node);
                serializeComments(node);
                visitChilds(node);
            } finally {
                json.writeEndObject();
            }
        } catch (IOException e) {
            enableErrorState(e);
            return PROCESS_ABORT;
        }
        return PROCESS_SKIP;
    }

    @Override
    public int visit(ICPPASTDesignator node) {
        try {
            json.writeStartObject();
            try {
                serializeCommonData(node);
                serializeComments(node);
                visitChilds(node);
            } finally {
                json.writeEndObject();
            }
        } catch (IOException e) {
            enableErrorState(e);
            return PROCESS_ABORT;
        }
        return PROCESS_SKIP;
    }

    @Override
    public int visit(ICPPASTVirtSpecifier node) {
        try {
            json.writeStartObject();
            try {
                serializeCommonData(node);
                serializeComments(node);
                visitChilds(node);
            } finally {
                json.writeEndObject();
            }
        } catch (IOException e) {
            enableErrorState(e);
            return PROCESS_ABORT;
        }
        return PROCESS_SKIP;
    }

}
