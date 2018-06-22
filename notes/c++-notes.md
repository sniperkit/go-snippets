---
layout: page
title: "C++ notes for self"
description: "Scattered notes about C/C++ stuff."
category : notes
published: true
---


* [SFINAE "Substitution Failure Is Not An Error (SFINAE)" idiom](https://en.wikibooks.org/wiki/More_C%2B%2B_Idioms/SFINAE)

## Articles out there

* [2012: Why should I have written ZeroMQ in C, not C++ (part I) – Martin Sústrik](http://250bpm.com/blog:4)


## IDEA(s)

### Class Dummy {...} :

Have a small CMake-based project with dummy "example" codes of interest, like connecting to databases, setting things up, etc.

* Focus on both 1) how to configure the CMake (sub)project/s, and 2) the examples, which should be _practical_ (no "school case" foo/bar stuff).

The idea is that it builds what it can, wits. what was found, like for ex. if `libpqxx` was found then build the PostgreSql database examples, likewise for MySql, etc.

* Boost.log
* PostgreSql libpqxx
* MySql: mysql-connector-c++

### Auto-rebuild upon file change (from a shell)

Have a command line tool that would monitor source code changes in the current
project dir., and :

* would run syntax checking on the updated sources ;
* re-build the affected targets.
* or rebuild the whole project ?
* or just rebuild the target that was specified on the command line.

## FYI

* [Removing Deprecated Exception Specifications from C++17 (2015-09-28)](http://www.open-std.org/jtc1/sc22/wg21/docs/papers/2015/p0003r0.html#1.0)
  — `throw(some_excetion_type)` may no longer appear as part of a function
  prototype. — it is "replaced" by `noexcept`
  – <http://stackoverflow.com/a/13841791/643087>
  – <http://en.cppreference.com/w/cpp/language/except_spec>

## Questions one asks oneself

* Now I have to find a stack impl. -> which? For ex. `std::vector` vs `llvm::SmallVector`?

### About concepts

* Placement ::new idiom.

## References

* [Containers library @ cppreference.com](http://en.cppreference.com/w/cpp/container)

  Includes a nice "member function table" for quick look-up of available methods
  per container -- as well as C++11 additions.


## e/Books

* [Optimizing C++ (Wikibooks.org)](https://en.wikibooks.org/wiki/Optimizing_C%2B%2B)

### Misc.

* <https://conan.io/> sort of a C/C++ package manager, free for open source.

#### Parsing C++

Is not a trivial task at all, many have tried.

* [Parsing C++ (March 2001, Andrew Birkett)](http://www.nobugs.org/developer/parsingcpp/)

### Idioms, tips'n'tricks

### Class inheritance, true polymorphic classes, compile-time polymorphism, ...

* A pure virtual class method, ex. `virtual doWork() =0;` _may_ provide (have) a
  definition.  Client code still have to implement the interface, but may rely on
  some sort of "default implementation" provided by its base class(es).


#### Assertions

```cpp
assert(1 == 2 && "Dude, this ain't no good, like not at all.");
```

##### static_assert()

* `decltype()` & `std::is_same<X, Y>`
    - Defensive programming in presence of `auto ...`, for ex. :

      ```cpp
      for(const auto& row : rset) {
        static_assert( std::is_same<decltype(row), pqxx::row const&>::value,
          "PQXX: (!!) we're expecting pqxx::row& from result sets.");
          ...
          ...
      }
      ```

#### explicit vs. implicit constructors

Mark implicit "converting" constructors with a comment like `/* implicit */` so as to assert intention.

```cpp
class StringRef {
public:
  ...

  /// Construct an empty string ref.
  /* implicit * / StringRef() = default;

  ...
};
```

#### Compile-time polymorphism

```cpp
LLVM/Clang has a nice example of compile-time polymorphism whereby no dynamic dispatch occurs at runtime.  The wiring is achieved upon compilation.

class PimplIOVisitor :
        public RecursiveASTVisitor< PimplIOVisitor > {

  Derived &getDerived() { return *static_cast<Derived *>(this); }

  ...
  ...
};
```

#### do { ... } while(false)

Seen that a few times: a `do { ... } while( false );` that wraps some code. What is it usefull for ?

* it does some sort of lexical scoping.
* names shadowing (e.g. defined variables at the parent lexical scope): one may redefine a variable with the same name (and any different type).
* automatic destruction of stuff ?

This one example is from LLVM/Clang's `RecursiveASTVisitor< Derived >` :

```cpp
// A helper macro to implement short-circuiting when recursing.  It
// invokes CALL_EXPR, which must be a method call, on the derived
// object (s.t. a user of RecursiveASTVisitor can override the method
// in CALL_EXPR).
#define TRY_TO(CALL_EXPR)       \
do {                            \
  if (!getDerived().CALL_EXPR)  \
    return false;               \
} while (false)
// Note:       ^ No trailing ';'
```

### initializer lists `{...}`

#### wrt. STL containers

Before :

```bash
auto pair = Artifacts_.insert( std::make_pair(K, Artifact(D)) );
```

After :

```bash
auto pair = Artifacts_.insert( {K, Artifact(D)} );
```

#### Preprocessor

```cpp
// This macro makes available a variable D, the passed-in decl.
#define DEF_TRAVERSE_DECL(DECL, CODE)                              \
  template <typename Derived>                                      \
  bool RecursiveASTVisitor<Derived>::Traverse##DECL(DECL *D) {     \
    bool ShouldVisitChildren = true;                               \
    bool ReturnValue = true;                                       \
    if (!getDerived().shouldTraversePostOrder())                   \
      TRY_TO(WalkUpFrom##DECL(D));                                 \
    { CODE; }                                                      \
    if (ReturnValue && ShouldVisitChildren)                        \
      TRY_TO(TraverseDeclContextHelper(dyn_cast<DeclContext>(D))); \
    if (ReturnValue && getDerived().shouldTraversePostOrder())     \
      TRY_TO(WalkUpFrom##DECL(D));                                 \
    return ReturnValue;                                            \
  }
```

```cpp
DEF_TRAVERSE_DECL(EmptyDecl, {})
DEF_TRAVERSE_DECL(AnotherEmptyDecl, { /* No code here. */ })

DEF_TRAVERSE_DECL(LabelDecl, {// There is no code in a LabelDecl.
                             })

DEF_TRAVERSE_DECL(StaticAssertDecl, {
  TRY_TO(TraverseStmt(D->getAssertExpr()));
  TRY_TO(TraverseStmt(D->getMessage()));
})

DEF_TRAVERSE_DECL(
    TranslationUnitDecl,
    {// Code in an unnamed namespace shows up automatically in
     // decls_begin()/decls_end().  Thus we don't need to recurse on
     // D->getAnonymousNamespace().
    })
```



#### Elements of taste and style

* Popular coding styles out there : Google, LLVM, LLDB, etc.
    - Clang-format ?

* Constness :
    - `pqxx::row const &` vs `const pqxx::row &` ?
    - `try {...} catch (const std::exception &e) { ... }`
    vs `try {...} catch (std::exception const &e) { ... }` ?

## Other

List __pkg-config__ packages :

```bash
 $ pkg-config --list-all | sort -k1,1 | less
```

## Pointers

* <https://wiki.gentoo.org/wiki/GCC_optimization>

__EOF__
