package gogs

import (
	"testing"
)

func TestCreateBranch(t *testing.T) {
	url := "http://10.2.96.250"
	token := "560e8c92bb4bc211bf028e0df95f2f4326408eaa"

	c := NewClient(url, token)

	opt := CreateBranchOption{
		Base:       "master",
		BranchName: "Feature-docker-k8s-gbren-test",
	}
	branch, err := c.CreateBranch("devops", "docker-k8s", opt)

	t.Log(err)
	t.Log(branch)

	if branch.Name != opt.BranchName {
		t.Error("failed")
	}
}
