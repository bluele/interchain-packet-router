package types

import (
	"github.com/bluele/interchain-simple-packet/types"
)

// PacketDataI defines the standard packet data.
type PacketDataI interface {
	GetHeader() HeaderI
	GetPayload() []byte
}

// HeaderI defines the standard header for a packet data.
type HeaderI = types.HeaderI

// PacketI defines the standard interface for IBC packets
type PacketI interface {
	GetSequence() uint64
	GetTimeoutHeight() uint64
	GetTimeoutTimestamp() uint64
	GetSourcePort() string
	GetSourceChannel() string
	GetDestPort() string
	GetDestChannel() string
	GetData() []byte
	ValidateBasic() error
}
