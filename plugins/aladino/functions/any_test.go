// Copyright (C) 2022 zola - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions_test

import (
	"strings"
	"testing"

	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	plugins_aladino "github.com/reviewpad/reviewpad/v3/plugins/aladino"
	"github.com/stretchr/testify/assert"
)

var any = plugins_aladino.PluginBuiltIns().Functions["any"].Code

func TestAny(t *testing.T) {
	mockedEnv := aladino.MockDefaultEnv(t, nil, nil, aladino.MockBuiltIns(), nil)

	testCases := []struct {
		name    string
		args    []aladino.Value
		res     aladino.Value
		wantErr error
	}{
		{
			name: "matches one",
			args: []aladino.Value{
				aladino.BuildArrayValue(
					[]aladino.Value{
						aladino.BuildStringValue("a"),
						aladino.BuildStringValue("e"),
						aladino.BuildStringValue("f"),
					},
				),
				aladino.BuildFunctionValue(func(args []aladino.Value) aladino.Value {
					val := args[0].(*aladino.StringValue).Val
					return aladino.BuildBoolValue(strings.Contains("abcd", val))
				}),
			},
			res:     aladino.BuildBoolValue(true),
			wantErr: nil,
		},
		{
			name: "matches two",
			args: []aladino.Value{
				aladino.BuildArrayValue(
					[]aladino.Value{
						aladino.BuildStringValue("a"),
						aladino.BuildStringValue("b"),
						aladino.BuildStringValue("f"),
					},
				),
				aladino.BuildFunctionValue(func(args []aladino.Value) aladino.Value {
					val := args[0].(*aladino.StringValue).Val
					return aladino.BuildBoolValue(strings.Contains("abcd", val))
				}),
			},
			res:     aladino.BuildBoolValue(true),
			wantErr: nil,
		},
		{
			name: "matches all",
			args: []aladino.Value{
				aladino.BuildArrayValue(
					[]aladino.Value{
						aladino.BuildStringValue("a"),
						aladino.BuildStringValue("b"),
						aladino.BuildStringValue("c"),
					},
				),
				aladino.BuildFunctionValue(func(args []aladino.Value) aladino.Value {
					val := args[0].(*aladino.StringValue).Val
					return aladino.BuildBoolValue(strings.Contains("abcd", val))
				}),
			},
			res:     aladino.BuildBoolValue(true),
			wantErr: nil,
		},
		{
			name: "matches none",
			args: []aladino.Value{
				aladino.BuildArrayValue(
					[]aladino.Value{
						aladino.BuildStringValue("e"),
						aladino.BuildStringValue("f"),
						aladino.BuildStringValue("g"),
						aladino.BuildStringValue("h"),
						aladino.BuildStringValue("i"),
					},
				),
				aladino.BuildFunctionValue(func(args []aladino.Value) aladino.Value {
					val := args[0].(*aladino.StringValue).Val
					return aladino.BuildBoolValue(strings.Contains("abcd", val))
				}),
			},
			res:     aladino.BuildBoolValue(false),
			wantErr: nil,
		},
		{
			name: "empty list",
			args: []aladino.Value{
				aladino.BuildArrayValue(
					[]aladino.Value{},
				),
				aladino.BuildFunctionValue(func(args []aladino.Value) aladino.Value {
					val := args[0].(*aladino.StringValue).Val
					return aladino.BuildBoolValue(strings.Contains("abcd", val))
				}),
			},
			res:     aladino.BuildBoolValue(false),
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := any(mockedEnv, tc.args)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.res, res)
		})
	}
}
