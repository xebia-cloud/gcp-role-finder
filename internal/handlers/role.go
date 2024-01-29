package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"xebia-cloud/gcp-role-finder/internal"
	"xebia-cloud/gcp-role-finder/internal/search/fulltext"

	"github.com/gofiber/fiber/v2"
)

type RoleHandler struct {
	roles      internal.Roles
	searcher   internal.RoleSearcher
	repository internal.RoleRepository
	mutex      sync.RWMutex
}

func NewRoleHandler(ctx context.Context, repository internal.RoleRepository) (*RoleHandler, error) {
	searcher, err := fulltext.NewSearcher(ctx)
	if err != nil {
		return nil, err
	}
	s := &RoleHandler{
		roles:      make(internal.Roles, 0),
		searcher:   searcher,
		repository: repository,
	}

	return s, nil
}

func (s *RoleHandler) RefreshRoles(ctx context.Context) error {
	slog.InfoContext(ctx, "refresh roles")
	roles, err := s.repository.GetRoles(ctx)
	if err != nil {
		return err
	}

	slog.InfoContext(ctx, "roles loaded", "count", len(roles))

	s.mutex.Lock()
	s.searcher.IndexRoles(ctx, roles)
	s.roles = roles
	s.mutex.Unlock()

	return nil
}

func (s *RoleHandler) List(c *fiber.Ctx) error {
	var err error
	var result internal.Roles

	start, end := parseRange(c.Query("range", "[0-9]"))
	sortBy, sortOrder := parseSort(c.Query("sort", "[\"id\",\"ASC\"]"))
	filters := parseFilters(c.Query("filter", "{}"))

	if q, ok := filters["q"]; ok && q != "" {
		s.mutex.Lock()
		searcher := s.searcher
		s.mutex.Unlock()
		result, err = searcher.FindRoles(c.Context(), q)
		if err != nil {
			return err
		}
	} else {
		s.mutex.Lock()
		result = make(internal.Roles, len(s.roles), len(s.roles))
		copy(result, s.roles)
		s.mutex.Unlock()
	}

	if start > len(result) {
		start = len(result)
	}
	if end > len(result) {
		end = len(result)
	}
	if start > end {
		start = end
	}
	result.Sort(sortBy, sortOrder)

	c.Append("Content-Range", fmt.Sprintf("roles %d-%d/%d", start, end, len(result)))
	if end+1 < len(result) {
		end = end + 1
	}
	return c.Status(fiber.StatusOK).JSON(result[start:end])
}

func (s *RoleHandler) GetRoleByID(c *fiber.Ctx) error {
	id := c.Params("id")
	s.mutex.Lock()
	searcher := s.searcher
	s.mutex.Unlock()

	if role, err := searcher.GetByID(c.Context(), id); err == nil {
		c.Status(fiber.StatusOK).JSON(role)
	} else {
		c.Status(fiber.StatusNotFound).JSON(map[string]string{"error": err.Error()})
	}
	return nil
}

func parseRange(input string) (start int, end int) {
	startEndRange := make([]int, 0, 2)

	err := json.Unmarshal([]byte(input), &startEndRange)
	if err != nil || len(startEndRange) != 2 {
		slog.Error("Invalid range input", "range", input)
	}

	if len(startEndRange) > 0 {
		start = startEndRange[0]
		end = startEndRange[len(startEndRange)-1]
	}
	if start < 0 {
		start = 0
	}
	if end < 0 {
		end = 0
	}

	if start > end {
		start = end
	}

	return
}

func parseSort(input string) (field string, order string) {
	sort := make([]string, 0, 2)

	err := json.Unmarshal([]byte(input), &sort)
	if err != nil || len(sort) != 2 {
		slog.Error("Invalid sort input", "sort", input)
	}

	if len(sort) > 0 {
		field = sort[0]
	}
	if len(sort) > 1 {
		order = strings.ToUpper(sort[len(sort)-1])
	}

	return
}

func parseFilters(input string) (filter map[string]string) {
	filter = make(map[string]string)
	err := json.Unmarshal([]byte(input), &filter)
	if err != nil {
		slog.Error("Invalid filter input", "filter", input)
	}
	return
}
