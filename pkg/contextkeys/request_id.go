package contextkeys

// RequestIDKeyType is the context key type used to store
// request IDs in a context.Context.
type RequestIDKeyType struct{}

// RequestIDKey is the context key used to store and retrieve
// the request ID associated with an HTTP request.
var RequestIDKey = RequestIDKeyType{}
