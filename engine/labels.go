// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/google/go-github/v48/github"
)

func validateLabelColor(label *PadLabel) error {
	if label.Color != "" {
		matched, _ := regexp.MatchString(`(?i)^#?([0-9A-F]{6}){1,2}$`, label.Color)
		if !matched {
			return execError("evalLabel: color code not valid")
		}
	}

	return nil
}

func createLabel(e *Env, labelName *string, label *PadLabel) error {
	err := validateLabelColor(label)
	if err != nil {
		return err
	}

	var labelColor *string
	if label.Color != "" {
		label.Color = strings.TrimPrefix(label.Color, "#")
		labelColor = &label.Color
	}

	ghLabel := &github.Label{
		Name:        labelName,
		Color:       labelColor,
		Description: &label.Description,
	}

	owner := e.TargetEntity.Owner
	repo := e.TargetEntity.Repo

	_, _, err = e.GithubClient.CreateLabel(e.Ctx, owner, repo, ghLabel)

	return err
}

func checkLabelHasUpdates(e *Env, label *PadLabel, ghLabel *github.Label) (bool, error) {
	if ghLabel == nil {
		return false, fmt.Errorf("checkLabelHasUpdates: impossible to check updates on a empty github label")
	}

	if label.Description == "" {
		return false, nil
	}

	if ghLabel.Description != nil && *ghLabel.Description == label.Description {
		return false, nil
	}

	return true, nil
}

func updateLabel(e *Env, labelName *string, label *PadLabel) error {
	owner := e.TargetEntity.Owner
	repo := e.TargetEntity.Repo

	updatedGithubLabel := &github.Label{
		Description: &label.Description,
	}

	_, _, err := e.GithubClient.EditLabel(e.Ctx, owner, repo, *labelName, updatedGithubLabel)

	return err
}

func checkLabelExists(e *Env, labelName string) (*github.Label, error) {
	owner := e.TargetEntity.Owner
	repo := e.TargetEntity.Repo

	ghLabel, _, err := e.GithubClient.GetLabel(e.Ctx, owner, repo, labelName)
	if err != nil {
		if err.(*github.ErrorResponse).Response.StatusCode == 404 {
			return ghLabel, nil
		}

		return ghLabel, err
	}

	return ghLabel, nil
}
