package fulltext

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"xebia-gcloud/gcp-role-finder/internal"

	"github.com/blevesearch/bleve"
)

type searcher struct {
	index         bleve.Index
	roles         map[string]*internal.Role
	excludedRoles map[string]string
}

func NewSearcher(_ context.Context, options ...func(r *searcher)) (internal.RoleSearcher, error) {
	s := &searcher{
		excludedRoles: make(map[string]string, 16),
	}
	for _, o := range options {
		o(s)
	}
	return s, nil
}

func WithExcludedRoles(roles []string) func(*searcher) {
	return func(r *searcher) {
		for _, role := range roles {
			r.excludedRoles[base64.StdEncoding.EncodeToString([]byte(role))] = role
		}
	}
}

func (s *searcher) IndexRoles(_ context.Context, roles internal.Roles) error {
	rolesByID := make(map[string]*internal.Role, len(roles))
	for _, role := range roles {
		rolesByID[role.ID] = role
	}

	index, err := bleve.NewMemOnly(bleve.NewIndexMapping())
	if err != nil {
		return err
	}
	if roles != nil {
		for _, role := range roles {
			if err := index.Index(role.ID, role); err != nil {
				return err
			}
		}
	}
	s.index = index
	s.roles = rolesByID

	return nil
}

func (s *searcher) GetByID(ctx context.Context, id string) (*internal.Role, error) {
	if s.roles == nil {
		return nil, errors.New("no roles where indexed yet")
	}
	if role, ok := s.roles[id]; ok {
		return role, nil
	}
	return nil, fmt.Errorf("no role found with id %s", id)
}

func (s *searcher) FindRoles(ctx context.Context, q string) (internal.Roles, error) {
	if s.index == nil {
		return nil, errors.New("no roles where indexed yet")
	}

	query := bleve.NewQueryStringQuery(fmt.Sprintf("%s*", q))
	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Size = 100
	searchResult, err := s.index.SearchInContext(ctx, searchRequest)
	if err != nil {
		slog.ErrorContext(ctx, "search failed", "error", err.Error())
		return nil, err
	}

	result := make(internal.Roles, 0, len(searchResult.Hits))
	for _, hit := range searchResult.Hits {
		if _, ok := s.excludedRoles[hit.ID]; !ok {
			result = append(result, s.roles[hit.ID])
		}
	}

	sort.SliceStable(result, func(i, j int) bool {
		return len(result[i].IncludedPermissions) < len(result[j].IncludedPermissions)
	})

	slog.InfoContext(ctx, "search result", "hits", len(searchResult.Hits), "q", q)

	return result, nil
}
