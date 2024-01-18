// package uplink defines the types used in groundstation <-> satiate transport.

package uplink

import "time"

const HeartbeatUrl = "/gs-api/heartbeat/"

type HeartbeatReq struct {
	Time time.Time
}
