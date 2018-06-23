/*
 * Copyright 2014-2015 The DevMine Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
package ch.devmine.javaparser;

import ch.devmine.javaparser.structures.ArrayExpression;
import ch.devmine.javaparser.structures.ArrayLit;
import ch.devmine.javaparser.structures.ArrayType;
import ch.devmine.javaparser.structures.AssignStatement;
import ch.devmine.javaparser.utils.Log;
import ch.devmine.javaparser.structures.ClassDecl;
import ch.devmine.javaparser.structures.Method;
import ch.devmine.javaparser.structures.Stmt;
import ch.devmine.javaparser.structures.Attribute;
import ch.devmine.javaparser.structures.AttributeRef;
import ch.devmine.javaparser.structures.BasicLit;
import ch.devmine.javaparser.structures.BinaryExpression;
import ch.devmine.javaparser.structures.CallExpression;
import ch.devmine.javaparser.structures.CaseClause;
import ch.devmine.javaparser.structures.CatchClause;
import ch.devmine.javaparser.structures.ClassRef;
import ch.devmine.javaparser.structures.ConstructorDecl;
import ch.devmine.javaparser.structures.ConstructorCallExpr;
import ch.devmine.javaparser.structures.DeclStatement;
import ch.devmine.javaparser.structures.EnumDecl;
import ch.devmine.javaparser.structures.Expr;
import ch.devmine.javaparser.structures.ExprStatement;
import ch.devmine.javaparser.structures.Field;
import ch.devmine.javaparser.structures.FuncRef;
import ch.devmine.javaparser.structures.FuncType;
import ch.devmine.javaparser.structures.Ident;
import ch.devmine.javaparser.structures.IfStatement;
import ch.devmine.javaparser.structures.IncDecExpr;
import ch.devmine.javaparser.structures.IndexExpression;
import ch.devmine.javaparser.structures.InterfaceDecl;
import ch.devmine.javaparser.structures.InterfaceRef;
import ch.devmine.javaparser.structures.LoopStatement;
import ch.devmine.javaparser.structures.OtherStatement;
import ch.devmine.javaparser.structures.ProtoDecl;
import ch.devmine.javaparser.structures.RangeLoopStatement;
import ch.devmine.javaparser.structures.ReturnStatement;
import ch.devmine.javaparser.structures.SwitchStatement;
import ch.devmine.javaparser.structures.TernaryExpression;
import ch.devmine.javaparser.structures.ThrowStatement;
import ch.devmine.javaparser.structures.TryStatement;
import ch.devmine.javaparser.structures.UnaryExpression;
import ch.devmine.javaparser.structures.ValueSpec;
import ch.devmine.javaparser.utils.ParserUtils;
import com.github.javaparser.JavaParser;
import com.github.javaparser.ParseException;
import com.github.javaparser.ast.CompilationUnit;
import com.github.javaparser.ast.ImportDeclaration;
import com.github.javaparser.ast.PackageDeclaration;
import com.github.javaparser.ast.body.BodyDeclaration;
import com.github.javaparser.ast.body.ConstructorDeclaration;
import com.github.javaparser.ast.body.FieldDeclaration;
import com.github.javaparser.ast.body.MethodDeclaration;
import com.github.javaparser.ast.body.ModifierSet;
import com.github.javaparser.ast.body.TypeDeclaration;
import java.io.FileNotFoundException;
import java.util.ArrayList;
import java.util.List;
import com.github.javaparser.ast.body.ClassOrInterfaceDeclaration;
import com.github.javaparser.ast.body.EnumDeclaration;
import com.github.javaparser.ast.body.Parameter;
import com.github.javaparser.ast.body.VariableDeclarator;
import com.github.javaparser.ast.comments.Comment;
import com.github.javaparser.ast.expr.ArrayAccessExpr;
import com.github.javaparser.ast.expr.ArrayCreationExpr;
import com.github.javaparser.ast.expr.ArrayInitializerExpr;
import com.github.javaparser.ast.expr.AssignExpr;
import com.github.javaparser.ast.expr.BinaryExpr;
import com.github.javaparser.ast.expr.BooleanLiteralExpr;
import com.github.javaparser.ast.expr.CastExpr;
import com.github.javaparser.ast.expr.CharLiteralExpr;
import com.github.javaparser.ast.expr.ClassExpr;
import com.github.javaparser.ast.expr.ConditionalExpr;
import com.github.javaparser.ast.expr.DoubleLiteralExpr;
import com.github.javaparser.ast.expr.EnclosedExpr;
import com.github.javaparser.ast.expr.Expression;
import com.github.javaparser.ast.expr.FieldAccessExpr;
import com.github.javaparser.ast.expr.InstanceOfExpr;
import com.github.javaparser.ast.expr.IntegerLiteralExpr;
import com.github.javaparser.ast.expr.LiteralExpr;
import com.github.javaparser.ast.expr.LongLiteralExpr;
import com.github.javaparser.ast.expr.MethodCallExpr;
import com.github.javaparser.ast.expr.NameExpr;
import com.github.javaparser.ast.expr.NullLiteralExpr;
import com.github.javaparser.ast.expr.ObjectCreationExpr;
import com.github.javaparser.ast.expr.StringLiteralExpr;
import com.github.javaparser.ast.expr.SuperExpr;
import com.github.javaparser.ast.expr.ThisExpr;
import com.github.javaparser.ast.expr.UnaryExpr;
import com.github.javaparser.ast.expr.VariableDeclarationExpr;
import com.github.javaparser.ast.stmt.BlockStmt;
import com.github.javaparser.ast.stmt.BreakStmt;
import com.github.javaparser.ast.stmt.ContinueStmt;
import com.github.javaparser.ast.stmt.DoStmt;
import com.github.javaparser.ast.stmt.EmptyStmt;
import com.github.javaparser.ast.stmt.ExplicitConstructorInvocationStmt;
import com.github.javaparser.ast.stmt.ExpressionStmt;
import com.github.javaparser.ast.stmt.ForStmt;
import com.github.javaparser.ast.stmt.ForeachStmt;
import com.github.javaparser.ast.stmt.IfStmt;
import com.github.javaparser.ast.stmt.LabeledStmt;
import com.github.javaparser.ast.stmt.ReturnStmt;
import com.github.javaparser.ast.stmt.Statement;
import com.github.javaparser.ast.stmt.SwitchEntryStmt;
import com.github.javaparser.ast.stmt.SwitchStmt;
import com.github.javaparser.ast.stmt.SynchronizedStmt;
import com.github.javaparser.ast.stmt.ThrowStmt;
import com.github.javaparser.ast.stmt.TryStmt;
import com.github.javaparser.ast.stmt.TypeDeclarationStmt;
import com.github.javaparser.ast.stmt.WhileStmt;
import com.github.javaparser.ast.type.ClassOrInterfaceType;
import com.github.javaparser.ast.type.Type;
import java.io.FileInputStream;
import java.io.InputStream;
import java.util.Map;

/**
 *
 * @author lweingart
 */
public class Parser {

    /**
     * class tag
     */
    public static final String TAG = "Parser class::";

    /**
     * The structure representing the parsed source file
     */
    private CompilationUnit unit;

    /**
     * The parsed source file
     */
    private final String sourceFile;

    /**
     * Classes in the source file
     */
    private final List<ClassDecl> classes;

   /**
    * Interfaces in the source file
    */
    private final List<InterfaceDecl> interfaces;

    /**
     * Enums in the source file
     */
    private final List<EnumDecl> enums;

    /**
     * Package name
     */
    private String packageName;

    /**
     * Imports of the source file
     */
    private List<String> imports;

    /**
     * Number of lines of source code
     */
    private int loc;

    /**
     * Indicates if a source file is a package-info file
     */
    private boolean packageInfo;

    /**
     * Contains informations about the package
     */
    private List<String> packageInfoComments;

    /**
     * The stream to parse
     */
    private InputStream inputStream;

    /**
     *
     * @param fileToParse
     *      path to the source file to parse
     * @param is
     *      can be null
     * @throws FileNotFoundException
     */
    public Parser(String fileToParse, InputStream is) throws FileNotFoundException {
        this.sourceFile = fileToParse;
        this.inputStream = (is != null) ? is : new FileInputStream(this.sourceFile);
        this.classes = new ArrayList<>();
        this.interfaces = new ArrayList<>();
        this.enums = new ArrayList<>();
        this.packageInfoComments = new ArrayList<>();
        this.packageInfo = false;
    }

