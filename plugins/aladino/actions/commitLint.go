// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"fmt"

	"github.com/reviewpad/go-conventionalcommits"
	"github.com/reviewpad/go-conventionalcommits/parser"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
)

func CommitLint() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:           aladino.BuildFunctionType([]aladino.Type{}, nil),
		Code:           commitLintCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func commitLintCode(e aladino.Env, _ []aladino.Value) error {
	entity := e.GetTarget().GetTargetEntity()

	prNum := entity.Number
	owner := entity.Owner
	repo := entity.Repo

	ghCommits, err := e.GetGithubClient().GetPullRequestCommits(e.GetCtx(), owner, repo, prNum)
	if err != nil {
		return err
	}

	for _, ghCommit := range ghCommits {
		commitMsg := ghCommit.Commit.GetMessage()
		res, err := parser.NewMachine(conventionalcommits.WithTypes(conventionalcommits.TypesConventional)).Parse([]byte(commitMsg))

		if err != nil || !res.Ok() {
			body := fmt.Sprintf("**Unconventional commit detected**: '%v' (%v)", commitMsg, ghCommit.GetSHA())
			reportedMessages := e.GetBuiltInsReportedMessages()
			reportedMessages[aladino.SEVERITY_ERROR] = append(reportedMessages[aladino.SEVERITY_ERROR], body)
		}
	}

	return nil
}
