// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"fmt"

	"github.com/reviewpad/reviewpad/v3/codehost/github/target"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func AssignTeamReviewer() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildArrayOfType(aladino.BuildStringType())}, nil),
		Code:           assignTeamReviewerCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func assignTeamReviewerCode(e aladino.Env, args []aladino.Value) error {
	t := e.GetTarget().(*target.PullRequestTarget)

	teamReviewers := args[0].(*aladino.ArrayValue).Vals

	if len(teamReviewers) < 1 {
		return fmt.Errorf("assignTeamReviewer: requires at least 1 team to request for review")
	}

	teamReviewersSlugs := make([]string, len(teamReviewers))

	for i, team := range teamReviewers {
		teamReviewersSlugs[i] = team.(*aladino.StringValue).Val
	}

	return t.RequestTeamReviewers(teamReviewersSlugs)
}
