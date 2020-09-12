package types

// PacketMiddlewareFunc is a function which receives a PacketHandlerFunc and returns another PacketHandlerFunc.
type PacketMiddlewareFunc func(PacketHandlerFunc) PacketHandlerFunc

// packetMiddleware interface is anything which implements a PacketMiddlewareFunc named Middleware.
type packetMiddleware interface {
	Middleware(handler PacketHandlerFunc) PacketHandlerFunc
}

// Middleware allows PacketHandlerFunc to implement the packetMiddleware interface.
func (mw PacketMiddlewareFunc) Middleware(handler PacketHandlerFunc) PacketHandlerFunc {
	return mw(handler)
}

// UsePacketMiddlewares appends a PacketMiddlewareFunc to the chain. Middleware can be used to intercept or otherwise modify requests and/or responses, and are executed in the order that they are applied to the Router.
func (r *Router) UsePacketMiddlewares(mwf ...PacketMiddlewareFunc) {
	for _, fn := range mwf {
		r.packetMiddlewares = append(r.packetMiddlewares, fn)
	}
}

// ACKMiddlewareFunc is a function which receives a ACKHandlerFunc and returns another ACKHandlerFunc.
type ACKMiddlewareFunc func(ACKHandlerFunc) ACKHandlerFunc

type ackMiddleware interface {
	Middleware(handler ACKHandlerFunc) ACKHandlerFunc
}

// Middleware allows MiddlewareFunc to implement the middleware interface.
func (mw ACKMiddlewareFunc) Middleware(handler ACKHandlerFunc) ACKHandlerFunc {
	return mw(handler)
}

// UseACKMiddlewares appends a ACKMiddlewareFunc to the chain. Middleware can be used to intercept or otherwise modify requests and/or responses, and are executed in the order that they are applied to the Router.
func (r *Router) UseACKMiddlewares(mwf ...ACKMiddlewareFunc) {
	for _, fn := range mwf {
		r.ackMiddlewares = append(r.ackMiddlewares, fn)
	}
}
