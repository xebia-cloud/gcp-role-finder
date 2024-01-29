package fs

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"xebia-cloud/gcp-role-finder/internal"
)

func Test_repository_GetRoles(t *testing.T) {
	ctx := context.Background()

	roles := []*internal.Role{
		internal.NewRole("accessapproval.approver",
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
		),
	}

	filename := fmt.Sprintf("%s/data.json", t.TempDir())
	fileRepository, err := NewRepository(ctx, WithFile(filename))
	if err = fileRepository.SaveRoles(ctx, roles); err != nil {
		t.Error(err)
	}

	savedRoles, err := fileRepository.GetRoles(ctx)
	if err != nil {
		t.Error(err)
	}

	if len(roles) != len(savedRoles) {
		t.Errorf("mismatch in number of roles retrieved")
	}

	for i, v := range roles {
		if !reflect.DeepEqual(v, savedRoles[i]) {
			t.Errorf("mismatch in retrieved roles %v, %v", v, savedRoles[i])
		}
	}
}