    /**
     *
     * @throws FileNotFoundException
     * @throws ParseException
     */
    public void parse() throws FileNotFoundException, ParseException {
        InputStream in = this.inputStream;

        // boolean arg in parse() is for taknig the javadoc into account
        // sadly it is not supported yet by the github javaparser
        this.unit = JavaParser.parse(in, "UTF-8", true);

        PackageDeclaration pkg = this.unit.getPackage();
        if (pkg != null && pkg.getName() != null) {
            this.packageName = pkg.getName().getName();
            this.loc++;
        }

        if (this.unit.getImports() != null) {
            this.imports = new ArrayList<>();
            for (ImportDeclaration imp : this.unit.getImports()) {
                this.imports.add(imp.getName().toString());
            }
        }

        // for each class in the source file
        if (this.unit.getTypes() != null) {
            for (TypeDeclaration td : this.unit.getTypes()) {
                if (td instanceof ClassOrInterfaceDeclaration) {
                    ClassOrInterfaceDeclaration classOrInt = (ClassOrInterfaceDeclaration) td;
                    if (classOrInt.isInterface()) {
                        this.interfaces.add(parseInterface(td));
                    } else {
                        this.classes.add(parseClass(td));
                    }
                    this.loc++;
                } else if (td instanceof EnumDeclaration) {
                    this.enums.add(parseEnum(td));
                    this.loc++;
                }
            }
        } else {
            // TODO: parse package-info.java
            // not done yet because of the non support of javadoc by github javaparser.parser
            if (this.unit.getAllContainedComments() != null) {
                for (Comment com : this.unit.getAllContainedComments()) {
                    this.packageInfoComments.addAll(ParserUtils.prepareComments(com.getContent()));
                }
            }
            this.packageInfo = true;
        }
    }

    /**
     *
     * @return
     *      the structure representing the parsed source file
     */
    public CompilationUnit getUnit() {
        return this.unit;
    }

    /**
     *
     * @return
     *      the package name
     */
    public String getPackageName() {
        return this.packageName;
    }

    /**
     *
     * @return
     *      the list of imports
     */
    public List<String> getImports() {
        return this.imports;
    }

    /**
     *
     * @return
     *      the list of classes
     */
    public List<ClassDecl> getClasses() {
        return this.classes.isEmpty() ? null : this.classes;
    }

    /**
     *
     * @return
     *      the list of interfaces
     */
    public List<InterfaceDecl> getInterfaces() {
        return this.interfaces.isEmpty() ? null : this.interfaces;
    }

    /**
     *
     * @return
     *      the list of enums
     */
    public List<EnumDecl> getEnums() {
        return this.enums.isEmpty() ? null : this.enums;
    }

    /**
     *
     * @return
     *      the number of source code lines
     */
    public int getLoc() {
        return this.loc;
    }

    /**
     *
     * @return
     *      the boolean packageInfo
     */
    public boolean isPackageInfo() {
        return this.packageInfo;
    }

    public List<String> getPackageInfoComments() {
        return this.packageInfoComments;
    }

    /**
     *
     * @param td
     *      the parsed class
     * @return
     *      a class structure
     */
    private ClassDecl parseClass(TypeDeclaration td) {
        ClassDecl cls = new ClassDecl();

        cls.setName(td.getName());
        cls.setVisibility(visibility(td.getModifiers()));
        ClassRef superClass = new ClassRef();
        if (td instanceof ClassOrInterfaceDeclaration) {
            ClassOrInterfaceDeclaration coid = (ClassOrInterfaceDeclaration) td;
            List<ClassRef> extensions = new ArrayList<>();
            if (coid.getExtends() != null) {
                for (ClassOrInterfaceType obj : coid.getExtends()) {
                    superClass.setClassName(obj.toString());
                    extensions.add(superClass);
                }
            }
            cls.setExtendedClasses(extensions);
            List<InterfaceRef> implems = new ArrayList<>();
            if (coid.getImplements() != null) {
                for (ClassOrInterfaceType obj : coid.getImplements()) {
                    InterfaceRef implInterface = new InterfaceRef(null, obj.toString());
                    implems.add(implInterface);
                }
            }
            cls.setImplementedInterfaces(implems);

            if (coid.getComment() != null) {
                List<String> javadoc = ParserUtils.prepareComments(coid.getComment().getContent());

                cls.setDoc(javadoc);
            }

            // getJavadoc returns null, the parser lib does not support it yet
            if (coid.getJavaDoc() != null) {
                List<String> javadoc = ParserUtils.prepareComments(coid.getJavaDoc().getContent());
                cls.setDoc(javadoc);
            }
        }
        List<Attribute> attributes = new ArrayList<>();
        List<ConstructorDecl> constructors = new ArrayList<>();
        List<Method> methods = new ArrayList<>();

        if (td.getMembers() != null) {
            for (BodyDeclaration bodyElmt : td.getMembers()) {
                if (bodyElmt instanceof FieldDeclaration) { // attribute
                    FieldDeclaration fd = (FieldDeclaration) bodyElmt;
                    attributes.addAll(parseAttribute(fd));
                    this.loc++;
                } else if (bodyElmt instanceof ConstructorDeclaration) { // constructor
                    ConstructorDeclaration cd = (ConstructorDeclaration) bodyElmt;
                    constructors.add(parseConstructor(cd, attributes, superClass.getClassName()));
                    this.loc++;
                } else if (bodyElmt instanceof MethodDeclaration) { // method
                    MethodDeclaration md = (MethodDeclaration) bodyElmt;
                    methods.add(parseMethod(md, attributes));
                    this.loc++;
                } else if (bodyElmt instanceof ClassOrInterfaceDeclaration) {
                    this.classes.add(parseClass((TypeDeclaration) bodyElmt));
                    this.loc++;
                }
            }
        }
        cls.setAttributes(attributes);
        cls.setConstructors(constructors);
        cls.setMethods(methods);

        return cls;
    }

    /**
     *
     * @param td
     *      the parsed interface
     * @return
     *      an interface structure
     */
    private InterfaceDecl parseInterface(TypeDeclaration td) {
        ClassOrInterfaceDeclaration classOrInt = (ClassOrInterfaceDeclaration) td;
        InterfaceDecl interfaceDecl = new InterfaceDecl();
        // getJavadoc returns null, the parser lib does not support it yet
        if (classOrInt.getComment() != null) {
            List<String> javadoc = ParserUtils.prepareComments(classOrInt.getComment().getContent());
            interfaceDecl.setDoc(javadoc);
        }
        if (classOrInt.getJavaDoc() != null) {
            List<String> javadoc = ParserUtils.prepareComments(classOrInt.getJavaDoc().getContent());
            interfaceDecl.setDoc(javadoc);
        }
        interfaceDecl.setName(td.getName());
        List<InterfaceRef> implems = new ArrayList<>();
        if (classOrInt.getImplements() != null) {
            for (ClassOrInterfaceType obj : classOrInt.getImplements()) {
                InterfaceRef implInterface = new InterfaceRef(null, obj.toString());
                implems.add(implInterface);
            }
        }
        interfaceDecl.setImplementedInterfaces(implems);
        if (td.getMembers() != null) {
            List<ProtoDecl> prototypes = new ArrayList<>();
            for (BodyDeclaration bodyElmt : td.getMembers()) {
                if (bodyElmt instanceof MethodDeclaration) {
                    MethodDeclaration md = (MethodDeclaration) bodyElmt;
                    ProtoDecl proto = new ProtoDecl();

                    if (md.getComment() != null) {
                        List<String> javadoc = ParserUtils.prepareComments(md.getComment().getContent());
                        proto.setDoc(javadoc);
                    }
                    // getJavadoc returns null, the parser lib does not support it yet
                    if (md.getJavaDoc() != null) {
                        List<String> javadoc = ParserUtils.prepareComments(md.getJavaDoc().getContent());
                        proto.setDoc(javadoc);
                    }
                    proto.setName(new Ident(md.getName()));
                    FuncType type = new FuncType();
                    type.setParameters(parseParameters(md.getParameters()));
                    List<Field> results = new ArrayList<>();
                    List<String> coms = null;
                    if (md.getType().getComment() != null) {
                        coms = ParserUtils.prepareComments(md.getType().getComment().getContent());
                    }
                    results.add(new Field(coms, null, md.getType().toString()));
                    type.setResults(results);
                    proto.setType(type);
                    proto.setVisibility(visibility(md.getModifiers()));
                    prototypes.add(proto);
                    this.loc++;
                }
            }
            interfaceDecl.setPrototypes(prototypes);
        }
        interfaceDecl.setVisibility(visibility(td.getModifiers()));

        return interfaceDecl;
    }

