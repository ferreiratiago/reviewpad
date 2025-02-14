// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransformAladinoExpression(t *testing.T) {
	tests := map[string]struct {
		arg     string
		wantVal string
	}{
		"id": {
			arg:     "$id()",
			wantVal: "$id()",
		},
		"merge": {
			arg:     "$merge()",
			wantVal: "$merge(\"merge\")",
		},
		"size": {
			arg:     "$size()",
			wantVal: "$size([])",
		},
		"issueCountBy simple": {
			arg:     "$issueCountBy(\"john\", \"open\") > 0",
			wantVal: "$issueCountBy(\"john\", \"open\") > 0",
		},
		"issueCountBy": {
			arg:     "$issueCountBy(\"john\", \"open\") > 0 && true && $issueCountBy(\"dev\") > 0",
			wantVal: "$issueCountBy(\"john\", \"open\") > 0 && true && $issueCountBy(\"dev\", \"all\") > 0",
		},
		"pullRequestCountBy simple": {
			arg:     "$pullRequestCountBy(\"john\") > 0",
			wantVal: "$pullRequestCountBy(\"john\", \"all\") > 0",
		},
		"pullRequestCountBy nil state": {
			arg:     "$pullRequestCountBy(\"john\", \"\") > 0",
			wantVal: "$pullRequestCountBy(\"john\", \"\") > 0",
		},
		"pullRequestCountBy nil dev": {
			arg:     "$pullRequestCountBy(\"\", \"closed\") > 0",
			wantVal: "$pullRequestCountBy(\"\", \"closed\") > 0",
		},
		"pullRequestCountBy nil dev and state": {
			arg:     "$pullRequestCountBy(\"\") > 0",
			wantVal: "$pullRequestCountBy(\"\", \"all\") > 0",
		},
		"pullRequestCountBy and issueCountBy": {
			arg:     "$pullRequestCountBy(\"john\") > 0 && true && $issueCountBy(\"dev\") > 0",
			wantVal: "$pullRequestCountBy(\"john\", \"all\") > 0 && true && $issueCountBy(\"dev\", \"all\") > 0",
		},
		"close": {
			arg:     "$close()",
			wantVal: `$close("", "completed")`,
		},
		"close with comment": {
			arg:     `$close("comment")`,
			wantVal: `$close("comment", "completed")`,
		},
		"close with comment and not_planned state reason": {
			arg:     `$close("comment", "not_planned")`,
			wantVal: `$close("comment", "not_planned")`,
		},
		"close with comment and completed state reason": {
			arg:     `$close("comment", "completed")`,
			wantVal: `$close("comment", "completed")`,
		},
		"close with empty comment and completed state reason": {
			arg:     `$close("", "completed")`,
			wantVal: `$close("", "completed")`,
		},
		"close with empty comment and not_planned state reason": {
			arg:     `$close("", "not_planned")`,
			wantVal: `$close("", "not_planned")`,
		},
		// TODO: test addDefaultTotalRequestedReviewers
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			gotVal := transformAladinoExpression(test.arg)
			assert.Equal(t, test.wantVal, gotVal)
		})
	}
}
