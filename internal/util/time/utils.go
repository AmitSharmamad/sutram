package time

import "time"

// TimeOnlyToTime attempts to convert a pointer to a custom TimeOnly type (or time.Time) to time.Time
// Accepts interface{} to allow for different possible types
func TimeOnlyToTime(val interface{}) time.Time {
	if val == nil {
		return time.Time{}
	}
	// Try direct time.Time pointer
	if t, ok := val.(*time.Time); ok && t != nil {
		return *t
	}
	// Try direct time.Time value
	if t, ok := val.(time.Time); ok {
		return t
	}
	// Try pointer to custom type with ToTime() method
	if tm, ok := val.(interface{ ToTime() time.Time }); ok {
		return tm.ToTime()
	}
	// Try value of custom type with ToTime() method
	if tm, ok := val.(interface{ ToTime() time.Time }); ok {
		return tm.ToTime()
	}
	// Try dereferencing pointer and checking again
	switch v := val.(type) {
	case *interface{}:
		return TimeOnlyToTime(*v)
	}
	return time.Time{}
}