    /**
     *
     * @param td
     *      the parsed enum
     * @return
     *      an enum structure
     */
    private EnumDecl parseEnum(TypeDeclaration td) {
        // enum
        EnumDeclaration enumeration = (EnumDeclaration) td;
        EnumDecl enumDecl = new EnumDecl();
        ClassRef superClass = new ClassRef();

        if (td.getComment()!= null) {
            List<String> javadoc = ParserUtils.prepareComments(td.getComment().getContent());
            enumDecl.setDoc(javadoc);
        }
        enumDecl.setName(td.getName());
        enumDecl.setVisibility(visibility(td.getModifiers()));
        if (enumeration.getEntries() != null) {
            List<Ident> entries = new ArrayList<>();
            for (Ident entry : entries) {
                Ident ident = new Ident(entry.getName());
                entries.add(ident);
                this.loc++;
            }
            enumDecl.setEnumConstants(entries);
        }
        List<InterfaceRef> implems = new ArrayList<>();
        if (enumeration.getImplements() != null) {
            for (ClassOrInterfaceType obj : enumeration.getImplements()) {
                InterfaceRef implInterface = new InterfaceRef(null, obj.toString());
                implems.add(implInterface);
            }
        }
        enumDecl.setImplementedInterfaces(implems);
        List<Attribute> attributes = new ArrayList<>();
        List<ConstructorDecl> constructors = new ArrayList<>();
        List<Method> methods = new ArrayList<>();

        if (td.getMembers() != null) {
            for (BodyDeclaration bodyElmt : td.getMembers()) {
                if (bodyElmt instanceof FieldDeclaration) { // attribute
                    FieldDeclaration fd = (FieldDeclaration) bodyElmt;
                    attributes.addAll(parseAttribute(fd));
                    this.loc++;
                } else if (bodyElmt instanceof ConstructorDeclaration) { // constructor
                    ConstructorDeclaration cd = (ConstructorDeclaration) bodyElmt;
                    constructors.add(parseConstructor(cd, attributes, superClass.getClassName()));
                    this.loc++;
                } else if (bodyElmt instanceof MethodDeclaration) { // method
                    MethodDeclaration md = (MethodDeclaration) bodyElmt;
                    methods.add(parseMethod(md, attributes));
                    this.loc++;
                } else if (bodyElmt instanceof ClassOrInterfaceDeclaration) {
                    this.classes.add(parseClass((TypeDeclaration) bodyElmt));
                    this.loc++;
                }
            }
        }
        enumDecl.setAttributes(attributes);
        enumDecl.setConstructors(constructors);
        enumDecl.setMethods(methods);

        return enumDecl;
    }

    /**
     *
     * @param visibility
     *      an int representing a visibility
     * @return
     *      a string representation of the visibility
     */
    private String visibility(int visibility) {
        List<String> result = new ArrayList<>();
        if (ModifierSet.hasModifier(visibility, ModifierSet.PUBLIC)) {
            return "public";
        }
        else if (ModifierSet.hasModifier(visibility, ModifierSet.PROTECTED)) {
            return "protected";
        }
        else if (ModifierSet.hasModifier(visibility, ModifierSet.PRIVATE)) {
            return "private";
        }

        return "package";
    }

    /**
     *
     * @param cd
     *      the github javaparser constructor declaration
     * @param attributes
     *      the list of attributes of the class,
     *      to potentially get a type from the name
     * @param superClass
     *      the name of the extended class
     * @return
     *      a constructor structure
     */
    private ConstructorDecl parseConstructor(ConstructorDeclaration cd, List<Attribute> attributes,
            String superClass) {
        ConstructorDecl constructor = new ConstructorDecl();

        if (cd.getComment() != null) {
            List<String> javadoc = ParserUtils.prepareComments(cd.getComment().getContent());
            constructor.setDoc(javadoc);
        }
        // getJavadoc returns null, the parser lib does not support it yet
        if (cd.getJavaDoc() != null) {
            List<String> javadoc = ParserUtils.prepareComments(cd.getJavaDoc().getContent());
            constructor.setDoc(javadoc);
        }
        constructor.setName(cd.getName());
        constructor.setParameters(parseParameters(cd.getParameters()));
        List<Stmt> body = new ArrayList<>();
        List<Statement> japaStmtsList = cd.getBlock().getStmts();
        if (japaStmtsList != null) {
            for (Statement statement : japaStmtsList) {

                // if statement is of the kind 'int foo = 42;' or 'this.bar = "bar";'
                if (statement instanceof ExpressionStmt) {
                    ExpressionStmt st = (ExpressionStmt) statement;
                    if (isAssignExpr(st.getExpression())) {
                        body.addAll(parseAssignExpression((AssignExpr) st.getExpression(), attributes, cd.getBeginLine()));
                    } else if (isVariableDeclarationExpr(st.getExpression())) {
                        body.add(parseVariableDeclarationExpression(
                                (VariableDeclarationExpr) st.getExpression(), attributes, cd.getBeginLine()));
                    } else {
                        ExprStatement exprStatement = new ExprStatement();
                        exprStatement.setExpression(parseExpression(st.getExpression(), attributes, cd.getBeginLine()));
                        body.add(exprStatement);
                        this.loc++;
                    }
                    // if statement == 'super(...);' or 'this(...);'
                } else if (statement instanceof ExplicitConstructorInvocationStmt) {
                    ExplicitConstructorInvocationStmt st = (ExplicitConstructorInvocationStmt) statement;
                    ExprStatement exprStatement = new ExprStatement();
                    CallExpression callExpr = new CallExpression();
                    FuncRef funcRef = new FuncRef();
                    if (superClass == null) {
                        superClass = "Object";
                    }
                    String funcName = st.isThis() ? cd.getName() : superClass;
                    funcRef.setFuncName(funcName);
                    callExpr.setFunction(funcRef);
                    List<Expr> arguments = new ArrayList<>();
                    if (st.getArgs() != null) {
                        for (Expression expression : st.getArgs()) {
                            arguments.add(parseExpression(expression, attributes, cd.getBeginLine()));
                        }
                    }
                    callExpr.setArguments(arguments);
                    callExpr.setLine(st.getBeginLine() - cd.getBeginLine());
                    exprStatement.setExpression(callExpr);
                    body.add(exprStatement);
                    this.loc++;
                } else {
                    body.addAll(parseStatement(statement, attributes, cd));
                    this.loc++;
                }
            }
        }
        constructor.setBody(body);
        constructor.setLoc(body.size());
        constructor.setVisibility(visibility(cd.getModifiers()));

        return constructor;
    }

    /**
     *
     * @param md
     *      the github javaparser method declaration
     * @param attributes
     *      the list of attributes of the class,
     *      to potentially get a type from the name
     * @return
     *      a method structure
     */
    private Method parseMethod(MethodDeclaration md, List<Attribute> attributes) {
        Method method = new Method();
        // getJavadoc returns null, the parser lib does not support it yet

        if (md.getComment() != null) {
            List<String> javadoc = ParserUtils.prepareComments(md.getComment().getContent());
            method.setDoc(javadoc);
        }
        if (md.getJavaDoc() != null) {
            List<String> javadoc = ParserUtils.prepareComments(md.getJavaDoc().getContent());
            method.setDoc(javadoc);
        }
        method.setName(md.getName());
        FuncType type = new FuncType();
        type.setParameters(parseParameters(md.getParameters()));
        List<Field> results = new ArrayList<>();
        List<String> coms = null;
        if (md.getType().getComment() != null) {
            coms = ParserUtils.prepareComments(md.getType().getComment().getContent());
        }
        results.add(new Field(coms, null, md.getType().toString()));
        type.setResults(results);
        method.setType(type);
        List<Stmt> body = new ArrayList<>();
        if (md.getBody() != null) {
            List<Statement> japaStmtsList = md.getBody().getStmts();
            if (japaStmtsList != null) {
                body.addAll(parseListStmt(japaStmtsList, attributes, md));
            }
        }

        method.setBody(body);
        method.setVisibility(visibility(md.getModifiers()));
        method.setOverride(false);
        if (md.getAnnotations() != null) {
            for (Expression expr : md.getAnnotations()) {
                if (expr.toString().equals("@Override"));
                method.setOverride(true);
            }
        }
        method.setLoc(body.size());

        return method;
    }

