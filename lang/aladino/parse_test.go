// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package aladino

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := map[string]struct {
		input    string
		wantExpr Expr
	}{
		"single line": {
			input: `$addLabel("small")`,
			wantExpr: BuildFunctionCall(
				BuildVariable("addLabel"),
				[]Expr{BuildStringConst("small")},
			),
		},
		"single line with spaces": {
			input: `$addLabel(  "medium multiline")`,
			wantExpr: BuildFunctionCall(
				BuildVariable("addLabel"),
				[]Expr{BuildStringConst("medium multiline")},
			),
		},
		"comment with a string": {
			input: `["hello <a href='https://www.google.com'></a> world", "hello world"]`,
			wantExpr: BuildArray([]Expr{
				BuildStringConst("hello <a href='https://www.google.com'></a> world"),
				BuildStringConst("hello world"),
			}),
		},
		"comment a string with spaced strings": {
			input: `$comment("hello \"world\" and \"world\" again")`,
			wantExpr: BuildFunctionCall(
				BuildVariable("comment"),
				[]Expr{BuildStringConst("hello \"world\" and \"world\" again")},
			),
		},
		"array of strings": {
			input: `["hello", "world"]`,
			wantExpr: BuildArray([]Expr{
				BuildStringConst("hello"),
				BuildStringConst("world"),
			}),
		},
		"array of typed expression": {
			input: `[($developer: String => $developer), "hello", ($dev: String => $dev == "hello")]`,
			wantExpr: BuildArray([]Expr{
				BuildLambda(
					[]Expr{
						BuildTypedExpr(
							BuildVariable("developer"),
							BuildStringType(),
						),
					},
					BuildVariable("developer"),
				),
				BuildStringConst("hello"),
				BuildLambda(
					[]Expr{
						BuildTypedExpr(
							BuildVariable("dev"),
							BuildStringType(),
						),
					},
					BuildBinaryOp(BuildVariable("dev"), eqOperator(), BuildStringConst("hello")),
				),
			}),
		},
		"lambda no arguments": {
			input: `( => 10)`,
			wantExpr: BuildLambda(
				[]Expr{},
				BuildIntConst(10),
			),
		},
		"lambda single argument": {
			input: `($dev: String => $totalCreatedPullRequests($dev) == 10)`,
			wantExpr: BuildLambda(
				[]Expr{BuildTypedExpr(BuildVariable("dev"), BuildStringType())},
				BuildBinaryOp(BuildFunctionCall(BuildVariable("totalCreatedPullRequests"), []Expr{BuildVariable("dev")}), eqOperator(), BuildIntConst(10)),
			),
		},
		"lambda single argument and operation": {
			input: `($dev: String => $dev == "hello")`,
			wantExpr: BuildLambda(
				[]Expr{BuildTypedExpr(BuildVariable("dev"), BuildStringType())},
				BuildBinaryOp(BuildVariable("dev"), eqOperator(), BuildStringConst("hello")),
			),
		},
		"lambda multiple arguments": {
			input: `($a: Int, $b: Int => $a > $b)`,
			wantExpr: BuildLambda(
				[]Expr{
					BuildTypedExpr(BuildVariable("a"), BuildIntType()),
					BuildTypedExpr(BuildVariable("b"), BuildIntType()),
				},
				BuildBinaryOp(BuildVariable("a"), greaterThanOperator(), BuildVariable("b")),
			),
		},
		"higher order functions": {
			input: `$any($reviewers(), ($dev: String => $isElementOf($dev, $team("security"))))`,
			wantExpr: BuildFunctionCall(
				BuildVariable("any"),
				[]Expr{
					BuildFunctionCall(
						BuildVariable("reviewers"),
						[]Expr{},
					),
					BuildLambda(
						[]Expr{BuildTypedExpr(BuildVariable("dev"), BuildStringType())},
						BuildFunctionCall(
							BuildVariable("isElementOf"),
							[]Expr{
								BuildVariable("dev"),
								BuildFunctionCall(
									BuildVariable("team"),
									[]Expr{BuildStringConst("security")},
								),
							},
						),
					),
				},
			),
		},
		"nested lambda": {
			input: `($a: Int => ($b: Int => $a > $b))`,
			wantExpr: BuildLambda(
				[]Expr{BuildTypedExpr(BuildVariable("a"), BuildIntType())},
				BuildLambda(
					[]Expr{BuildTypedExpr(BuildVariable("b"), BuildIntType())},
					BuildBinaryOp(BuildVariable("a"), greaterThanOperator(), BuildVariable("b")),
				),
			),
		},
		"typed expression lambda": {
			input: `($a: Int, $b: Int => $a > $b)`,
			wantExpr: BuildLambda(
				[]Expr{BuildTypedExpr(BuildVariable("a"), BuildIntType()),
					BuildTypedExpr(BuildVariable("b"), BuildIntType())},
				BuildBinaryOp(BuildVariable("a"), greaterThanOperator(), BuildVariable("b")),
			),
		},
		"multi parameter lambda": {
			input: `($a: []Int, $b: []Bool, $c: []String, $d: Func(String) String, $e: Func(Int) String  => $concat($a, $b, $c))`,
			wantExpr: BuildLambda(
				[]Expr{
					BuildTypedExpr(BuildVariable("a"), BuildArrayOfType(BuildIntType())),
					BuildTypedExpr(BuildVariable("b"), BuildArrayOfType(BuildBoolType())),
					BuildTypedExpr(BuildVariable("c"), BuildArrayOfType(BuildStringType())),
					BuildTypedExpr(BuildVariable("d"), BuildFunctionType([]Type{BuildStringType()}, BuildStringType())),
					BuildTypedExpr(BuildVariable("e"), BuildFunctionType([]Type{BuildIntType()}, BuildStringType())),
				},
				BuildFunctionCall(
					BuildVariable("concat"),
					[]Expr{
						BuildVariable("a"),
						BuildVariable("b"),
						BuildVariable("c"),
					},
				),
			),
		},
		"string array typed expression lambda": {
			input: `($orgs: []String => $isElementOf("reviewpad", $orgs))`,
			wantExpr: BuildLambda(
				[]Expr{
					BuildTypedExpr(BuildVariable("orgs"), BuildArrayOfType(BuildStringType())),
				},
				BuildFunctionCall(
					BuildVariable("isElementOf"),
					[]Expr{
						BuildStringConst("reviewpad"),
						BuildVariable("orgs"),
					},
				),
			),
		},
		"int array typed expression lambda": {
			input: `($nums: []Int => $isElementOf(5, $nums))`,
			wantExpr: BuildLambda(
				[]Expr{
					BuildTypedExpr(BuildVariable("nums"), BuildArrayOfType(BuildIntType())),
				},
				BuildFunctionCall(
					BuildVariable("isElementOf"),
					[]Expr{
						BuildIntConst(5),
						BuildVariable("nums"),
					},
				),
			),
		},
		"boolean array typed expression lambda": {
			input: `($flags: []Bool => $isElementOf(true, $flags))`,
			wantExpr: BuildLambda(
				[]Expr{
					BuildTypedExpr(BuildVariable("flags"), BuildArrayOfType(BuildBoolType())),
				},
				BuildFunctionCall(
					BuildVariable("isElementOf"),
					[]Expr{
						BuildBoolConst(true),
						BuildVariable("flags"),
					},
				),
			),
		},
		"array of string array typed expression lambda": {
			input: `[($orgs: []String => $isElementOf("reviewpad", $orgs)), "hello", ($users: []String => $isElementOf("user", $users))]`,
			wantExpr: BuildArray(
				[]Expr{
					BuildLambda(
						[]Expr{
							BuildTypedExpr(BuildVariable("orgs"), BuildArrayOfType(BuildStringType())),
						},
						BuildFunctionCall(
							BuildVariable("isElementOf"),
							[]Expr{
								BuildStringConst("reviewpad"),
								BuildVariable("orgs"),
							},
						),
					),
					BuildStringConst("hello"),
					BuildLambda(
						[]Expr{
							BuildTypedExpr(BuildVariable("users"), BuildArrayOfType(BuildStringType())),
						},
						BuildFunctionCall(
							BuildVariable("isElementOf"),
							[]Expr{
								BuildStringConst("user"),
								BuildVariable("users"),
							},
						),
					),
				},
			),
		},
		"array of int array typed expression lambda": {
			input: `[($nums: []Int => $isElementOf(5, $nums)), "hi", 1, ($ids: []Int => $isElementOf(5, $ids))]`,
			wantExpr: BuildArray(
				[]Expr{
					BuildLambda(
						[]Expr{
							BuildTypedExpr(BuildVariable("nums"), BuildArrayOfType(BuildIntType())),
						},
						BuildFunctionCall(
							BuildVariable("isElementOf"),
							[]Expr{
								BuildIntConst(5),
								BuildVariable("nums"),
							},
						),
					),
					BuildStringConst("hi"),
					BuildIntConst(1),
					BuildLambda(
						[]Expr{
							BuildTypedExpr(BuildVariable("ids"), BuildArrayOfType(BuildIntType())),
						},
						BuildFunctionCall(
							BuildVariable("isElementOf"),
							[]Expr{
								BuildIntConst(5),
								BuildVariable("ids"),
							},
						),
					),
				},
			),
		},
		"array of boolean array typed expression lambda": {
			input: `[($flags: []Bool => $isElementOf(true, $flags)), "hi", 1, ($enabled: []Bool => $isElementOf(false, $enabled))]`,
			wantExpr: BuildArray(
				[]Expr{
					BuildLambda(
						[]Expr{
							BuildTypedExpr(BuildVariable("flags"), BuildArrayOfType(BuildBoolType())),
						},
						BuildFunctionCall(
							BuildVariable("isElementOf"),
							[]Expr{
								BuildBoolConst(true),
								BuildVariable("flags"),
							},
						),
					),
					BuildStringConst("hi"),
					BuildIntConst(1),
					BuildLambda(
						[]Expr{
							BuildTypedExpr(BuildVariable("enabled"), BuildArrayOfType(BuildBoolType())),
						},
						BuildFunctionCall(
							BuildVariable("isElementOf"),
							[]Expr{
								BuildBoolConst(false),
								BuildVariable("enabled"),
							},
						),
					),
				},
			),
		},
		"array of string lambda": {
			input: `($a: []String => $length($a))`,
			wantExpr: BuildLambda(
				[]Expr{
					BuildTypedExpr(BuildVariable("a"), BuildArrayOfType(BuildStringType())),
				},
				BuildFunctionCall(
					BuildVariable("length"),
					[]Expr{
						BuildVariable("a"),
					},
				),
			),
		},
		"array of int lambda": {
			input: `($a: []Int => $length($a))`,
			wantExpr: BuildLambda(
				[]Expr{
					BuildTypedExpr(BuildVariable("a"), BuildArrayOfType(BuildIntType())),
				},
				BuildFunctionCall(
					BuildVariable("length"),
					[]Expr{
						BuildVariable("a"),
					},
				),
			),
		},
		"array of bool lambda": {
			input: `($a: []Bool => $length($a))`,
			wantExpr: BuildLambda(
				[]Expr{
					BuildTypedExpr(BuildVariable("a"), BuildArrayOfType(BuildBoolType())),
				},
				BuildFunctionCall(
					BuildVariable("length"),
					[]Expr{
						BuildVariable("a"),
					},
				),
			),
		},
		"nested array of string lambda": {
			input: `($a: [][]String => $length($a))`,
			wantExpr: BuildLambda(
				[]Expr{
					BuildTypedExpr(BuildVariable("a"), BuildArrayOfType(BuildArrayOfType(BuildStringType()))),
				},
				BuildFunctionCall(
					BuildVariable("length"),
					[]Expr{
						BuildVariable("a"),
					},
				),
			),
		},
		"thriple nested array of string lambda": {
			input: `($a: [][][]String => $length($a))`,
			wantExpr: BuildLambda(
				[]Expr{
					BuildTypedExpr(BuildVariable("a"), BuildArrayOfType(BuildArrayOfType(BuildArrayOfType(BuildStringType())))),
				},
				BuildFunctionCall(
					BuildVariable("length"),
					[]Expr{
						BuildVariable("a"),
					},
				),
			),
		},
		"nested array of int lambda": {
			input: `($a: [][]Int => $length($a))`,
			wantExpr: BuildLambda(
				[]Expr{
					BuildTypedExpr(BuildVariable("a"), BuildArrayOfType(BuildArrayOfType(BuildIntType()))),
				},
				BuildFunctionCall(
					BuildVariable("length"),
					[]Expr{
						BuildVariable("a"),
					},
				),
			),
		},
		"thriple nested array of int lambda": {
			input: `($a: [][][]Int => $length($a))`,
			wantExpr: BuildLambda(
				[]Expr{
					BuildTypedExpr(BuildVariable("a"), BuildArrayOfType(BuildArrayOfType(BuildArrayOfType(BuildIntType())))),
				},
				BuildFunctionCall(
					BuildVariable("length"),
					[]Expr{
						BuildVariable("a"),
					},
				),
			),
		},
		"nested array of bool lambda": {
			input: `($a: [][]Bool => $length($a))`,
			wantExpr: BuildLambda(
				[]Expr{
					BuildTypedExpr(BuildVariable("a"), BuildArrayOfType(BuildArrayOfType(BuildBoolType()))),
				},
				BuildFunctionCall(
					BuildVariable("length"),
					[]Expr{
						BuildVariable("a"),
					},
				),
			),
		},
		"thriple nested array of bool lambda": {
			input: `($a: [][][]Bool => $length($a))`,
			wantExpr: BuildLambda(
				[]Expr{
					BuildTypedExpr(BuildVariable("a"), BuildArrayOfType(BuildArrayOfType(BuildArrayOfType(BuildBoolType())))),
				},
				BuildFunctionCall(
					BuildVariable("length"),
					[]Expr{
						BuildVariable("a"),
					},
				),
			),
		},
		"function typed expression lambda": {
			input: `($a: Func(Int, String) String, $b: Func(Int) Int => $a(1, "a") > $b(2))`,
			wantExpr: BuildLambda(
				[]Expr{
					BuildTypedExpr(BuildVariable("a"), BuildFunctionType([]Type{BuildIntType(), BuildStringType()}, BuildStringType())),
					BuildTypedExpr(BuildVariable("b"), BuildFunctionType([]Type{BuildIntType()}, BuildIntType())),
				},
				BuildBinaryOp(
					BuildFunctionCall(BuildVariable("a"), []Expr{BuildIntConst(1), BuildStringConst("a")}),
					greaterThanOperator(),
					BuildFunctionCall(BuildVariable("b"), []Expr{BuildIntConst(2)}),
				),
			),
		},
		"array of func(string) string lambda": {
			input: `($a: []Func(String) String => $length($a))`,
			wantExpr: BuildLambda(
				[]Expr{
					BuildTypedExpr(BuildVariable("a"), BuildArrayOfType(BuildFunctionType([]Type{BuildStringType()}, BuildStringType()))),
				},
				BuildFunctionCall(
					BuildVariable("length"),
					[]Expr{
						BuildVariable("a"),
					},
				),
			),
		},
		"array of func(string) int lambda": {
			input: `($a: []Func(String) Int => $length($a))`,
			wantExpr: BuildLambda(
				[]Expr{
					BuildTypedExpr(BuildVariable("a"), BuildArrayOfType(BuildFunctionType([]Type{BuildStringType()}, BuildIntType()))),
				},
				BuildFunctionCall(
					BuildVariable("length"),
					[]Expr{
						BuildVariable("a"),
					},
				),
			),
		},
		"array of func(string) bool lambda": {
			input: `($a: []Func(String) Bool => $length($a))`,
			wantExpr: BuildLambda(
				[]Expr{
					BuildTypedExpr(BuildVariable("a"), BuildArrayOfType(BuildFunctionType([]Type{BuildStringType()}, BuildBoolType()))),
				},
				BuildFunctionCall(
					BuildVariable("length"),
					[]Expr{
						BuildVariable("a"),
					},
				),
			),
		},
		"array of func(string) func(string) string lambda": {
			input: `($a: []Func(String) Func(String) String => $length($a))`,
			wantExpr: BuildLambda(
				[]Expr{
					BuildTypedExpr(
						BuildVariable("a"),
						BuildArrayOfType(
							BuildFunctionType(
								[]Type{BuildStringType()},
								BuildFunctionType(
									[]Type{BuildStringType()},
									BuildStringType(),
								),
							),
						),
					),
				},
				BuildFunctionCall(
					BuildVariable("length"),
					[]Expr{
						BuildVariable("a"),
					},
				),
			),
		},
		"array of func(string) func(string) int lambda": {
			input: `($a: []Func(String) Func(String) Int => $length($a))`,
			wantExpr: BuildLambda(
				[]Expr{
					BuildTypedExpr(
						BuildVariable("a"),
						BuildArrayOfType(
							BuildFunctionType(
								[]Type{BuildStringType()},
								BuildFunctionType(
									[]Type{BuildStringType()},
									BuildIntType(),
								),
							),
						),
					),
				},
				BuildFunctionCall(
					BuildVariable("length"),
					[]Expr{
						BuildVariable("a"),
					},
				),
			),
		},
		"array of func(string) func(string) bool lambda": {
			input: `($a: []Func(String) Func(String) Bool => $length($a))`,
			wantExpr: BuildLambda(
				[]Expr{
					BuildTypedExpr(
						BuildVariable("a"),
						BuildArrayOfType(
							BuildFunctionType(
								[]Type{BuildStringType()},
								BuildFunctionType(
									[]Type{BuildStringType()},
									BuildBoolType(),
								),
							),
						),
					),
				},
				BuildFunctionCall(
					BuildVariable("length"),
					[]Expr{
						BuildVariable("a"),
					},
				),
			),
		},
		"array of func(int) string lambda": {
			input: `($a: []Func(Int) String => $length($a))`,
			wantExpr: BuildLambda(
				[]Expr{
					BuildTypedExpr(BuildVariable("a"), BuildArrayOfType(BuildFunctionType([]Type{BuildIntType()}, BuildStringType()))),
				},
				BuildFunctionCall(
					BuildVariable("length"),
					[]Expr{
						BuildVariable("a"),
					},
				),
			),
		},
		"array of func(int) int lambda": {
			input: `($a: []Func(Int) Int => $length($a))`,
			wantExpr: BuildLambda(
				[]Expr{
					BuildTypedExpr(BuildVariable("a"), BuildArrayOfType(BuildFunctionType([]Type{BuildIntType()}, BuildIntType()))),
				},
				BuildFunctionCall(
					BuildVariable("length"),
					[]Expr{
						BuildVariable("a"),
					},
				),
			),
		},
		"array of func(int) bool lambda": {
			input: `($a: []Func(Int) Bool => $length($a))`,
			wantExpr: BuildLambda(
				[]Expr{
					BuildTypedExpr(BuildVariable("a"), BuildArrayOfType(BuildFunctionType([]Type{BuildIntType()}, BuildBoolType()))),
				},
				BuildFunctionCall(
					BuildVariable("length"),
					[]Expr{
						BuildVariable("a"),
					},
				),
			),
		},
		"array of func(bool) string lambda": {
			input: `($a: []Func(Bool) String => $length($a))`,
			wantExpr: BuildLambda(
				[]Expr{
					BuildTypedExpr(BuildVariable("a"), BuildArrayOfType(BuildFunctionType([]Type{BuildBoolType()}, BuildStringType()))),
				},
				BuildFunctionCall(
					BuildVariable("length"),
					[]Expr{
						BuildVariable("a"),
					},
				),
			),
		},
		"array of func(bool) int lambda": {
			input: `($a: []Func(Bool) Int => $length($a))`,
			wantExpr: BuildLambda(
				[]Expr{
					BuildTypedExpr(BuildVariable("a"), BuildArrayOfType(BuildFunctionType([]Type{BuildBoolType()}, BuildIntType()))),
				},
				BuildFunctionCall(
					BuildVariable("length"),
					[]Expr{
						BuildVariable("a"),
					},
				),
			),
		},
		"array of func(bool) bool lambda": {
			input: `($a: []Func(Bool) Bool => $length($a))`,
			wantExpr: BuildLambda(
				[]Expr{
					BuildTypedExpr(BuildVariable("a"), BuildArrayOfType(BuildFunctionType([]Type{BuildBoolType()}, BuildBoolType()))),
				},
				BuildFunctionCall(
					BuildVariable("length"),
					[]Expr{
						BuildVariable("a"),
					},
				),
			),
		},
		"nested array of func(string) string lambda": {
			input: `($a: [][]Func(String) String => $length($a))`,
			wantExpr: BuildLambda(
				[]Expr{
					BuildTypedExpr(BuildVariable("a"), BuildArrayOfType(BuildArrayOfType(BuildFunctionType([]Type{BuildStringType()}, BuildStringType())))),
				},
				BuildFunctionCall(
					BuildVariable("length"),
					[]Expr{
						BuildVariable("a"),
					},
				),
			),
		},
		"func(func(string) string) func(string) int lambda": {
			input: `($a: Func(Func(String) String) Func(String) Int => $length($a))`,
			wantExpr: BuildLambda(
				[]Expr{
					BuildTypedExpr(BuildVariable("a"), BuildFunctionType(
						[]Type{BuildFunctionType([]Type{BuildStringType()}, BuildStringType())},
						BuildFunctionType([]Type{BuildStringType()}, BuildIntType()),
					)),
				},
				BuildFunctionCall(
					BuildVariable("length"),
					[]Expr{
						BuildVariable("a"),
					},
				),
			),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			gotExpr, err := Parse(test.input)
			assert.Nil(t, err)
			assert.Equal(t, test.wantExpr, gotExpr)
		})
	}
}
