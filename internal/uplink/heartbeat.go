// package uplink defines the types used in groundstation <-> satiate transport.

package uplink

const HeartbeatUrl = "/gs-api/heartbeat/"

type HeartbeatReq struct {
	Hostname string
}