    /**
     *
     * @param japaStmtsList
     *      a list of statement (from github javaparser library)
     * @param attributes
     *      the list of attributes of the class,
     *      to potentially get a type from the name
     * @param bd
     *      a body declaration
     * @return
     *      a list of statements
     */
    private List<Stmt> parseListStmt(List<Statement> japaStmtsList, List<Attribute> attributes, BodyDeclaration bd) {
        List<Stmt> stmtList = new ArrayList<>();
        for (Statement statement : japaStmtsList) {
            stmtList.addAll(parseStatement(statement, attributes, bd));
        }

        return stmtList;
    }

    /**
     *
     * @param statement
     *      a statement (from the github javaparser library)
     * @param attributes
     *      the list of attributes of the class,
     *      to potentially get a type from the name
     * @param bd
     *      a body declaration
     * @return
     *      a list of statements
     */
    private List<Stmt> parseStatement(Statement statement, List<Attribute> attributes, BodyDeclaration bd) {
        List<Stmt> stmtList = new ArrayList<>();

        if (statement instanceof BlockStmt) {
            BlockStmt st = (BlockStmt) statement;
            if (st.getStmts() != null) {
                stmtList.addAll(parseListStmt(st.getStmts(), attributes, bd));
            }
        } else if (statement instanceof ExpressionStmt) {
            ExpressionStmt st = (ExpressionStmt) statement;
            if (isAssignExpr(st.getExpression())) {
                stmtList.addAll(parseAssignExpression((AssignExpr) st.getExpression(), attributes, bd.getBeginLine()));
            } else if (isVariableDeclarationExpr(st.getExpression())) {
                stmtList.add(parseVariableDeclarationExpression(
                        (VariableDeclarationExpr) st.getExpression(), attributes, bd.getBeginLine()));
            } else {
                ExprStatement expressionStatement = new ExprStatement();
                expressionStatement.setExpression(parseExpression(st.getExpression(), attributes, bd.getBeginLine()));
                stmtList.add(expressionStatement);
                this.loc++;
            }
        } else if (statement instanceof ForStmt) {
            ForStmt fs = (ForStmt) statement;
            LoopStatement loopStatement = parseForStatement(fs, attributes, bd);
            stmtList.add(loopStatement);
            this.loc++;
        } else if (statement instanceof ForeachStmt) {
            ForeachStmt fe = (ForeachStmt) statement;
            RangeLoopStatement rangeLoop = parseForeachStatement(fe, attributes, bd);
            stmtList.add(rangeLoop);
            this.loc++;
        } else if (statement instanceof WhileStmt) {
            WhileStmt ws = (WhileStmt) statement;
            LoopStatement ls = parseWhileStatement(ws, attributes, bd);
            stmtList.add(ls);
            this.loc++;
        } else if (statement instanceof BreakStmt) {
            OtherStatement os = parseBreakStatement(statement, bd);
            stmtList.add(os);
            this.loc++;
        } else if (statement instanceof ContinueStmt) {
            OtherStatement os = parseContinueStatement(statement, bd);
            stmtList.add(os);
            this.loc++;
        } else if (statement instanceof IfStmt) {
            IfStmt is = (IfStmt) statement;
            IfStatement ifSt = parseIfStatement(is, attributes, bd);
            stmtList.add(ifSt);
            this.loc++;
        } else if (statement instanceof ReturnStmt) {
            ReturnStmt rs = (ReturnStmt) statement;
            ReturnStatement returnStatement = parseReturnStatement(rs, attributes, bd);
            stmtList.add(returnStatement);
            this.loc++;
        } else if (statement instanceof TryStmt) {
            TryStmt ts = (TryStmt) statement;
            TryStatement tryStatement = parseTryStatement(ts, attributes, bd);
            stmtList.add(tryStatement);
            this.loc++;
        } else if (statement instanceof SwitchStmt) {
            SwitchStmt switchStmt = (SwitchStmt) statement;
            SwitchStatement switchStatement = parseSwitchStatement(switchStmt, attributes, bd);
            stmtList.add(switchStatement);
            this.loc++;
        } else if (statement instanceof ThrowStmt) {
            ThrowStmt throwStmt = (ThrowStmt) statement;
            ThrowStatement throwStatement = new ThrowStatement();
            throwStatement.setExpression(parseExpression(throwStmt.getExpr(), attributes, bd.getBeginLine()));
            stmtList.add(throwStatement);
            this.loc++;
        } else if (statement instanceof DoStmt) {
            DoStmt doStmt = (DoStmt) statement;
            LoopStatement loopStmt = parseDoStatement(doStmt, attributes, bd);
            stmtList.add(loopStmt);
            this.loc++;
        } else if (statement instanceof SynchronizedStmt) {
            SynchronizedStmt ss = (SynchronizedStmt) statement;
            stmtList.addAll(parseStatement(ss.getBlock(), attributes, bd));
            this.loc++;
        } else if (statement instanceof LabeledStmt) {
            LabeledStmt ls = (LabeledStmt) statement;
            stmtList.addAll(parseStatement(ls.getStmt(), attributes, bd));
            this.loc++;
        } else if (statement instanceof EmptyStmt || statement instanceof TypeDeclarationStmt) {
            // not implemented, does not fit in the API yet
            // info logs are currently disabled in order to get a clean output
            this.loc++;
//            Log.i(TAG, "The type of statement '"
//                    .concat(statement.getClass().toString())
//                    .concat("' is not managed by the parser. ")
//                    .concat("Statement:: ")
//                    .concat(statement.toString()));
        } else {
            // missed case ?());
            this.loc++;
//            Log.i(TAG, "The type of statement '"
//                    .concat(statement.getClass().toString())
//                    .concat("' is not managed by the parser")
//                    .concat("Statement:: ")
//                    .concat(statement.toString()));
        }

        return stmtList;
    }

    /**
     *
     * @param fs
     *      a github javaparser ForStatement
     * @param attributes
     *      the list of attributes of the class,
     *      to potentially get a type from the name
     * @param bd
     *      a body declaration
     * @return
     *  a loopStatement structure
     */
    private LoopStatement parseForStatement(ForStmt fs, List<Attribute> attributes, BodyDeclaration bd) {
        LoopStatement loopStatement = new LoopStatement();
        if (fs.getInit() != null) {
            List<Stmt> initialisation = new ArrayList<>();
            for (Expression exp : fs.getInit()) {
                if (exp instanceof VariableDeclarationExpr) {
                    initialisation.add(parseVariableDeclarationExpression(
                            (VariableDeclarationExpr) exp, attributes, bd.getBeginLine()));
                } else if (exp instanceof AssignExpr) {
                    initialisation.addAll(parseAssignExpression((AssignExpr) exp, attributes, bd.getBeginLine()));
                }
            }
            loopStatement.setInitialization(initialisation);
        }
        if (fs.getUpdate() != null) {
            List<Stmt> postItstmt = new ArrayList<>();
            for (Expression exp : fs.getUpdate()) {
                if (exp instanceof UnaryExpr) {
                    UnaryExpr unEx = (UnaryExpr) exp;
                    ExprStatement expStmt = new ExprStatement();
                    IncDecExpr incDec = new IncDecExpr();
                    String operator = unaryOperator(unEx.getOperator())[0];
                    boolean prefix = unaryOperator(unEx.getOperator())[1].equals("true");
                    switch (operator) {
                        case "DEC":
                        case "INC":
                            incDec.setOperator(operator);
                            incDec.setOperand(parseExpression(unEx.getExpr(), attributes, bd.getBeginLine()));
                            incDec.setPrefix(prefix);
                            expStmt.setExpression(incDec);
                            break;
                        default:
                            Log.e(TAG, "For loop :".concat("update expression not an").concat("IncdecExpression."));
                    }
                    postItstmt.add(expStmt);
                } else if (exp instanceof AssignExpr) {
                    AssignExpr assEx = (AssignExpr) exp;
                    postItstmt.addAll(parseAssignExpression(assEx, attributes, bd.getBeginLine()));
                }
            }
            loopStatement.setPostItStmt(postItstmt);
        }
        loopStatement.setBody(parseStatement(fs.getBody(), attributes, bd));
        loopStatement.setPostEvaluated(true);
        loopStatement.setLine(fs.getBeginLine() - bd.getBeginLine());

        return loopStatement;
    }

