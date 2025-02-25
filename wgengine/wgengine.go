package wgengine

type Engine interface {
	// Reconfig reconfigures WireGuard and makes sure it's running.
	// This also handles setting up any kernel routes.
	//
	// This is called whenever tailcontrol (the control plane)
	// sends an updated network map.
	//
	// The returned error is ErrNoChanges if no changes were made.
	//Reconfig(*wgcfg.Config, *router.Config) error

	// GetFilter returns the current packet filter, if any.
	//GetFilter() *filter.Filter

	// SetFilter updates the packet filter.
	//SetFilter(*filter.Filter)

	// SetDNSMap updates the DNS map.
	//SetDNSMap(*tsdns.Map)

	// SetStatusCallback sets the function to call when the
	// WireGuard status changes.
	//SetStatusCallback(StatusCallback)

	// RequestStatus requests a WireGuard status update right
	// away, sent to the callback registered via SetStatusCallback.
	//RequestStatus()

	// Close shuts down this wireguard instance, remove any routes
	// it added, etc. To bring it up again later, you'll need a
	// new Engine.
	Close()

	// Wait waits until the Engine's Close method is called or the
	// engine aborts with an error.
	Wait() chan struct{}

	// LinkChange informs the engine that the system network
	// link has changed. The isExpensive parameter is set on links
	// where sending packets uses substantial power or money,
	// such as mobile data on a phone.
	//
	// LinkChange should be called whenever something changed with
	// the network, no matter how minor. The implementation should
	// look at the state of the network and decide whether the
	// change from before is interesting enough to warrant taking
	// action on.
	//LinkChange(isExpensive bool)

	// SetDERPMap controls which (if any) DERP servers are used.
	// If nil, DERP is disabled. It starts disabled until a DERP map
	// is configured.
	//SetDERPMap(*tailcfg.DERPMap)

	// SetNetworkMap informs the engine of the latest network map
	// from the server. The network map's DERPMap field should be
	// ignored as as it might be disabled; get it from SetDERPMap
	// instead.
	// The network map should only be read from.
	//SetNetworkMap(*controlclient.NetworkMap)

	// SetNetInfoCallback sets the function to call when a
	// new NetInfo summary is available.
	//SetNetInfoCallback(NetInfoCallback)

	// DiscoPublicKey gets the public key used for path discovery
	// messages.
	//DiscoPublicKey() tailcfg.DiscoKey

	// UpdateStatus populates the network state using the provided
	// status builder.
	//UpdateStatus(*ipnstate.StatusBuilder)
}
