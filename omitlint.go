// Package omitlint defines an Analyzer that checks for fields with a json
// 'omitempty' option, but which have no meaningful zero value, meaning they
// will never be omitted.
package omitlint

import (
	"flag"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"reflect"
	"slices"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

type omitlint struct{}

// NewAnalyzer returns a new omitlint analyzer.
func NewAnalyzer() *analysis.Analyzer {
	omitlint := &omitlint{}

	a := &analysis.Analyzer{
		Name:     "omitlint",
		Doc:      "Checks for fields with a json 'omitempty' option, but which have no meaningful zero value, meaning they will never be omitted.",
		Requires: []*analysis.Analyzer{inspect.Analyzer},
		Run:      omitlint.run,
	}
	// false, 0, a nil pointer, a nil interface value, and any empty array, slice, map, or string
	a.Flags.Init("omitlint", flag.ExitOnError)
	a.Flags.Var(versionFlag{}, "V", "print version and exit")

	return a
}

func (o *omitlint) run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{
		(*ast.StructType)(nil),
	}
	inspect.Preorder(nodeFilter, func(node ast.Node) {
		var s *ast.StructType
		var ok bool
		if s, ok = node.(*ast.StructType); !ok {
			return
		}
		if tv, ok := pass.TypesInfo.Types[s]; ok {
			checkTags(pass, s, tv.Type.(*types.Struct))
		}
	})
	return nil, nil
}

func checkTags(pass *analysis.Pass, node *ast.StructType, typ *types.Struct) {
	for fi := range typ.NumFields() {
		field := typ.Field(fi)
		// ignore unexported fields
		if !field.Exported() {
			continue
		}
		tag := typ.Tag(fi)
		jsonTag, found := reflect.StructTag(tag).Lookup("json")
		// ignore fields that
		//  1. Do not have a "json" tag
		//  2. Explicitly exclude the field from encoding
		//  3. Do not specify an "omitempty" option
		if !found || jsonTag == "-" || !hasOmitempty(jsonTag) {
			continue
		}
		checkField(pass, field, node.Fields.List[fi])
	}
}

func checkField(pass *analysis.Pass, field *types.Var, fieldNode *ast.Field) {
	// Ignore types with a zero value known to encoding/json to be omittable, or
	// those that are unencodable and will be reported elsewhere.
	switch typeDef := field.Type().Underlying().(type) {
	case *types.Basic,
		*types.Slice,
		*types.Pointer,
		*types.Map,
		*types.Chan,
		*types.Signature:
		return
	case *types.Interface:
		// TODO: while interfaces can always be nil, the user could attempt to
		// store a value in an interface that cannot be omitted, so we should
		// track down assignments of unomittable types and report these as well.
		// Once that's in place, we might as well revisit unencodable types and
		// report these as an error as well, regardless of whether they have
		// "omitempty" defined, saving the user the trouble of a runtime error.
		return
	case *types.Array:
		// Array types can only be omitted if their length is 0.
		if typeDef.Len() == 0 {
			return
		}
	default:
	}

	message := fmt.Sprintf(`field %q is marked "omitempty", but cannot be omitted`, field.Name())

	pass.Report(analysis.Diagnostic{
		Pos:     field.Pos(),
		End:     field.Pos() + token.Pos(len(field.Name())),
		Message: message,
		SuggestedFixes: []analysis.SuggestedFix{{
			Message: "Make field a pointer type",
			// TODO: fix assignments
			TextEdits: []analysis.TextEdit{{
				Pos:     fieldNode.Type.Pos(),
				End:     fieldNode.Type.Pos(),
				NewText: []byte("*"),
			}},
		}},
	})
}

func hasOmitempty(tagStr string) bool {
	_, options, _ := strings.Cut(tagStr, ",")
	return slices.Contains(strings.Split(options, ","), "omitempty")
}
