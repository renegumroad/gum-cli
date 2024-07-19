package updater

import (
	"net/http"

	"github.com/minio/selfupdate"
	"github.com/pkg/errors"
)

type UpdaterArgs struct {
	URL string
}

type Updater struct {
	args *UpdaterArgs
	http http.Client
}

func New(args *UpdaterArgs) *Updater {
	return newWithComponents(
		args,
		http.Client{},
	)
}

func newWithComponents(args *UpdaterArgs, client http.Client) *Updater {
	return &Updater{
		args: args,
		http: client,
	}
}

func (u *Updater) Update() error {
	resp, err := u.http.Get(u.args.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = selfupdate.Apply(resp.Body, selfupdate.Options{})
	if err != nil {
		if rerr := selfupdate.RollbackError(err); rerr != nil {
			return errors.Errorf("Failed to rollback from bad update. Reach out to DevOps for help: %v", rerr)
		}
	}
	return err
}
