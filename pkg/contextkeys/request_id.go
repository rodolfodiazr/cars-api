package contextkeys

// RequestIDKeyType is the context key type used to store
// request IDs in a context.Context.
type RequestIDKeyType struct{}

var RequestIDKey = RequestIDKeyType{}