    /**
     *
     * @param fe
     *      a github javaparser ForeachStatement
     * @param attributes
     *      the list of attributes of the class,
     *      to potentially get a type from the name
     * @param bd
     *      a body declaration
     * @return
     *      a rangeLoopStatement
     */
    private RangeLoopStatement parseForeachStatement(ForeachStmt fe, List<Attribute> attributes, BodyDeclaration bd) {
        RangeLoopStatement rangeLoop = new RangeLoopStatement();
        List<Expr> variables = new ArrayList<>();
        variables.add(parseExpression(fe.getIterable(), attributes, bd.getBeginLine()));
        rangeLoop.setVariables(variables);
        rangeLoop.setIterable(new Ident(fe.getVariable().toString()));
        rangeLoop.setBody(parseStatement(fe.getBody(), attributes, bd));
        rangeLoop.setLine(fe.getBeginLine() - bd.getBeginLine());

        return rangeLoop;
    }

    /**
     *
     * @param ws
     *      a github javaparser whileStatement
     * @param attributes
     *      the list of attributes of the class,
     *      to potentially get a type from the name
     * @param bd
     *      a body declaration
     * @return
     *      a LoopStatement structure
     */
    private LoopStatement parseWhileStatement(WhileStmt ws, List<Attribute> attributes, BodyDeclaration bd) {
        LoopStatement ls = new LoopStatement();
        List<Stmt> initialisation = null;
        if (ws.getCondition() != null) {
            if (ws.getCondition() instanceof AssignExpr) {
                initialisation = new ArrayList<>();
                initialisation.addAll(parseAssignExpression((AssignExpr) ws.getCondition(), attributes, bd.getBeginLine()));
            } else {
                ls.setCondition(parseExpression(ws.getCondition(), attributes, bd.getBeginLine()));
            }
        }
        ls.setInitialization(initialisation);
        ls.setBody(parseStatement(ws.getBody(), attributes, bd));
        ls.setLine(ws.getBeginLine() - bd.getBeginLine());

        return ls;
    }

    /**
     *
     * @param statement
     *      a statement (in the case of a github javaparser BreakStatement)
     * @param bd
     *      a body declaration
     * @return
     *      an OtherStatement structure
     */
    private OtherStatement parseBreakStatement(Statement statement, BodyDeclaration bd) {
        OtherStatement os = new OtherStatement();
        List<Stmt> body = new ArrayList<>();
        ExprStatement expStmt = new ExprStatement(new Ident("break"));
        body.add(expStmt);
        os.setBody(body);
        os.setLine(statement.getBeginLine() - bd.getBeginLine());

        return os;
    }

    /**
     *
     * @param statement
     *      a statement (in the case of a github javaparser ContinueStatement
     * @param bd
     *      a body declaration
     * @return
     *      an OtherStatement structure
     */
    private OtherStatement parseContinueStatement(Statement statement, BodyDeclaration bd) {
        OtherStatement os = new OtherStatement();
        List<Stmt> body = new ArrayList<>();
        ExprStatement expStmt = new ExprStatement(new Ident("continue"));
        body.add(expStmt);
        os.setBody(body);
        os.setLine(statement.getBeginLine() - bd.getBeginLine());

        return os;
    }

    /**
     *
     * @param is
     *      a github javaparser IfStatement
     * @param attributes
     *      the list of attributes of the class,
     *      to potentially get a type from the name
     * @param bd
     *      a body declaration
     * @return
     *      an IfStatement structure
     */
    private IfStatement parseIfStatement(IfStmt is, List<Attribute> attributes, BodyDeclaration bd) {
        IfStatement ifSt = new IfStatement();
        List<Stmt> body = new ArrayList<>();
        ifSt.setCondition(parseExpression(is.getCondition(), attributes, bd.getBeginLine()));
        body.addAll(parseStatement(is.getThenStmt(), attributes, bd));
        if (is.getElseStmt() != null) {
            List<Stmt> elsStmts = new ArrayList<>();
            elsStmts.addAll(parseStatement(is.getElseStmt(), attributes, bd));
            ifSt.setElseStmt(elsStmts);
        }
        ifSt.setBody(body);
        ifSt.setLine(is.getBeginLine() - bd.getBeginLine());

        return ifSt;
    }

    /**
     *
     * @param rs
     *      a github javaparser ReturnStatement
     * @param attributes
     *      the list of attributes of the class,
     *      to potentially get a type from the name
     * @param bd
     *      a body declaration
     * @return
     *      a ReturnStatement structure
     */
    private ReturnStatement parseReturnStatement(ReturnStmt rs, List<Attribute> attributes, BodyDeclaration bd) {
        ReturnStatement returnStatement = new ReturnStatement();
        if (rs.getExpr() != null) {
            List<Expr> results = new ArrayList<>();
            results.add(parseExpression(rs.getExpr(), attributes, loc));
            returnStatement.setResults(results);
        }
        returnStatement.setLine(rs.getBeginLine() - bd.getBeginLine());

        return returnStatement;
    }

    /**
     *
     * @param ts
     *      a github javaparser TryStatement
     * @param attributes
     *      the list of attributes of the class,
     *      to potentially get a type from the name
     * @param bd
     *      a body declaration
     * @return
     *      a TryStatement structure
     */
    private TryStatement parseTryStatement(TryStmt ts, List<Attribute> attributes, BodyDeclaration bd) {
        TryStatement tryStatement = new TryStatement();
        tryStatement.setBody(parseStatement(ts.getTryBlock(), attributes, bd));
        if (ts.getCatchs() != null) {
            List<CatchClause> catchClauses = new ArrayList<>();
            List<com.github.javaparser.ast.stmt.CatchClause> japaCatchClauses = ts.getCatchs();
            for (com.github.javaparser.ast.stmt.CatchClause japaCatchClause : japaCatchClauses) {
                CatchClause catchClause = new CatchClause();
                List<Field> params = new ArrayList<>();
                if (japaCatchClause.getExcept().getTypes() != null) {
                    for (Type t : japaCatchClause.getExcept().getTypes()) {
                        String type = ParserUtils.parseException(t.toString()).get("type");
                        String name = ParserUtils.parseException(t.toString()).get("name");
                        List<String> doc = null;
                        if (t.getComment() != null) {
                            doc = ParserUtils.prepareComments(t.getComment().getContent());
                        }
                        Field f = new Field(doc, name, type);
                        params.add(f);
                    }
                }
                catchClause.setParameters(params);
                List<Stmt> body = new ArrayList<>();
                body.addAll(parseStatement(japaCatchClause.getCatchBlock(), attributes, bd));
                catchClause.setBody(body);
                catchClauses.add(catchClause);
                this.loc++;
            }
            tryStatement.setCatchClauses(catchClauses);
        }
        if (ts.getFinallyBlock() != null) {
            tryStatement.setFinallyStmts(parseStatement(ts.getFinallyBlock(), attributes, bd));
        }

        return tryStatement;
    }

    /**
     *
     * @param switchStmt
     *      a github javaparser SwithStatement
     * @param attributes
     *      the list of attributes of the class,
     *      to potentially get a type from the name
     * @param bd
     *      a body declaration
     * @return
     *      a SwitchStatement structure
     */
    private SwitchStatement parseSwitchStatement(SwitchStmt switchStmt, List<Attribute> attributes, BodyDeclaration bd) {
        SwitchStatement switchStatement = new SwitchStatement();
        switchStatement.setCondition(parseExpression(switchStmt.getSelector(), attributes, bd.getBeginLine()));
        if (switchStmt.getEntries() != null) {
            List<CaseClause> caseClauses = new ArrayList<>();
            for (SwitchEntryStmt switchEntry : switchStmt.getEntries()) {
                CaseClause caseClause = new CaseClause();
                List<Expr> conditions = new ArrayList<>();
                if (switchEntry.getLabel() != null) {
                    conditions.add(parseExpression(switchEntry.getLabel(), attributes, bd.getBeginLine()));
                    caseClause.setConditions(conditions);
                } else {
                    if (switchEntry.getStmts() != null) {
                        switchStatement.setDefaultStmt(parseListStmt(switchEntry.getStmts(), attributes, bd));
                    }
                }
                if (switchEntry.getStmts() != null) {
                    List<Stmt> body = new ArrayList<>();
                    body.addAll(parseListStmt(switchEntry.getStmts(), attributes, bd));
                    caseClause.setBody(body);
                }
                caseClauses.add(caseClause);
            }
            switchStatement.setCaseClauses(caseClauses);
        }
        return switchStatement;
    }

    /**
     *
     * @param doStmt
     *      a github javaparser DoStatement
     * @param attributes
     *      the list of attributes of the class,
     *      to potentially get a type from the name
     * @param bd
     *      a body declaration
     * @return
     *      a LoopStatement structure
     */
    private LoopStatement parseDoStatement(DoStmt doStmt, List<Attribute> attributes, BodyDeclaration bd) {
        LoopStatement loopStmt = new LoopStatement();
        loopStmt.setCondition(parseExpression(doStmt.getCondition(), attributes, bd.getBeginLine()));
        List<Stmt> list = new ArrayList<>();
        list.addAll(parseStatement(doStmt.getBody(), attributes, bd));
        loopStmt.setBody(list);
        loopStmt.setPostEvaluated(true);
        loopStmt.setLine(doStmt.getBeginLine() - bd.getBeginLine());
        return loopStmt;
    }

