package webapi

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gpuctl/gpuctl/internal/broadcast"
	"github.com/gpuctl/gpuctl/internal/onboard"
)

func (a *api) onboard(data broadcast.OnboardReq, _ *http.Request, log *slog.Logger) error {

	hostname := data.Hostname

	conf := a.onboardConf

	if hostname == "" {
		// TODO: return 400 bad request
		return errors.New("hostname cannot be empty")
	}

	// We error here, instead of on startup, so the rest of the API
	// methods will still work.
	if conf.DataDir == "" {
		return errors.New("`Onboard.datadir` must be set")
	}
	if conf.Signer == nil {
		return errors.New("no ssh key given to server")
	}
	if conf.Username == "" {
		return errors.New("`Onboard.username` must be set")
	}

	return onboard.Onboard(hostname,
		conf.Username,
		conf.DataDir,
		conf.Signer,
		conf.KeyCallback,
		conf.RemoteConf,
	)
}

func (a *api) deboard(data broadcast.RemoveMachineInfo,
	_ *http.Request,
	log *slog.Logger) error {
	hostname := data.Hostname

	conf := a.onboardConf

	if hostname == "" {
		// TODO: return 400 bad request
		return errors.New("hostname cannot be empty")
	}

	// We error here, instead of on startup, so the rest of the API
	// methods will still work.
	if conf.Signer == nil {
		return errors.New("no ssh key given to server")
	}
	if conf.Username == "" {
		return errors.New("`Onboard.username` must be set")
	}

	return onboard.Deboard(hostname, conf.Username, conf.Signer, conf.KeyCallback)
}
