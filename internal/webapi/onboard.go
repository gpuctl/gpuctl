package webapi

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gpuctl/gpuctl/internal/broadcast"
	"github.com/gpuctl/gpuctl/internal/femto"
	"github.com/gpuctl/gpuctl/internal/tunnel"
	"github.com/gpuctl/gpuctl/internal/types"
)

func (a *Api) onboard(data broadcast.OnboardReq, _ *http.Request, log *slog.Logger) (*femto.EmptyBodyResponse, error) {

	hostname := data.Hostname
	conf := a.onboardConf

	if hostname == "" {
		// TODO: return 400 bad request
		return nil, errors.New("hostname cannot be empty")
	}

	err := tunnel.Onboard(hostname, conf)
	if err != nil {
		return nil, err
	}

	return femto.Ok(types.Unit{})
}

func (a *Api) deboard(data broadcast.RemoveMachineInfo,
	_ *http.Request,
	log *slog.Logger) error {
	hostname := data.Hostname

	conf := a.onboardConf

	if hostname == "" {
		// TODO: return 400 bad request
		return errors.New("hostname cannot be empty")
	}

	return tunnel.Deboard(hostname, conf)
}