    /**
     *
     * @param params
     *      list of parameters to parse
     * @return
     *      a list of Field structures
     */
    private List<Field> parseParameters(List<Parameter> params) {
        List<Field> fields = new ArrayList<>();
        if (params != null) {
            for (Parameter parameter : params) {
                Field field = new Field();
                field.setName(parameter.getId().getName());
                field.setType(parameter.getType().toString());
                if (parameter.getComment() != null) {
                    List<String> doc = ParserUtils.prepareComments(parameter.getComment().getContent());
                    field.setDoc(doc);
                }
                fields.add(field);
            }
        }

        return fields;
    }

    /**
     *
     * @param fd
     *      a github javaparser FieldDeclaration
     * @return
     *      an Attribute structure
     */
    private List<Attribute> parseAttribute(FieldDeclaration fd) {
        List<Attribute> attributes = new ArrayList<>();
        List<VariableDeclarator> vars = fd.getVariables();
        if (vars.size() == 1) {
            Attribute attribute = new Attribute();
            attribute.setType(fd.getType().toString());
            attribute.setName(vars.get(0).getId().getName());
            if (vars.get(0).getInit() != null) {
                attribute.setValue(vars.get(0).getInit().toString());
            }
            // getJavadoc returns null, the parser lib does not support it yet
            if (fd.getComment()!= null) {
                List<String> javadoc = ParserUtils.prepareComments(fd.getComment().getContent());
                attribute.setDoc(javadoc);
            }
            if (fd.getJavaDoc() != null) {
                List<String> javadoc = ParserUtils.prepareComments(fd.getJavaDoc().getContent());
                attribute.setDoc(javadoc);
            }
            String visibility = visibility(fd.getModifiers());
            attribute.setVisibility(visibility);
            attribute.setConstant(ModifierSet.isFinal(fd.getModifiers()));
            attribute.setStatick(ModifierSet.isStatic(fd.getModifiers()));
            attributes.add(attribute);
        } else {
            for (VariableDeclarator var : vars) {
                Attribute attribute = new Attribute();
                attribute.setType(fd.getType().toString());
                attribute.setName(var.getId().getName());
                attribute.setValue(ParserUtils.extractValue(fd.toString()));
                String visibility = visibility(fd.getModifiers());
                attribute.setVisibility(visibility);
                attribute.setConstant(ModifierSet.isFinal(fd.getModifiers()));
                attribute.setStatick(ModifierSet.isStatic(fd.getModifiers()));
                attributes.add(attribute);
            }
        }

        return attributes;
    }

    /**
     *
     * @param expression
     *  a github javaparser Expression
     * @param attributes
     *      the list of attributes of the class,
     *      to potentially get a type from the name
     * @param lineNumber
     *      the starting line number of the parse method or constructor
     * @return
     *      an Expression structure
     */
    public Expr parseExpression(Expression expression, List<Attribute> attributes, int lineNumber) {
        if (expression instanceof AssignExpr) { // this.bar = "bar";
            AssignExpr assExpr = (AssignExpr) expression;
            return parseExpression(assExpr.getTarget(), attributes, lineNumber);
        } else if (expression instanceof MethodCallExpr) {
            MethodCallExpr mcEx = (MethodCallExpr) expression;
            return parseMethodCallExpression(mcEx, attributes, lineNumber);
        } else if (expression instanceof NameExpr) { // ident
            NameExpr nEx = (NameExpr) expression;
            Ident nameExpr = new Ident(nEx.getName());
            return nameExpr;
        } else if (expression instanceof LiteralExpr) { // basic lit
            return parseLiteralExpr((LiteralExpr) expression);
        } else if (expression instanceof FieldAccessExpr) {
            FieldAccessExpr fieldExpr = (FieldAccessExpr) expression;
            Ident ident = new Ident(ParserUtils.parseTarget(expression.toString()).get("name"));
            AttributeRef attrRef = new AttributeRef(ident);
            return attrRef;
        } else if (expression instanceof ObjectCreationExpr) {
            ObjectCreationExpr objConExpr = (ObjectCreationExpr) expression;
            return parseObjectCreationExpression(objConExpr, attributes, lineNumber);
        } else if (expression instanceof ArrayAccessExpr) {
            ArrayAccessExpr arryExpr = (ArrayAccessExpr) expression;
            return parseArrayAccessExpression(arryExpr);
        } else if (expression instanceof UnaryExpr) {
            UnaryExpr unExpr = (UnaryExpr) expression;
            return parseUnaryExpression(unExpr, attributes, lineNumber);
        } else if (expression instanceof ConditionalExpr) {
            ConditionalExpr condExpr = (ConditionalExpr) expression;
            return parseConditionalExpression(condExpr, attributes, lineNumber);
        } else if (expression instanceof CastExpr) {
            CastExpr castExpr = (CastExpr) expression;
            return parseExpression(castExpr.getExpr(), attributes, lineNumber);
        } else if (expression instanceof BinaryExpr) {
            BinaryExpr binEx = (BinaryExpr) expression;
            return parseBinaryExpression(binEx, attributes, lineNumber);
        } else if (expression instanceof EnclosedExpr) {
            EnclosedExpr enclosedExpr = (EnclosedExpr) expression;
            return parseExpression(enclosedExpr.getInner(), attributes, lineNumber);
        } else if (expression instanceof InstanceOfExpr) {
            InstanceOfExpr intExpr = (InstanceOfExpr) expression;
            return parseInstanceOfExpression(intExpr, attributes, lineNumber);
        } else if (expression instanceof ArrayCreationExpr) {
            ArrayCreationExpr arryCreaExpr = (ArrayCreationExpr) expression;
            return parseArrayCreationExpression(arryCreaExpr, attributes, lineNumber);
        } else if (expression instanceof ArrayInitializerExpr) {
            ArrayInitializerExpr arryInEx = (ArrayInitializerExpr) expression;
            return parseArrayInitializerExpression(arryInEx, attributes, lineNumber);
        } else if (expression instanceof ThisExpr) {
            return new Ident("this");
        } else if (expression instanceof SuperExpr) {
            return new Ident("super");
        } else if (expression instanceof ClassExpr) {
            return new Ident(expression.toString());
        } else if (expression instanceof VariableDeclarationExpr) { // int foo = 42;
            // should be parsed by parseVariableDeclarationExpression()
            Log.e(TAG, "Unreachable case :: expression : ".concat(expression.toString()));
            return null;
        } else {
            Log.e(TAG, "The type of expression '".concat(expression.getClass().toString()).concat("' is not managed by the parser"));
            return null;
        }
    }

    /**
     *
     * @param arryExpr
     *      a github javaparser ArrayExpression
     * @return
     *      an IndexExpression structure
     */
    private IndexExpression parseArrayAccessExpression(ArrayAccessExpr arryExpr) {
        IndexExpression indExpr = new IndexExpression();
        Ident expr = new Ident(arryExpr.getName().toString());
        Ident index = new Ident(arryExpr.getIndex().toString());
        indExpr.setExpression(expr);
        indExpr.setIndex(index);
        return indExpr;
    }

    /**
     *
     * @param unExpr
     *      a github javaparser UnaryExpression
     * @param attributes
     *      the list of attributes of the class,
     *      to potentially get a type from the name
     * @param lineNumber
     *      the starting line number of the parse method or constructor
     * @return
     *      an Expression structure
     */
    private Expr parseUnaryExpression(UnaryExpr unExpr, List<Attribute> attributes, int lineNumber) {
        String operator = (unaryOperator(unExpr.getOperator())[0]);
        boolean prefix = unaryOperator(unExpr.getOperator())[1].equals("true");
        switch (operator) {
            case "DEC":
            case "INC": {
                IncDecExpr incDecExpr = new IncDecExpr();
                incDecExpr.setOperator(operator);
                incDecExpr.setOperand(parseExpression(unExpr.getExpr(), attributes, lineNumber));
                incDecExpr.setPrefix(prefix);
                return incDecExpr;
            }
            default:
                UnaryExpression unaryExpression = new UnaryExpression();
                unaryExpression.setOperator(unaryOperator(unExpr.getOperator())[0]);
                unaryExpression.setOperand(parseExpression(unExpr.getExpr(), attributes, lineNumber));
                return unaryExpression;
        }
    }

