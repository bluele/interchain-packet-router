package types

const (
	HeaderKeyServiceID = "@service"
)

// GetServiceID returns a service ID from a given header
func GetServiceID(h HeaderI) (string, bool) {
	v, ok := h.Get(HeaderKeyServiceID)
	if !ok {
		return "", false
	}
	return string(v), true
}

// SetServiceID sets ID to a given header
func SetServiceID(h HeaderI, id string) {
	h.Set(HeaderKeyServiceID, []byte(id))
}
