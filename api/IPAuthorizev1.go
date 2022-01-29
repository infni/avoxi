package api

import (
	"context"
	"time"

	"ipauthorize/internal/pkg/countycodes"
)

type IPAuthorizev1 struct {
	UnimplementedIpAuthorizeServer
	comparer countycodes.CountryCodeComparer
}

func NewIPAuthorizev1(comparer countycodes.CountryCodeComparer) *IPAuthorizev1 {
	svc := new(IPAuthorizev1)
	svc.comparer = comparer
	return svc
}

func (*IPAuthorizev1) Health(context.Context, *HealthRequest) (*HealthResponse, error) {
	return &HealthResponse{Now: time.Now().Format(time.RFC3339Nano)}, nil
}

func (svc *IPAuthorizev1) IpAuhtorize(_ context.Context, reqesut *IpAuthorizeRequest) (response *IpAuthorizeResponse, err error) {
	isInCountry, err := svc.comparer.IsInCountry(reqesut.Ip, reqesut.CountryNames)

	return &IpAuthorizeResponse{IsAuthorized: isInCountry}, err

}
