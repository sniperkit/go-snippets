---
layout: page
title: "LLVM/Clang API pointers"
description: "API documentation pointers, diagrams, notes."
category : notes
tags: [wide, draft, C++, Clang]
published: true
---

* __[clang.llvm.org/doxygen/](http://clang.llvm.org/doxygen/index.html)
  &ndash; API doc.__
* [clang.llvm.org/docs/IntroductionToTheClangAST.html](https://clang.llvm.org/docs/IntroductionToTheClangAST.html)
* [clang.llvm.org/docs/ClangPlugins.html](https://clang.llvm.org/docs/ClangPlugins.html)
* ~~<https://github.com/fabic/llvm-clang>~~
* <https://railsware.com/blog/2014/02/28/creation-and-using-clang-plugin-with-xcode/>

<table class="fs-smaller pre">
  <thead>
    <tr>
      <th>What</th>
      <th>Method / attributes</th>
      <th>&ndash;</th>
      <th>&ndash;</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>ASTContext</td>
      <td>QualType ASTContext::getTagDeclType (const TagDecl *Decl) const;</td>
      <td>TagDecl * Type::getAsTagDecl () const;</td>
      <td></td>
    </tr>
    <tr>
      <td>ASTContext</td>
      <td>QualType ASTContext::getTypeDeclType (
        const TypeDecl *Decl,
        const TypeDecl *PrevDecl=nullptr) const;</td>
      </td>
      <td>Type (as QualType) &ndash;given&ndash; TypeDecl</td>
      <td></td>
    </tr>
    <tr>
      <td>Type</td>
      <td>template&lt; typename T &gt;<br>const T * Type::getAs&lt;T&gt;() const;</td>
      <td></td>
      <td></td>
    </tr>
    <tr>
      <td>Type</td>
      <td>CXXRecordDecl * getAsCXXRecordDecl () const;</td>
      <td></td>
      <td></td>
    </tr>
    <tr>
      <td>Type</td>
      <td>const CXXRecordDecl * getPointeeCXXRecordDecl () const;</td>
      <td></td>
      <td></td>
    </tr>
    <tr>
      <td>Decl</td>
      <td></td>
      <td></td>
      <td></td>
    </tr>
    <tr>
      <td>NamedDecl</td>
      <td></td>
      <td></td>
      <td></td>
    </tr>
    <tr>
      <td>TypeDecl</td>
      <td></td>
      <td></td>
      <td></td>
    </tr>
    <tr>
      <td>TagDecl</td>
      <td></td>
      <td></td>
      <td></td>
    </tr>
    <tr>
      <td>ValueDecl</td>
      <td></td>
      <td></td>
      <td></td>
    </tr>
    <tr>
      <td></td>
      <td></td>
      <td></td>
      <td></td>
    </tr>
    <tr>
      <td></td>
      <td></td>
      <td></td>
      <td></td>
    </tr>
    <tr>
      <td></td>
      <td></td>
      <td></td>
      <td></td>
    </tr>
    <!--
    <tr>
      <td></td>
      <td></td>
      <td></td>
      <td></td>
    </tr>
    -->
  </tbody>
</table>

### List all Decl nodes

```bash
cat <<EOF | clang++ -E -x c++ - | sed -e '/^$/d;/^#/d' | sort -k2 | column -t -s :
#define DECL(CLASS, BASE) CLASS  :  BASE
#include "clang/AST/DeclNodes.inc"
EOF
```

More convenient with a one-liner :

```bash
echo -e '#define DECL(CLASS, BASE) CLASS : BASE\n#include "clang/AST/DeclNodes.inc"\n' \
  | clang++ -E -x c++ -  \
  | sed -e '/^$/d;/^#/d' \
  | sort -k2             \
  | column -t -s :
```

Likewise for the `Type` hierarchy :

```bash
echo -e '#define TYPE(CLASS, BASE) CLASS : BASE\n#include "clang/AST/TypeNodes.def"\n'  \
  | clang++ -E -x c++ -  \
  | sed -e '/^$/d;/^#/d' \
  | sort -k2             \
  | column -t -s :
```

The __BuiltinType__/s :

TODO


And the `Stmt` hierarchy :

```bash
echo -e '#define STMT(CLASS, BASE) CLASS : BASE\n#include "clang/AST/StmtNodes.inc"\n'  \
  | clang++ -E -x c++ -  \
  | sed -e '/^$/d;/^#/d' \
  | sort -k2             \
  | column -t -s :
```

* [TypeSourceInfo][TSI]
* [QualType][QualType]

## Pointers

### [ASTContext](http://clang.llvm.org/doxygen/classclang_1_1ASTContext.html)

  <img src="http://clang.llvm.org/doxygen/classclang_1_1ASTContext__coll__graph.png" />

### [Type](http://clang.llvm.org/doxygen/classclang_1_1Type.html)

  <img src="http://clang.llvm.org/doxygen/classclang_1_1Type__inherit__graph.png" />

See also [TypeSourceInfo][TSI] that we typically obtain from __[Decl/s][Decl]__

### [Decl][Decl]

  <img src="http://clang.llvm.org/doxygen/classclang_1_1Decl__inherit__graph.png" />

### [NamedDecl][NamedDecl]

  <img src="http://clang.llvm.org/doxygen/classclang_1_1NamedDecl__inherit__graph.png" />

* [IdentifierInfo](http://clang.llvm.org/doxygen/classclang_1_1IdentifierInfo.html) * getIdentifier () const;

### [TypeDecl][TypeDecl]

  <img src="http://clang.llvm.org/doxygen/classclang_1_1TypeDecl__inherit__graph.png" />

### [ValueDecl][ValueDecl]

  <img src="http://clang.llvm.org/doxygen/classclang_1_1ValueDecl__inherit__graph.png" />

### [DeclaratorDecl][DeclaratorDecl]

  <img src="http://clang.llvm.org/doxygen/classclang_1_1DeclaratorDecl__inherit__graph.png" />

* [TypeSourceInfo][TypeSourceInfo] * DeclaratorDecl::getTypeSourceInfo () const;
* SourceRange     DeclaratorDecl::getSourceRange () const override LLVM_READONLY;
* SourceLocation  DeclaratorDecl::getLocStart () const LLVM_READONLY;
* SourceLocation  DeclaratorDecl::getTypeSpecStartLoc () const;

#### [DeclContext](http://clang.llvm.org/doxygen/classclang_1_1DeclContext.html)

  <img src="http://clang.llvm.org/doxygen/classclang_1_1DeclContext__inherit__graph.png" />

##### [TagDecl][TagDecl]

  <img src="http://clang.llvm.org/doxygen/classclang_1_1TagDecl__inherit__graph.png" />

##### [FunctionDecl](http://clang.llvm.org/doxygen/classclang_1_1FunctionDecl.html)


  <img src="http://clang.llvm.org/doxygen/classclang_1_1FunctionDecl__inherit__graph.png" />

* <tt>[QualType](http://clang.llvm.org/doxygen/classclang_1_1QualType.html) [getReturnType()](http://clang.llvm.org/doxygen/classclang_1_1FunctionDecl.html#a2935d37297596c57aa383de101a55030)</tt>
* `SourceRange  getReturnTypeSourceRange () const`


[TSI]: http://clang.llvm.org/doxygen/classclang_1_1TypeSourceInfo.html "TypeSourceInfo [Decl.h]"
[Decl]: http://clang.llvm.org/doxygen/classclang_1_1Decl.html "Decl"
[NamedDecl]: http://clang.llvm.org/doxygen/classclang_1_1NamedDecl.html "NamedDecl < Decl"
[TypeDecl]: http://clang.llvm.org/doxygen/classclang_1_1TypeDecl.html "TypeDecl < NamedDecl < Decl"
[ValueDecl]: http://clang.llvm.org/doxygen/classclang_1_1ValueDecl.html "ValueDecl < NamedDecl < Decl"
[DeclaratorDecl]: http://clang.llvm.org/doxygen/classclang_1_1DeclaratorDecl.html "DeclaratorDecl < ValueDecl < NamedDecl < Decl"
[QualType]: http://clang.llvm.org/doxygen/classclang_1_1QualType.html "QualType is-nearly-a Type (smart-pointer-like behaviour)"
[TagDecl]: http://clang.llvm.org/doxygen/classclang_1_1TagDecl.html "TagDecl < TypeDecl < NamedDecl < Decl"
