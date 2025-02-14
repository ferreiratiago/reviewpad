// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_actions

import (
	"fmt"

	"github.com/reviewpad/reviewpad/v3/codehost/github/target"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	"github.com/reviewpad/reviewpad/v3/utils"
)

func AssignRandomReviewer() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type:           aladino.BuildFunctionType([]aladino.Type{}, nil),
		Code:           assignRandomReviewerCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func assignRandomReviewerCode(e aladino.Env, _ []aladino.Value) error {
	t := e.GetTarget().(*target.PullRequestTarget)

	reviewers, err := t.GetReviewers()
	if err != nil {
		return err
	}

	// When there's already assigned reviewers, do nothing
	if len(reviewers.Users) > 0 {
		return nil
	}

	ghUsers, err := t.GetAvailableAssignees()
	if err != nil {
		return err
	}

	filteredGhUsers := []string{}

	user, err := t.GetAuthor()
	if err != nil {
		return err
	}

	for _, ghUser := range ghUsers {
		if ghUser.Login != user.Login {
			filteredGhUsers = append(filteredGhUsers, ghUser.Login)
		}
	}

	if len(filteredGhUsers) == 0 {
		return fmt.Errorf("can't assign a random user because there is no users")
	}

	lucky := utils.GenerateRandom(len(filteredGhUsers))
	ghUser := filteredGhUsers[lucky]

	return t.RequestReviewers([]string{ghUser})
}
