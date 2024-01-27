package fulltext

import (
	"context"
	"encoding/base64"
	"testing"
	"xebia-gcloud/gcp-role-finder/internal"
)

func Test_repository_GetRoles(t *testing.T) {
	ctx := context.Background()

	roles := []*internal.Role{
		{
			base64.StdEncoding.EncodeToString([]byte("accessapproval.approver")),
			"accessapproval.approver",
			"Access Approval Approver",
			"Ability to view or act on access approval requests and view configuration",
			[]string{
				"accessapproval.requests.approve",
				"accessapproval.requests.dismiss",
				"accessapproval.requests.get",
				"accessapproval.requests.invalidate",
				"accessapproval.requests.list",
				"accessapproval.serviceAccounts.get",
				"accessapproval.settings.get",
				"resourcemanager.projects.get",
				"resourcemanager.projects.list",
			},
			"GA",
			"AA==",
			9,
		},
	}

	searcher, err := NewSearcher(ctx)
	if err != nil {
		t.Error(err)
	}
	if err = searcher.IndexRoles(ctx, roles); err != nil {
		t.Error(err)
	}
	matchingRoles, err := searcher.FindRoles(ctx, "accessapproval.requests.get")
	if err != nil {
		t.Error(err)
	}
	if len(matchingRoles) != 1 {
		t.Error("expected a single result")
	}

	searcher, err = NewSearcher(ctx, WithExcludedRoles([]string{"accessapproval.approver"}))
	if err != nil {
		t.Error(err)
	}
	if err = searcher.IndexRoles(ctx, roles); err != nil {
		t.Error(err)
	}
	matchingRoles, err = searcher.FindRoles(ctx, "accessapproval.requests.get")
	if err != nil {
		t.Error(err)
	}
	if len(matchingRoles) != 0 {
		t.Error("expected no result")
	}
}
