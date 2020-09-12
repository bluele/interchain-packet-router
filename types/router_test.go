package types

import (
	"errors"
	"fmt"
	"testing"

	sptypes "github.com/bluele/interchain-simple-packet/types"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/cosmos/cosmos-sdk/x/ibc/04-channel/exported"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

func TestRouter(t *testing.T) {
	require := require.New(t)

	r := New()
	r.HandleFunc(
		"/srv0",
		func(ctx sdk.Context, p PacketI, pd PacketDataI, sender PacketSender) (*sdk.Result, []byte, error) {
			return &sdk.Result{}, nil, nil
		},
		nil,
	)
	r.HandleFunc(
		"/srv1",
		func(ctx sdk.Context, p PacketI, pd PacketDataI, sender PacketSender) (*sdk.Result, []byte, error) {
			return nil, nil, errors.New("some error")
		},
		nil,
	)
	r.HandleFunc(
		"/srv3",
		func(ctx sdk.Context, p PacketI, pd PacketDataI, sender PacketSender) (*sdk.Result, []byte, error) {
			return &sdk.Result{}, []byte("ok"), nil
		},
		func(ctx sdk.Context, p PacketI, pd PacketDataI, ack []byte, sender PacketSender) (*sdk.Result, error) {
			return &sdk.Result{}, nil
		},
	)

	ctx := makeTestContext()

	h0 := sptypes.Header{}
	SetServiceID(&h0, "/srv0")
	_, _, err := r.ServePacket(ctx, channeltypes.Packet{}, sptypes.PacketData{Header: h0}, &mockPacketSender{})
	require.NoError(err)

	h1 := sptypes.Header{}
	SetServiceID(&h1, "/srv1")
	_, _, err = r.ServePacket(ctx, channeltypes.Packet{}, sptypes.PacketData{Header: h1}, &mockPacketSender{})
	require.Error(err)

	// not found handler
	h2 := sptypes.Header{}
	SetServiceID(&h2, "/srv2")
	_, _, err = r.ServePacket(ctx, channeltypes.Packet{}, sptypes.PacketData{Header: h2}, &mockPacketSender{})
	require.Error(err)

	h3 := sptypes.Header{}
	SetServiceID(&h3, "/srv3")
	_, ack, err := r.ServePacket(ctx, channeltypes.Packet{}, sptypes.PacketData{Header: h3}, &mockPacketSender{})
	require.NoError(err)
	_, err = r.ServeACK(ctx, channeltypes.Packet{}, sptypes.PacketData{Header: h3}, ack, &mockPacketSender{})
	require.NoError(err)
}

func TestMiddleware(t *testing.T) {
	require := require.New(t)

	var authMiddleware = func(next PacketHandlerFunc) PacketHandlerFunc {
		return func(ctx sdk.Context, p PacketI, pd PacketDataI, sender PacketSender) (*sdk.Result, []byte, error) {
			if p.GetSourceChannel() != "root" {
				return nil, nil, fmt.Errorf("unexpected channel id '%v'", p.GetSourceChannel())
			}
			sender = NewSendingPacketHandler(
				sender,
				func(ctx sdk.Context, channelCap *capabilitytypes.Capability, packet exported.PacketI) (exported.PacketI, error) {
					return packet, nil
				},
			)
			return next(ctx, p, pd, sender)
		}
	}

	r := New()
	r.UsePacketMiddlewares(authMiddleware)
	r.HandleFunc(
		"/srv0",
		func(ctx sdk.Context, p PacketI, pd PacketDataI, sender PacketSender) (*sdk.Result, []byte, error) {
			return &sdk.Result{}, nil, nil
		},
		nil,
	)
	r.HandleFunc(
		"/srv1/send",
		func(ctx sdk.Context, p PacketI, pd PacketDataI, sender PacketSender) (*sdk.Result, []byte, error) {
			if err := sender.SendPacket(ctx, nil, nil); err != nil {
				return nil, nil, err
			}
			return &sdk.Result{}, nil, nil
		},
		nil,
	)

	ctx := makeTestContext()
	h0 := sptypes.Header{}
	SetServiceID(&h0, "/srv0")
	_, _, err := r.ServePacket(ctx, channeltypes.Packet{SourceChannel: "root"}, sptypes.PacketData{Header: h0}, &mockPacketSender{})
	require.NoError(err)
	_, _, err = r.ServePacket(ctx, channeltypes.Packet{SourceChannel: "user0"}, sptypes.PacketData{Header: h0}, &mockPacketSender{})
	require.Error(err)

	h1 := sptypes.Header{}
	SetServiceID(&h1, "/srv1/send")
	_, _, err = r.ServePacket(ctx, channeltypes.Packet{SourceChannel: "root"}, sptypes.PacketData{Header: h1}, &mockPacketSender{})
	require.NoError(err)
}

type mockPacketSender struct{}

var _ PacketSender = (*mockPacketSender)(nil)

func (s *mockPacketSender) SendPacket(
	ctx sdk.Context,
	channelCap *capabilitytypes.Capability,
	packet exported.PacketI,
) error {
	return nil
}

func makeTestContext() sdk.Context {
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	return sdk.NewContext(cms, tmproto.Header{}, false, log.NewNopLogger())
}
