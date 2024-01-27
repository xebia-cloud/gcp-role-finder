package internal

import (
	"context"
	"encoding/base64"
	"fmt"
	"slices"
	"sort"
	"strings"
)

type Role struct {
	ID                  string   `json:"id"`
	Name                string   `json:"name"`
	Title               string   `json:"title"`
	Description         string   `json:"description"`
	IncludedPermissions []string `json:"includedPermissions"`
	Stage               string   `json:"stage"`
	ETag                string   `json:"etag"`
	PermissionCount     int      `json:"permissionCount"`
}

func NewRole(
	name string,
	title string,
	description string,
	includedPermissions []string,
	stage string,
	etag string,
) *Role {
	return &Role{
		ID:                  base64.StdEncoding.EncodeToString([]byte(name)),
		Name:                name,
		Title:               title,
		Description:         description,
		IncludedPermissions: includedPermissions,
		Stage:               stage,
		ETag:                etag,
		PermissionCount:     len(includedPermissions),
	}
}

type Roles []*Role

func (r Roles) Len() int      { return len(r) }
func (r Roles) Swap(i, j int) { r[i], r[j] = r[j], r[i] }

type (
	ByID              struct{ Roles }
	ByTitle           struct{ Roles }
	ByDescription     struct{ Roles }
	ByPermissionCount struct{ Roles }
	ByStage           struct{ Roles }
	ByETag            struct{ Roles }
)

func (a ByID) Less(i, j int) bool          { return a.Roles[i].ID < a.Roles[j].ID }
func (a ByTitle) Less(i, j int) bool       { return a.Roles[i].Title < a.Roles[j].Title }
func (a ByDescription) Less(i, j int) bool { return a.Roles[i].Description < a.Roles[j].Description }
func (a ByPermissionCount) Less(i, j int) bool {
	return a.Roles[i].PermissionCount < a.Roles[j].PermissionCount
}
func (a ByStage) Less(i, j int) bool { return a.Roles[i].Stage < a.Roles[j].Stage }
func (a ByETag) Less(i, j int) bool  { return a.Roles[i].ETag < a.Roles[j].ETag }

type RoleRepository interface {
	GetRoles(context.Context) (Roles, error)
	SaveRoles(context.Context, Roles) error
}

type RoleSearcher interface {
	FindRoles(context.Context, string) (Roles, error)
	IndexRoles(context.Context, Roles) error
	GetByID(context.Context, string) (*Role, error)
}

func (r Roles) Sort(sortBy string, sortOrder string) {
	switch sortBy {
	case "permissionCount":
		sort.Sort(ByPermissionCount{r})
	case "description":
		sort.Sort(ByDescription{r})
	case "title":
		sort.Sort(ByTitle{r})
	case "etag":
		sort.Sort(ByETag{r})
	case "stage":
		sort.Sort(ByStage{r})
	default:
		sort.Sort(ByID{r})
	}
	if sortOrder == "DESC" {
		slices.Reverse(r)
	}
}

func ValidateRoleParents(parents []string) error {
	for _, parent := range parents {
		if !(parent == "" ||
			strings.HasPrefix(parent, "projects/") ||
			strings.HasPrefix(parent, "organizations/")) {

			return fmt.Errorf("invalid parent: %s", parent)
		}
	}
	return nil
}
