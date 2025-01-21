package bus

// GetParam functions like url.Values.Get(). It returns the first value
// matching that key or empty string.
func (wcr *WebhookCallRequest) GetParam(key string) string {
	values := wcr.Params[key]
	if values == nil || len(values.Values) == 0 {
		return ""
	}
	return values.Values[0]
}
