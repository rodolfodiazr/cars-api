package contextkeys

// RequestIDKeyType prevents context key collisions.
type RequestIDKeyType string

const RequestIDKey RequestIDKeyType = "requestID"
