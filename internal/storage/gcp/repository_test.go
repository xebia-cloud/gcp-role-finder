package gcp

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"xebia-cloud/gcp-role-finder/internal/storage/fs"

	"golang.org/x/oauth2/google"
)

func Test_repository_GetRoles(t *testing.T) {
	ctx := context.Background()
	credentials, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		t.Error(err)
	}
	repository, err := NewRepository(ctx, WithCredentials(credentials))
	roles, err := repository.GetRoles(ctx)
	if err != nil {
		t.Error(err)
	}

	if roles.Len() < 1500 {
		t.Errorf("expected more than 1500 roles got %d", roles.Len())
	}

	filename := fmt.Sprintf("%s/data.json", t.TempDir())

	fileRepository, err := fs.NewRepository(ctx, fs.WithFile(filename))
	if err = fileRepository.SaveRoles(ctx, roles); err != nil {
		t.Error(err)
	}

	savedRoles, err := fileRepository.GetRoles(ctx)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(roles, savedRoles) {
		t.Errorf("mismatch in retrieved roles")
	}
}
