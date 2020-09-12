package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Handler interface {
	ServePacket(ctx sdk.Context, p PacketI, pd PacketDataI, sender PacketSender) (res *sdk.Result, ack []byte, err error)
	ServeACK(ctx sdk.Context, p PacketI, pd PacketDataI, ack []byte, sender PacketSender) (res *sdk.Result, err error)
}

type HandlerFunc struct {
	PacketHandlerFunc PacketHandlerFunc
	ACKHandlerFunc    ACKHandlerFunc
}

type (
	PacketHandlerFunc func(ctx sdk.Context, p PacketI, pd PacketDataI, sender PacketSender) (*sdk.Result, []byte, error)

	ACKHandlerFunc func(ctx sdk.Context, p PacketI, pd PacketDataI, ack []byte, sender PacketSender) (*sdk.Result, error)
)

var _ Handler = (*HandlerFunc)(nil)

func (f HandlerFunc) ServePacket(ctx sdk.Context, p PacketI, pd PacketDataI, sender PacketSender) (*sdk.Result, []byte, error) {
	return f.PacketHandlerFunc(ctx, p, pd, sender)
}

func (f HandlerFunc) ServeACK(ctx sdk.Context, p PacketI, pd PacketDataI, ack []byte, sender PacketSender) (*sdk.Result, error) {
	return f.ACKHandlerFunc(ctx, p, pd, ack, sender)
}
