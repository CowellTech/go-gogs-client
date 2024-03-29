// Copyright 2016 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package gogs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Branch represents a repository branch.
type Branch struct {
	Name   string         `json:"name"`
	Commit *PayloadCommit `json:"commit"`
}

func (c *Client) ListRepoBranches(user, repo string) ([]*Branch, error) {
	branches := make([]*Branch, 0, 10)
	return branches, c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/branches", user, repo), nil, nil, &branches)
}

func (c *Client) GetRepoBranch(user, repo, branch string) (*Branch, error) {
	b := new(Branch)
	return b, c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/branches/%s", user, repo, branch), nil, nil, &b)
}

type CreateBranchOption struct {
	BranchName string `json:"branchname" binding:"Required"`
	Base       string `json:"base" binding:"Required"`
}

type DiffBranchInfo struct {
	ChangeInfo      string
	Owner           string
	Repo            string
	Branch1         string
	Branch2         string
	Branch1CommitId string
	Branch2CommitId string
	FileList        []DiffBranchChangeList
	Error           string
}

type DiffBranchChangeList struct {
	File     string
	IsBinary bool
}

type BranchList struct {
	BranchList []ProjectBranch
}

type ProjectBranch struct {
	Owner   string `json:"owner"`
	Repo    string `json:"repo"`
	Branch1 string `json:"branch1"`
	Branch2 string `json:"branch2"`
}

type DiffFileList struct {
	File                   string `json:"file"`
	IsBinary               bool   `json:"isBinary"`
	Project                string `json:"project"`
	ProjectOwner           string `json:"projectOwner"`
	BaseDiffBranchCommitID string `json:"baseDiffBranchCommitId"`
	DeployBranchCommitID   string `json:"deployBranchCommitId"`
}

type ReturnDiffFile struct {
	BaseInfo       DiffFileList
	BaseDiffFile   string
	BranchDiffFile string
	ErrorInfo      string
}

func (c *Client) GetBranchDiff(user, repo, branch1 string, branch2 string) (*DiffBranchInfo, error) {
	b := new(DiffBranchInfo)
	return b, c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/branch/diff/%s/%s", user, repo, branch1, branch2), nil, nil, &b)
}

func (c *Client) GetBranchsDiff(projectList []ProjectBranch) (*[]DiffBranchInfo, error) {
	body, err := json.Marshal(&projectList)
	if err != nil {
		return nil, err
	}
	b := new([]DiffBranchInfo)
	header := http.Header{
		"Content-Type": []string{"application/json"},
	}
	return b, c.getParsedResponse("POST", "/repos/branchs/diff", header, bytes.NewReader(body), &b)
}

func (c *Client) GetBranchsDiffFile(fileList []DiffFileList) (*[]ReturnDiffFile, error) {
	body, err := json.Marshal(&fileList)
	if err != nil {
		return nil, err
	}
	b := new([]ReturnDiffFile)
	header := http.Header{
		"Content-Type": []string{"application/json"},
	}
	return b, c.getParsedResponse("POST", "/repos/raw", header, bytes.NewReader(body), &b)
}

func (c *Client) CreateBranch(user, repo string, opt CreateBranchOption) (*Branch, error) {
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}

	b := new(Branch)
	header := http.Header{
		"Content-Type": []string{"application/json"},
	}
	return b, c.getParsedResponse("POST", fmt.Sprintf("/repos/%s/%s/branch", user, repo), header, bytes.NewReader(body), &b)
}

func (c *Client) DeleteBranch(user, repo, branch string, force bool) error {
	var urlFormat string
	urlFormat = "/repos/%s/%s/branch/%s"
	if force {
		urlFormat = "/repos/%s/%s/branch/%s?force=true"
	}

	_, err := c.getResponse("DELETE", fmt.Sprintf(urlFormat, user, repo, branch), nil, nil)
	return err
}

type CommitUserResponse struct {
	Name  string    `json:"Name"`
	Email string    `json:"Email"`
	When  time.Time `json:"When"`
}
type CommitResponse struct {
	ID        string             `json:"id"`
	Author    CommitUserResponse `json:"author"`
	Committer CommitUserResponse `json:"committer"`
	Message   string             `json:"message"`
}

func (c *Client) GetCommitsOfBranch(user, repo, branch, pagesize string) (*[]CommitResponse, error) {
	b := new([]CommitResponse)
	return b, c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/branch/commits/%s/%s", user, repo, pagesize, branch), nil, nil, &b)
}

type TestPatchRequest struct {
	BaseBranch   string `json:"baseBranch"`
	HeadBranch   string `json:"headBranch"`
	PatchContent string `json:"patchContent"`
}

type TestPatchResponse struct {
	BaseCommit  string `json:"baseCommit"`
	HeadCommit  string `json:"headCommit"`
	HasConflict bool   `json:"hasConflict"`
}

type AutoMergeRequest struct {
	BaseBranch        string `json:"baseBranch"`
	HeadBranch        string `json:"headBranch"`
	MergeStyle        string `json:"mergeStyle"`
	CommitDescription string `json:"commitDescription"`
}

type AutoMergeResponse struct {
	MergedCommit string `json:"mergedCommit"`
	Result       string `json:"result"` // MergeSuccess|MergeFailed|PushSuccess|PushFailed
	Msg          string `json:"msg"`
}

type DiffBinaryInfo struct {
	Head         string `json:"head"`
	BaseCommit   string `json:"baseCommit"`
	HeadCommit   string `json:"headCommit"`
	PatchContent []byte `json:"patchContent"`
}

func (c *Client) UpdatePatch(user, repo, base, head string) (*DiffBinaryInfo, error) {
	b := new(DiffBinaryInfo)
	return b, c.getParsedResponse("GET", fmt.Sprintf("/repos/%s/%s/branch/patch/%s/%s", user, repo, base, head), nil, nil, &b)
}

func (c *Client) TestPatch(user, repo string, opt TestPatchRequest) (*TestPatchResponse, error) {
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}

	b := new(TestPatchResponse)
	header := http.Header{
		"Content-Type": []string{"application/json"},
	}
	return b, c.getParsedResponse("POST", fmt.Sprintf("/repos/%s/%s/branches/patch/test", user, repo), header, bytes.NewReader(body), &b)
}

func (c *Client) AutoMerge(user, repo string, opt AutoMergeRequest) (*AutoMergeResponse, error) {
	body, err := json.Marshal(&opt)
	if err != nil {
		return nil, err
	}

	b := new(AutoMergeResponse)
	header := http.Header{
		"Content-Type": []string{"application/json"},
	}
	return b, c.getParsedResponse("POST", fmt.Sprintf("/repos/%s/%s/branches/automerge", user, repo), header, bytes.NewReader(body), &b)
}
