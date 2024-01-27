package fs

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"path"
	"xebia-gcloud/gcp-role-finder/internal"
)

type repository struct {
	fileName string
}

func NewRepository(ctx context.Context, options ...func(r *repository)) (internal.RoleRepository, error) {
	r := &repository{"data/roles.json"}
	for _, o := range options {
		o(r)
	}
	return r, nil
}

func WithFile(fileName string) func(*repository) {
	return func(r *repository) {
		r.fileName = fileName
	}
}

func (r *repository) GetRoles(ctx context.Context) (internal.Roles, error) {
	file, err := os.Open(r.fileName)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			slog.WarnContext(ctx, "failed to close file", "filename", r.fileName, "error", err.Error())
		}
	}()

	roles := make([]*internal.Role, 0, 2000)
	err = json.NewDecoder(file).Decode(&roles)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (r *repository) SaveRoles(ctx context.Context, roles internal.Roles) error {
	dir := path.Dir(r.fileName)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	file, err := os.OpenFile(r.fileName, os.O_RDWR|os.O_CREATE, 0o755)
	if err != nil {
		return err
	}

	defer func() {
		err := file.Close()
		if err != nil {
			slog.WarnContext(ctx, "failed to close file", "filename", r.fileName, "error", err.Error())
		}
	}()

	return json.NewEncoder(file).Encode(roles)
}
