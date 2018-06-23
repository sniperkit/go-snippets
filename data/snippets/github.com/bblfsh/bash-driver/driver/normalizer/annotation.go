package normalizer

import (
	"github.com/bblfsh/bash-driver/driver/normalizer/intellij"

	"gopkg.in/bblfsh/sdk.v1/uast"
	. "gopkg.in/bblfsh/sdk.v1/uast/ann"
	"gopkg.in/bblfsh/sdk.v1/uast/transformer"
	"gopkg.in/bblfsh/sdk.v1/uast/transformer/annotatter"
	"gopkg.in/bblfsh/sdk.v1/uast/transformer/positioner"
	"gopkg.in/src-d/go-errors.v0"
)

// Transformers is the of list `transformer.Transfomer` to apply to a UAST, to
// learn more about the Transformers and the available ones take a look to:
// https://godoc.org/gopkg.in/bblfsh/sdk.v1/uast/transformers
var Transformers = []transformer.Tranformer{
	annotatter.NewAnnotatter(AnnotationRules),
	positioner.NewFillLineColFromOffset(),
}

var ErrRootMustBeFile = errors.NewKind("root must have internal type FILE")

var conditionOperators = On(intellij.ConditionOperator).Roles(uast.Operator).Children(
	On(HasToken("-eq")).Roles(uast.Relational, uast.Equal),
	On(HasToken("-ne")).Roles(uast.Relational, uast.Not, uast.Equal),
	On(HasToken("-gt")).Roles(uast.Relational, uast.GreaterThan),
	On(HasToken("-ge")).Roles(uast.Relational, uast.GreaterThanOrEqual),
	On(HasToken("-lt")).Roles(uast.Relational, uast.LessThan),
	On(HasToken("-le")).Roles(uast.Relational, uast.LessThanOrEqual),
)

// XXX expressions/statements