    /**
     *
     * @param condExpr
     *      a github javaparser ConditionalExpression
     * @param attributes
     *      the list of attributes of the class,
     *      to potentially get a type from the name
     * @param lineNumber
     *      the starting line number of the parse method or constructor
     * @return
     *      a TernaryExpression structure
     */
    private TernaryExpression parseConditionalExpression(ConditionalExpr condExpr, List<Attribute> attributes, int lineNumber) {
        TernaryExpression ternaryExpr = new TernaryExpression();
        ternaryExpr.setCondition(parseExpression(condExpr.getCondition(), attributes, lineNumber));
        ternaryExpr.setThen(parseExpression(condExpr.getThenExpr(), attributes, lineNumber));

        if (condExpr.getElseExpr() != null) {
            ternaryExpr.setEls(parseExpression(condExpr.getElseExpr(), attributes, lineNumber));
        }
        return ternaryExpr;
    }

    /**
     *
     * @param binEx
     *      a github javaparser BinaryExpression
     * @param attributes
     *      the list of attributes of the class,
     *      to potentially get a type from the name
     * @param lineNumber
     *      the starting line number of the parse method or constructor
     * @return
     *      a BinaryExpression structure
     */
    private BinaryExpression parseBinaryExpression(BinaryExpr binEx, List<Attribute> attributes, int lineNumber) {
        BinaryExpression binaryExpression = new BinaryExpression();
        binaryExpression.setLeftExpr(parseExpression(binEx.getLeft(), attributes, lineNumber));
        binaryExpression.setOperator(binaryOperator(binEx.getOperator()));
        binaryExpression.setRightExpr(parseExpression(binEx.getRight(), attributes, lineNumber));
        return binaryExpression;
    }

    /**
     *
     * @param intExpr
     *      a github javaparserInstanceOfExpression
     * @param attributes
     *      the list of attributes of the class,
     *      to potentially get a type from the name
     * @param lineNumber
     *      the starting line number of the parse method or constructor
     * @return
     *      a BinaryExpression structure
     */
    private BinaryExpression parseInstanceOfExpression(InstanceOfExpr intExpr, List<Attribute> attributes, int lineNumber) {
        BinaryExpression binayExpr = new BinaryExpression();
        Ident ident = new Ident(intExpr.getType().toString());
        binayExpr.setLeftExpr(ident);
        binayExpr.setOperator("TYPE_EQ");
        binayExpr.setRightExpr(parseExpression(intExpr.getExpr(), attributes, lineNumber));
        return binayExpr;
    }

    /**
     *
     * @param arryCreaExpr
     *      a github javaparser ArrayCreationExpression
     * @param attributes
     *      the list of attributes of the class,
     *      to potentially get a type from the name
     * @param lineNumber
     *      the starting line number of the parse method or constructor
     * @return
     *      an Expression structure
     * @throws NumberFormatException
     */
    private Expr parseArrayCreationExpression(ArrayCreationExpr arryCreaExpr, List<Attribute> attributes, int lineNumber)
                                              throws NumberFormatException
    {
        ArrayExpression arryExpr = new ArrayExpression();
        if (arryCreaExpr.getInitializer() != null) {
            ArrayLit aryLit = new ArrayLit();
            ArrayType aryType = new ArrayType();
            aryType.setElementType(new Ident(arryCreaExpr.getType().toString()));
            List<Integer> dimensions = new ArrayList<>();
            if (arryCreaExpr.getDimensions() != null) {
                for (Expression exp : arryCreaExpr.getDimensions()) {
                    dimensions.add(Integer.parseInt(exp.toString()));
                }
            } else {
                if (arryCreaExpr.getInitializer().getValues() != null) {
                    dimensions.addAll(ParserUtils.extractDimensions(arryCreaExpr.getInitializer().getValues()));
                }
            }
            aryType.setDimensions(dimensions);
            aryLit.setType(aryType);
            List<Expr> elements = new ArrayList<>();
            if (arryCreaExpr.getInitializer().getValues() != null) {
                for (Expression expr : arryCreaExpr.getInitializer().getValues()) {
                    elements.add(parseExpression(expr, attributes, lineNumber));
                }
            }
            aryLit.setElements(elements);
            return aryLit;
        }
        ArrayType arrayType = new ArrayType();
        arrayType.setElementType(new Ident(arryCreaExpr.getType().toString()));
        List<Integer> dimensions = new ArrayList<>();
        for (Expression exp : arryCreaExpr.getDimensions()) {
            try {
                dimensions.add(Integer.parseInt(exp.toString()));
            } catch (NumberFormatException e) {
                dimensions.add(-1);
            }
        }
        arrayType.setDimensions(dimensions);
        arryExpr.setType(arrayType);
        return arryExpr;
    }

    /**
     *
     * @param arryInEx
     *      a github javaparser ArrayInitializerExpression
     * @param attributes
     *      the list of attributes of the class,
     *      to potentially get a type from the name
     * @param lineNumber
     *      the starting line number of the parse method or constructor
     * @return
     *      an Array literal structure
     */
    private ArrayLit parseArrayInitializerExpression(ArrayInitializerExpr arryInEx, List<Attribute> attributes, int lineNumber) {
        ArrayLit arryLit = new ArrayLit();
        ArrayType arrayType = new ArrayType();
        List<Integer> dimensions = new ArrayList<>();
        if (arryInEx.getValues() != null) {
            dimensions.addAll(ParserUtils.extractDimensions(arryInEx.getValues()));
        }
        arrayType.setDimensions(dimensions);
        arryLit.setType(arrayType);
        List<Expr> elements = new ArrayList<>();
        if (arryInEx.getValues() != null) {
            for (Expression expr : arryInEx.getValues()) {
                elements.add(parseExpression(expr, attributes, lineNumber));
            }
        }
        arryLit.setElements(elements);
        return arryLit;
    }

    /**
     *
     * @param mcEx
     *      a github javaparser MethodCallExpression
     * @param attributes
     *      the list of attributes of the class,
     *      to potentially get a type from the name
     * @param lineNumber
     *      the starting line number of the parse method or constructor
     * @return
     *      a CallExpression structure
     */
    private CallExpression parseMethodCallExpression(MethodCallExpr mcEx, List<Attribute> attributes, int lineNumber) {
        CallExpression callExpr = new CallExpression();
        FuncRef funcRef = new FuncRef();
        funcRef.setFuncName(mcEx.getName());
        if (mcEx.getScope() != null) {
            funcRef.setNamespace(mcEx.getScope().toString());
        }
        callExpr.setFunction(funcRef);
        List<Expr> arguments = new ArrayList<>();
        if (mcEx.getArgs() != null) {
            for (Expression expr : mcEx.getArgs()) {
                arguments.add(parseExpression(expr, attributes, lineNumber));
            }
        }
        callExpr.setArguments(arguments);
        callExpr.setLine(mcEx.getBeginLine() - lineNumber);
        return callExpr;
    }

    /**
     *
     * @param objConExpr
     *      a github javaparser ObjectCreationExpression
     * @param attributes
     *      the list of attributes of the class,
     *      to potentially get a type from the name
     * @param lineNumber
     *      the starting line number of the parse method or constructor
     * @return
     *      a ConstructorCallExpression structure
     */
    private ConstructorCallExpr parseObjectCreationExpression(ObjectCreationExpr objConExpr, List<Attribute> attributes, int lineNumber) {
        ConstructorCallExpr constr = new ConstructorCallExpr();
        FuncRef funcRef = new FuncRef();
        funcRef.setNamespace(objConExpr.getType().getName());
        funcRef.setFuncName(objConExpr.getType().getName());
        constr.setFunction(funcRef);
        List<Expr> arguments = new ArrayList<>();
        if (objConExpr.getArgs() != null) {
            for (Expression exp : objConExpr.getArgs()) {
                arguments.add(parseExpression(exp, attributes, lineNumber));
            }
        }
        constr.setArguments(arguments);
        constr.setLine(objConExpr.getBeginLine() - lineNumber);
        return constr;
    }

    /**
     *
     * @param expression
     *      a github javaparser Expression
     * @return
     *      a boolean depending on if it's a github javaparser AssignExpression
     */
    private boolean isAssignExpr(Expression expression) {
        return (expression instanceof AssignExpr);
    }

    /**
     *
     * @param expression
     *  a github javaparser Expression
     * @return
     *      a boolean depending on if it's a github javaparser
     *      VariableDeclarationExpression
     */
    private boolean isVariableDeclarationExpr(Expression expression) {
        return (expression instanceof VariableDeclarationExpr);
    }

