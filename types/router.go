package types

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Router struct {
	// Routes to be matched, in order.
	routes []*Route
	// Slice of middlewares to be called after a match is found
	packetMiddlewares []packetMiddleware
	ackMiddlewares    []ackMiddleware

	// Configurable Handler to be used when no route matches.
	NotFoundHandler HandlerFunc
}

var _ Handler = (*Router)(nil)

func New() *Router {
	return &Router{}
}

// NewRoute registers an empty route.
func (r *Router) NewRoute() *Route {
	route := &Route{}
	r.routes = append(r.routes, route)
	return route
}

func (r *Router) HandleFunc(path string, pf PacketHandlerFunc, af ACKHandlerFunc) *Route {
	return r.NewRoute().Path(path).HandlerFunc(HandlerFunc{pf, af})
}

func (r Router) ServePacket(ctx sdk.Context, p PacketI, pd PacketDataI, sender PacketSender) (*sdk.Result, []byte, error) {
	for _, route := range r.routes {
		if route.Match(ctx, pd) {
			handler := route.f.ServePacket
			for i := len(r.packetMiddlewares) - 1; i >= 0; i-- {
				handler = r.packetMiddlewares[i].Middleware(handler)
			}
			return handler(ctx, p, pd, sender)
		}
	}
	if r.NotFoundHandler.PacketHandlerFunc != nil {
		return r.NotFoundHandler.PacketHandlerFunc(ctx, p, pd, sender)
	}
	return nil, nil, errors.New("route not found")
}

func (r Router) ServeACK(ctx sdk.Context, p PacketI, pd PacketDataI, ack []byte, sender PacketSender) (*sdk.Result, error) {
	for _, route := range r.routes {
		if route.Match(ctx, pd) {
			handler := route.f.ServeACK
			for i := len(r.packetMiddlewares) - 1; i >= 0; i-- {
				handler = r.ackMiddlewares[i].Middleware(handler)
			}
			return handler(ctx, p, pd, ack, sender)
		}
	}
	if r.NotFoundHandler.PacketHandlerFunc != nil {
		return r.NotFoundHandler.ACKHandlerFunc(ctx, p, pd, ack, sender)
	}
	return nil, errors.New("route not found")
}

type Route struct {
	path string
	f    Handler
}

func (r *Route) Path(path string) *Route {
	r.path = path
	return r
}

// HandlerFunc sets a handler function for the route.
func (r *Route) HandlerFunc(f HandlerFunc) *Route {
	r.f = f
	return r
}

func (r Route) Match(ctx sdk.Context, pd PacketDataI) bool {
	id, found := GetServiceID(pd.GetHeader())
	if !found {
		return false
	}
	return r.path == id
}
