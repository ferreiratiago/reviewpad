// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"github.com/google/go-github/v48/github"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func PullRequestCountBy() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType(), aladino.BuildStringType()}, aladino.BuildIntType()),
		Code:           pullRequestCountByCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest, handler.Issue},
	}
}

func pullRequestCountByCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	loginArg := args[0].(*aladino.StringValue).Val
	stateArg := args[1].(*aladino.StringValue).Val

	state := "all"
	if stateArg != "" {
		state = stateArg
	}

	opts := &github.IssueListByRepoOptions{
		State: state,
	}
	if loginArg != "" {
		opts.Creator = loginArg
	}

	return issueCountBy(e, opts, func(issue *github.Issue) bool {
		return issue.IsPullRequest()
	})
}