    /**
     *
     * @param assignExpr
     *      a github javaparser AssignExpression
     * @param attributes
     *      the list of attributes of the class,
     *      to potentially get a type from the name
     * @param lineNumber
     * @return
     */
    private List<AssignStatement> parseAssignExpression(AssignExpr assignExpr, List<Attribute> attributes, int lineNumber) {
        List<AssignStatement> assignStatements = new ArrayList<>();

        AssignStatement assignStmt = new AssignStatement();

        Map<String, String> leftPart = ParserUtils.parseTarget(assignExpr.getTarget().toString());
        List<Expr> leftHandExprList = new ArrayList<>();
        Ident ident = new Ident();
        ident.setName(leftPart.get("name"));
        if (leftPart.get("this") == null) {
            leftHandExprList.add(ident);
        } else {
            AttributeRef attributeRef = new AttributeRef();
            attributeRef.setName(ident);
            leftHandExprList.add(attributeRef);
        }
        assignStmt.setLeftHandExpr(leftHandExprList);

        List<Expr> rightHandExprList = new ArrayList<>();
        if (assignExpr.getValue() instanceof AssignExpr) {
            AssignExpr assign = (AssignExpr) assignExpr.getValue();
            Map<String, String> lp = ParserUtils.parseTarget(assign.getTarget().toString());
            List<Expr> leftH2 = new ArrayList<>();
            Ident id = new Ident(lp.get("name"));
            if (leftPart.get("this") == null) {
                leftH2.add(ident);
            } else {
                AttributeRef attributeRef = new AttributeRef();
                attributeRef.setName(ident);
                leftH2.add(attributeRef);
            }
            assignStmt.setRightHandExpr(leftH2);
            assignStatements.add(assignStmt);
            assignStatements.addAll(parseAssignExpression((AssignExpr) assignExpr.getValue(), attributes, lineNumber));
        } else {
            rightHandExprList.add(parseExpression(assignExpr.getValue(), attributes, lineNumber));
            assignStmt.setRightHandExpr(rightHandExprList);
            assignStmt.setLine(assignExpr.getBeginLine() - lineNumber);
            assignStatements.add(assignStmt);
        }
        this.loc++;
        return assignStatements;
    }

    /**
     *
     * @param varDecExpr
     * @param attributes
     *      the list of attributes of the class,
     *      to potentially get a type from the name
     * @param lineNumber
     *      the starting line number of the parse method or constructor
     * @return
     *      a DeclarationStatement structure
     */
    private DeclStatement parseVariableDeclarationExpression(VariableDeclarationExpr varDecExpr, List<Attribute> attributes, int lineNumber) {
        DeclStatement declStatement = new DeclStatement();

        List<Expr> leftHandExprList = new ArrayList<>();
        List<Expr> rightHandExprList = new ArrayList<>();
        if (varDecExpr.getVars().size() == 1) {
            ValueSpec valueSpec = new ValueSpec();
            Ident name = new Ident(varDecExpr.getVars().get(0).getId().getName());
            Ident type = new Ident(varDecExpr.getType().toString());
            valueSpec.setName(name);
            valueSpec.setType(type);
            leftHandExprList.add(valueSpec);
            declStatement.setKind(type.getName());

            if (varDecExpr.getVars().get(0).getInit() != null) {
                rightHandExprList.add(parseExpression(varDecExpr.getVars().get(0).getInit(), attributes, lineNumber));
            }
            declStatement.setLine(varDecExpr.getBeginLine() - lineNumber);
        } else {
            for (int i = 0; i < varDecExpr.getVars().size(); i++) {
                ValueSpec valueSpec = new ValueSpec();
                Ident name = new Ident(varDecExpr.getVars().get(i).getId().getName());
                Ident type = new Ident(varDecExpr.getType().toString());
                valueSpec.setName(name);
                valueSpec.setType(type);
                leftHandExprList.add(valueSpec);
                declStatement.setKind(type.getName());

                if (varDecExpr.getVars().get(0).getInit() != null) {
                    rightHandExprList.add(parseExpression(varDecExpr.getVars().get(0).getInit(), attributes, lineNumber));
                }
                declStatement.setLine(varDecExpr.getBeginLine() - lineNumber);
            }
        }
        declStatement.setLeftHandExpr(leftHandExprList);
        declStatement.setRightHandExpr(rightHandExprList);
        declStatement.setLine(varDecExpr.getBeginLine() - lineNumber);
        this.loc++;
        return declStatement;
    }

    /**
     *
     * @param literalExpr
     *      a github javaparser LiteralExpression
     * @return
     *      an Basic literal structure
     */
    private BasicLit parseLiteralExpr(LiteralExpr literalExpr) {
        if (literalExpr instanceof BooleanLiteralExpr) {
            BooleanLiteralExpr boolLit = (BooleanLiteralExpr) literalExpr;
            BasicLit basicLit = new BasicLit("BOOL", String.valueOf(boolLit.getValue()));
            return basicLit;
        } else if (literalExpr instanceof NullLiteralExpr) {
            NullLiteralExpr nullLit = (NullLiteralExpr) literalExpr;
            BasicLit basicLit = new BasicLit("NIL", "null");
            return basicLit;
        } else if (literalExpr instanceof IntegerLiteralExpr) {
            IntegerLiteralExpr intLit = (IntegerLiteralExpr) literalExpr;
            BasicLit basicLit = new BasicLit("INT", intLit.getValue());
            return basicLit;
        } else if (literalExpr instanceof LongLiteralExpr) {
            LongLiteralExpr longLit = (LongLiteralExpr) literalExpr;
            BasicLit basicLit = new BasicLit("LONG", longLit.getValue());
            return basicLit;
        } else if (literalExpr instanceof DoubleLiteralExpr) {
            DoubleLiteralExpr doubleLit = (DoubleLiteralExpr) literalExpr;
            BasicLit basicLit = new BasicLit("DOUBLE", doubleLit.getValue());
            return basicLit;
        } else if (literalExpr instanceof CharLiteralExpr) {
            CharLiteralExpr charLit = (CharLiteralExpr) literalExpr;
            BasicLit basicLit = new BasicLit("CHAR", charLit.getValue());
            return basicLit;
        } else {
            StringLiteralExpr stringLit = (StringLiteralExpr) literalExpr;
            BasicLit basicLit = new BasicLit("STRING", stringLit.getValue());
            return basicLit;
        }

    }

    /**
     *
     * @param operator
     *      a github javaparser unary operator
     * @return
     *      a string array with a string representation of the operator
     *      and a string representation of a boolean if the operator is
     *      a post or a pre increment
     */
    private String[] unaryOperator(UnaryExpr.Operator operator) {
        String[] result = new String[2];
        // boolean string value is only used to differentiate preValue from posValue
        switch (operator) {
            case inverse:
                result[0] = "NOT";
                result[1] = "true";
                return result;
            case negative:
                result[0] = "NEG";
                result[1] = "true";
                return result;
            case not:
                result[0] = "NOT";
                result[1] = "true";
                return result;
            case posDecrement:
                result[0] = "DEC";
                result[1] = "false";
                return result;
            case posIncrement:
                result[0] = "INC";
                result[1] = "false";
                return result;
            case positive:
                result[0] = "POS";
                result[1] = "true";
                return result;
            case preDecrement:
                result[0] = "DEC";
                result[1] = "true";
                return result;
            case preIncrement:
                result[0] = "INC";
                result[1] = "true";
                return result;
        }
        return result;
    }

    /**
     *
     * @param operator
     *      a github javaparser binay operator
     * @return
     *      a string representation of the operator
     */
    private String binaryOperator(BinaryExpr.Operator operator) {
        switch (operator) {
            case and:
                return "LAND";
            case binAnd:
                return "AND";
            case binOr:
                return "OR";
            case divide:
                return "QUO";
            case equals:
                return "EQ";
            case greater:
                return "GTR";
            case greaterEquals:
                return "GEQ";
            case less:
                return "LSS";
            case lessEquals:
                return "LEQ";
            case lShift:
                return "SHIFT_LEFT";
            case minus:
                return "SUB";
            case notEquals:
                return "NEQ";
            case or:
                return "LOR";
            case plus:
                return "ADD";
            case remainder:
                return "MOD";
            case rSignedShift:
                return "SHIFT_RIGHT";
            case rUnsignedShift:
                return "SHIFT_RIGHT";
            case times:
                return "MUL";
            default:
                return "XOR";
        }
    }
}