// AnnotationRules describes how a UAST should be annotated with `uast.Role`.
//
// https://godoc.org/gopkg.in/bblfsh/sdk.v1/uast/ann
var AnnotationRules = On(Any).Self(
	On(Not(intellij.File)).Error(ErrRootMustBeFile.New()),
	On(intellij.File).Roles(uast.File).Descendants(
		On(intellij.Comment).Roles(uast.Comment, uast.Noop),
		On(intellij.LineFeed).Roles(uast.Whitespace, uast.Noop),
		On(intellij.Whitespace).Roles(uast.Whitespace, uast.Noop),
		On(intellij.IntLiteral).Roles(uast.Number, uast.Literal, uast.Primitive),

		On(Or(intellij.VarUse, intellij.Variable)).Roles(uast.Expression, uast.Variable,
			uast.Identifier),
		On(Or(intellij.ComposedVariable, intellij.VarSubstitution)).Roles(uast.Variable,
			uast.Expression, uast.Identifier, uast.Incomplete),

		On(intellij.UnEvalString).Roles(uast.Expression, uast.String, uast.Literal),
		On(intellij.StringBegin).Roles(uast.Expression, uast.String, uast.Block),
		On(intellij.StringEnd).Roles(uast.Expression, uast.String, uast.Incomplete),
		On(intellij.StringContent).Roles(uast.Expression, uast.String, uast.Literal),
		On(intellij.String).Roles(uast.Expression, uast.String, uast.Literal, uast.Block),
		// FIXME: needs uast node "Execute" or "Shell" or similar
		On(intellij.BackQuoteShellCommand).Roles(uast.Expression, uast.String, uast.Literal,
			uast.Call, uast.Incomplete),
		On(Or(intellij.SubShellCmd, intellij.PipelineCmd)).Roles(uast.Expression, uast.Call,
			uast.Incomplete),
		On(Or(intellij.Shebang, intellij.ShebangElement)).Roles(uast.Comment, uast.Pathname,
			uast.Incomplete),
		// variable declaration
		On(intellij.SimpleCommand).Roles(uast.Expression).Children(
			On(intellij.VarDefElement).Roles(uast.Expression, uast.Assignment, uast.Binary).Children(
				On(intellij.AssignmentWord).Roles(uast.Identifier, uast.Left),
				On(intellij.OperatorAssign).Roles(uast.Operator, uast.Assignment),
				On(And(Not(intellij.AssignmentWord), Not(intellij.OperatorAssign))).Roles(uast.Right),
			),
		),
		// function declaration
		On(intellij.FunctionDefElement).Roles(uast.Function, uast.Declaration, uast.Block).Children(
			On(intellij.Function).Roles(uast.Function, uast.Declaration),
			On(intellij.NamedSymbol).Roles(uast.Function, uast.Declaration, uast.Name),
			On(intellij.GroupElement).Roles(uast.Function, uast.Declaration, uast.Body, uast.Block),
		),
		// let statement / expression (unfortunately it doesnt produce a subtree...)
		On(intellij.LetExpression).Roles(uast.Expression, uast.Assignment, uast.Incomplete),
		On(intellij.LetStatement).Roles(uast.Statement, uast.Incomplete),

		// if statement
		conditionOperators,
		On(intellij.ArithmeticCmd).Roles(uast.Expression, uast.Arithmetic, uast.Incomplete),
		On(intellij.ArithmeticSimple).Roles(uast.Expression, uast.Arithmetic, uast.Incomplete),
		On(Or(intellij.OperatorArithLess, intellij.OperatorLess)).Roles(uast.Operator, uast.Relational, uast.LessThan),
		On(Or(intellij.OperatorArithMore, intellij.OperatorMore)).Roles(uast.Operator, uast.Relational, uast.GreaterThan),
		On(intellij.OperatorArithEqual).Roles(uast.Operator, uast.Relational, uast.Equal),
		On(intellij.OperatorArithNotEqual).Roles(uast.Operator, uast.Relational, uast.Not, uast.Equal),
		On(intellij.OperatorNotEqual).Roles(uast.Operator, uast.Relational, uast.Not, uast.Equal),
		On(intellij.OperatorLessEqual).Roles(uast.Operator, uast.Relational, uast.LessThanOrEqual),
		On(intellij.OperatorMoreEqual).Roles(uast.Operator, uast.Relational, uast.GreaterThanOrEqual),
		On(intellij.OperatorBoolOr).Roles(uast.Operator, uast.Boolean, uast.Or),
		On(intellij.OperatorBoolAnd).Roles(uast.Operator, uast.Boolean, uast.And),
		On(intellij.OperatorBoolNot).Roles(uast.Operator, uast.Boolean, uast.Not),
		On(intellij.If).Roles(uast.Statement, uast.If),
		On(intellij.ConditionalShellCommand).Roles(uast.Expression, uast.Condition),
		On(intellij.IfShellCommand).Roles(uast.If, uast.Expression, uast.Block).Children(
			On(intellij.ElIf).Roles(uast.Statement, uast.If, uast.Else, uast.Incomplete),
			On(intellij.Else).Roles(uast.Statement, uast.Else, uast.Incomplete),
			On(intellij.SimpleCommand).Roles(uast.Expression, uast.If, uast.Condition),
		),
		On(intellij.Then).Roles(uast.Statement, uast.Incomplete),
		On(intellij.LogicalBlock).Roles(uast.Expression, uast.If, uast.Then),
		// FIXME: no role in the uast for "end"
		On(intellij.Fi).Roles(uast.Statement, uast.Incomplete),
		// for statement
		On(intellij.Do).Roles(uast.Statement, uast.Block),
		On(intellij.In).Roles(uast.Expression, uast.Binary, uast.Operator, uast.Relational, uast.Contains),
		On(intellij.GenericBash).Roles(uast.Incomplete).Children(
			On(HasToken("break")).Roles(uast.Statement, uast.Break),
			On(HasToken("continue")).Roles(uast.Statement, uast.Continue),
		),
		On(intellij.Done).Roles(uast.Statement, uast.Incomplete),

		// These are more tokens that real semantic nodes, but they're in the AST tree
		// so we must tag them.
		On(intellij.BracketArithLeft).Roles(uast.Incomplete),
		On(intellij.BracketArithRight).Roles(uast.Incomplete),
		On(intellij.LeftConditional).Roles(uast.Incomplete),
		On(intellij.RightConditional).Roles(uast.Incomplete),
		On(intellij.SemiColon).Roles(uast.Incomplete),
		On(intellij.BraceOpen).Roles(uast.Incomplete),
		On(intellij.BraceClose).Roles(uast.Incomplete),
		On(intellij.SemiColon).Roles(uast.Incomplete),
		On(intellij.LeftSquare).Roles(uast.Incomplete),
		On(intellij.RightSquare).Roles(uast.Incomplete),
		On(intellij.DoubleSemiColon).Roles(uast.Incomplete),
		On(intellij.LeftBracket).Roles(uast.Incomplete),
		On(intellij.RightBracket).Roles(uast.Incomplete),
		On(intellij.ParenOpen).Roles(uast.Incomplete),
		On(intellij.DoubleParenOpen).Roles(uast.Incomplete),
		On(intellij.ParenClose).Roles(uast.Incomplete),
		On(intellij.DoubleParenClose).Roles(uast.Incomplete),
		On(intellij.BackQuote).Roles(uast.Incomplete),
		On(intellij.Dollar).Roles(uast.Incomplete),
		On(intellij.ErrorElement).Roles(uast.Incomplete),
		On(intellij.RedirectList).Roles(uast.Incomplete),
		On(intellij.RedirectElement).Roles(uast.Incomplete),

		On(intellij.FileDescriptors).Roles(uast.Identifier, uast.Receiver, uast.Incomplete),
		On(intellij.Word).Roles(uast.Expression, uast.Identifier),
		On(intellij.CombinedWord).Roles(uast.Expression, uast.String, uast.Identifier, uast.Incomplete),
		On(intellij.ForShellCommand).Roles(uast.For, uast.Statement).Children(
			On(intellij.For).Roles(uast.Statement, uast.Incomplete),
			On(intellij.VarDefElement).Roles(uast.Expression, uast.For, uast.Iterator),
			On(intellij.LogicalBlock).Roles(uast.Expression, uast.For, uast.Body),
			On(intellij.CombinedWord).Roles(uast.Expression, uast.For, uast.Update),
		),
		// while and until statement
		On(Or(intellij.WhileLoop, intellij.While)).Roles(uast.Statement, uast.While),
		On(Or(intellij.UntilLoop, intellij.Until)).Roles(uast.Statement, uast.While, uast.Incomplete),
		On(Or(intellij.WhileLoop, intellij.UntilLoop)).Children(
			On(intellij.LogicalBlock).Roles(uast.Expression, uast.While, uast.Body),
			On(intellij.ConditionalShellCommand).Roles(uast.Expression, uast.While, uast.Condition).Children(
				On(intellij.LeftConditional).Roles(uast.Incomplete),
				On(intellij.RightConditional).Roles(uast.Incomplete),
				// "-a" and such
				On(intellij.ConditionOperator).Roles(uast.Operator, uast.Incomplete),
			),
		),
		On(intellij.OperatorEqual).Roles(uast.Expression, uast.Relational, uast.Operator, uast.Equal),
		On(intellij.OperatorBoolOr).Roles(uast.Expression, uast.Binary, uast.Operator, uast.Boolean, uast.Or),
		On(intellij.OperatorBoolAnd).Roles(uast.Expression, uast.Binary, uast.Operator, uast.Boolean, uast.And),
		On(intellij.ParameterExpOperatorA).Roles(uast.Operator, uast.Incomplete),
		On(intellij.ParameterExpOperatorH).Roles(uast.Operator, uast.Incomplete),
		On(intellij.ParameterExpOperatorP).Roles(uast.Operator, uast.Incomplete),
		On(intellij.ParameterExpOperatorC).Roles(uast.Operator, uast.Incomplete),
		On(intellij.ParameterExpOperatorS).Roles(uast.Operator, uast.Incomplete),
		On(intellij.AppendOperator).Roles(uast.Operator, uast.Incomplete),

		// case pattern
		On(Or(intellij.Case, intellij.CasePattern2)).Roles(uast.Expression, uast.Case),
		On(intellij.CasePatternList).Roles(uast.Expression, uast.Case, uast.List).Children(
			On(intellij.LogicalBlock).Roles(uast.Expression, uast.Case, uast.Body),
		),
		On(intellij.CasePattern).Roles(uast.Statement, uast.Switch),
		On(intellij.CaseEnd).Roles(uast.Statement, uast.Incomplete),

		// source command (import-ish)
		On(intellij.IncludeCommand).Roles(uast.Expression, uast.Import).Children(
			On(HasToken("source")).Roles(uast.Statement, uast.Import),
			On(intellij.SourceFileReference).Roles(uast.Expression, uast.Import, uast.Pathname, uast.Identifier),
		),
	),
)
