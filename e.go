package e

import "encoding/json"

type wrapper struct {
	cause   error
	message string
	context map[string]interface{}
}

type mapper struct {
	Message string                 `json:"message"`
	Context map[string]interface{} `json:"context"`
}

func (w *wrapper) Error() string {
	return w.message
}

func (w *wrapper) Unwrap() error {
	return w.cause
}

func (w *wrapper) Map() mapper {
	return mapper{
		Message: w.message,
		Context: w.context,
	}
}

func (w *wrapper) JSON() ([]byte, error) {
	buf, err := json.Marshal(w.Map())

	return buf, Wrap(err, "failed to marshal error to json", map[string]interface{}{
		"error": w.Error(),
	})
}

func Is(err error) bool {
	if err == nil {
		return false
	}

	_, ok := err.(*wrapper)

	return ok
}

func As(err error) *wrapper {
	return err.(*wrapper)
}

func New(msg string, ctx map[string]interface{}) error {
	return &wrapper{
		cause:   nil,
		message: msg,
		context: ctx,
	}
}

func Wrap(err error, msg string, ctx map[string]interface{}) error {
	if err == nil {
		return nil
	}

	return &wrapper{
		cause:   err,
		message: msg,
		context: ctx,
	}
}

func Cause(err error) error {
	if !Is(err) {
		return nil
	}

	return As(err).cause
}

func Root(err error) error {
	if !Is(err) {
		return err
	}

	return Root(As(err).cause)
}

func Primary(err error) error {
	if !Is(err) {
		return nil
	}

	w := As(err)

	if !Is(w.cause) {
		return err
	}

	return Primary(w.cause)
}

func Context(err error) map[string]interface{} {
	if !Is(err) {
		return nil
	}

	return As(err).context
}

func Depth(err error) int {
	if err == nil {
		return 0
	}

	if !Is(err) {
		return 1
	}

	return 1 + Depth(As(err).cause)
}

func Array(err error) []error {
	if err == nil {
		return nil
	}

	arr := []error{
		err,
	}

	if !Is(err) {
		return arr
	}

	return append(arr, Array(As(err).cause)...)
}

func Map(err error) []mapper {
	arr := Array(err)

	if arr == nil {
		return nil
	}

	out := make([]mapper, len(arr))

	for i, err := range arr {
		mpr := mapper{
			Message: err.Error(),
		}

		if Is(err) {
			mpr.Context = As(err).context
		}

		out[i] = mpr
	}

	return out
}

func JSON(err error) ([]byte, error) {
	if err == nil {
		return nil, nil
	}

	buf, jme := json.Marshal(Map(err))

	return buf, Wrap(jme, "failed to marshal error stack into json", map[string]interface{}{
		"error": err.Error(),
	})
}
