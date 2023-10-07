package dynamic_ip_matcher

import (
	"encoding/json"
	"net"
	"net/http"
	"net/netip"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
)

func init() {
	caddy.RegisterModule(MatchDynamicRemoteIP{})
}

// MatchDynamicRemoteIP matchers the requests by the remote IP address.
// The IP ranges are provided by modules to allow for dynamic ranges.
type MatchDynamicRemoteIP struct {
	// A module which provides a source of IP ranges, from which
	// requests are matched.
	ProvidersRaw json.RawMessage         `json:"providers,omitempty" caddy:"namespace=http.ip_sources inline_key=source"`
	Providers    caddyhttp.IPRangeSource `json:"-"`

	logger *zap.Logger
}

func (MatchDynamicRemoteIP) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID: "http.matchers.dynamic_remote_ip",
		New: func() caddy.Module {
			return new(MatchDynamicRemoteIP)
		},
	}
}

func (m *MatchDynamicRemoteIP) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	d.Next() // consume the directive name

	if !d.NextArg() {
		return d.ArgErr()
	}

	if m.Providers != nil {
		return d.Err("providers already specified")
	}

	dynModule := d.Val()
	modID := "http.ip_sources." + dynModule
	mod, err := caddyfile.UnmarshalModule(d, modID)

	if err != nil {
		return err
	}

	provider, ok := mod.(caddyhttp.IPRangeSource)

	if !ok {
		return d.Errf("module %s (%T) is not an IPRangeSource", modID, mod)
	}

	m.ProvidersRaw = caddyconfig.JSONModuleObject(provider, "source", dynModule, nil)

	return nil
}

func (m *MatchDynamicRemoteIP) Provision(ctx caddy.Context) error {
	m.logger = ctx.Logger()

	if m.ProvidersRaw != nil {
		val, err := ctx.LoadModule(m, "ProvidersRaw")

		if err != nil {
			return err
		}

		m.Providers = val.(caddyhttp.IPRangeSource)
	}

	return nil
}

func (m MatchDynamicRemoteIP) Match(r *http.Request) bool {
	address := r.RemoteAddr
	remoteIP, err := parseIPZoneFromString(address)

	if err != nil {
		m.logger.Error("getting remote IP", zap.Error(err))
		return false
	}

	matches := m.matchIP(r, remoteIP)

	return matches
}

func parseIPZoneFromString(address string) (netip.Addr, error) {
	ipStr, _, err := net.SplitHostPort(address)
	if err != nil {
		ipStr = address // OK; probably didn't have a port
	}

	// Remote IP may contain a zone if IPv6, so we need
	// to pull that out before parsing the IP
	ipStr, _, _ = strings.Cut(ipStr, "%")

	ipAddr, err := netip.ParseAddr(ipStr)
	if err != nil {
		return netip.IPv4Unspecified(), err
	}

	return ipAddr, nil
}

func (m *MatchDynamicRemoteIP) matchIP(r *http.Request, remoteIP netip.Addr) bool {
	if m.Providers == nil {
		// We have no prover, So we can't match anything
		return false
	}

	cidrs := m.Providers.GetIPRanges(r)

	// TODO move this to an address map for performance
	// https://stackoverflow.com/questions/53397369/fastest-way-search-ip-in-large-ip-subnet-list-on-golang
	for _, ipRange := range cidrs {
		if ipRange.Contains(remoteIP) {
			return true
		}
	}
	return false
}

// Interface guards
var (
	_ caddy.Module             = (*MatchDynamicRemoteIP)(nil)
	_ caddy.Provisioner        = (*MatchDynamicRemoteIP)(nil)
	_ caddyfile.Unmarshaler    = (*MatchDynamicRemoteIP)(nil)
	_ caddyhttp.RequestMatcher = (*MatchDynamicRemoteIP)(nil)
)
