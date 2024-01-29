package gcp

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"xebia-cloud/gcp-role-finder/internal"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iam/v1"
)

type repository struct {
	service     *iam.Service
	credentials *google.Credentials
	parents     []string
}

func NewRepository(ctx context.Context, options ...func(r *repository)) (internal.RoleRepository, error) {
	r := &repository{}
	for _, o := range options {
		o(r)
	}

	var err error
	if r.credentials == nil {
		r.credentials, err = google.FindDefaultCredentials(ctx)
		if err != nil {
			return nil, err
		}
	}

	if len(r.parents) == 0 {
		// set sensible defaults here
		r.parents = append(r.parents, "")
		if r.credentials.ProjectID != "" {
			r.parents = append(r.parents, fmt.Sprintf("projects/%s", r.credentials.ProjectID))
		}
	}

	client := oauth2.NewClient(ctx, r.credentials.TokenSource)
	r.service, err = iam.New(client)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func WithCredentials(credentials *google.Credentials) func(*repository) {
	return func(r *repository) {
		r.credentials = credentials
	}
}

func WithParents(parents []string) func(*repository) {
	return func(r *repository) {
		r.parents = parents
	}
}

func (r *repository) GetRoles(ctx context.Context) (internal.Roles, error) {
	roles := make([]*internal.Role, 0, 2000)

	for _, parent := range r.parents {
		slog.InfoContext(ctx, "loading roles from parent", "parent", parent)
		err := r.service.Roles.List().View("FULL").Parent(parent).Context(ctx).Pages(ctx, func(r *iam.ListRolesResponse) error {
			for _, role := range r.Roles {
				roles = append(roles, internal.NewRole(
					role.Name,
					role.Title,
					role.Description,
					role.IncludedPermissions,
					role.Etag,
					role.Stage,
				))
			}
			return nil
		})
		if err != nil {
			slog.ErrorContext(ctx, "failed to retrieve roles from IAM", "error", err.Error())
			return nil, err
		}
	}

	return roles, nil
}

func (r *repository) SaveRoles(ctx context.Context, roles internal.Roles) error {
	return errors.New("I cannot write roles into IAM")
}
