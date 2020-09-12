package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/cosmos/cosmos-sdk/x/ibc/04-channel/exported"
	channelkeeper "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/keeper"
)

type PacketSender interface {
	SendPacket(
		ctx sdk.Context,
		channelCap *capabilitytypes.Capability,
		packet exported.PacketI,
	) error
}

var _ PacketSender = (*channelkeeper.Keeper)(nil)

type SendingPacketCallback func(
	ctx sdk.Context,
	channelCap *capabilitytypes.Capability,
	packet exported.PacketI,
) (exported.PacketI, error)

type SendingPacketHandler struct {
	sender   PacketSender
	callback SendingPacketCallback
}

var _ PacketSender = (*SendingPacketHandler)(nil)

func NewSendingPacketHandler(sender PacketSender, callback SendingPacketCallback) SendingPacketHandler {
	return SendingPacketHandler{
		sender:   sender,
		callback: callback,
	}
}

func (h SendingPacketHandler) SendPacket(
	ctx sdk.Context,
	channelCap *capabilitytypes.Capability,
	packet exported.PacketI,
) error {
	p, err := h.callback(ctx, channelCap, packet)
	if err != nil {
		return err
	}
	return h.sender.SendPacket(ctx, channelCap, p)
}
