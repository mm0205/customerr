package customerr

import (
	"errors"
	"fmt"
	"reflect"
)

// Tag is a type for error tag.
// the tag is intended to be  by error handler to identify error kinds.
type Tag string

// customError has an inner error and additional information.
type customError struct {
	// innerError is the original error wrapped by the `customError`.
	innerError error

	// msg is a message for this error.
	msg string

	// tags is a collection of tag.
	// This is intended to be used by error handlers to identify error kinds.
	tags []Tag
}

// `customError` must implements the `error` interface.
var _ error = (*customError)(nil)

// New returns new custom error instance.
func New(
	inner error,
	tags []Tag,
	format string,
	args ...any,
) error {
	return &customError{
		innerError: inner,
		tags:       tags,
		msg:        fmt.Sprintf(format, args...),
	}
}

// Error returns an error message.
func (ce *customError) Error() string {
	return fmt.Errorf("%s: [%w]", ce.msg, ce.innerError).Error()
}

// Is returns `true` when one of the following conditions is met.
//	* the receiver and the `target` argument are the same pointer
//	* `customError.innerError` and `target` are the same pointer
//	* `customError.innerError` `wraps` `target`.
func (ce *customError) Is(target error) bool {
	if ce == target {
		return true
	}
	return errors.Is(ce.innerError, target)
}

// As assigns the receiver or `customError.innerError`
// or the error that `customError.innerError` recursively wraps in `target` and returns `true`.
// If none of the above is assignable to `target`, return `false`.
func (ce *customError) As(target any) bool {
	val := reflect.ValueOf(target)
	valKind := val.Kind()

	if valKind != reflect.Pointer ||
		val.IsNil() {
		panic("errors: target must be a non-nil pointer")
	}

	targetType := val.Type().Elem()
	if targetType.Kind() != reflect.Interface &&
		!targetType.Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		panic("errors: *target must be interface or implement error")
	}
	if reflect.TypeOf(ce).AssignableTo(targetType) {
		val.Elem().Set(reflect.ValueOf(ce))
	}

	//goland:noinspection GoErrorsAs
	return errors.As(ce.innerError, target)
}

// Unwrap returns `customError.innerError`
func (ce *customError) Unwrap() error {
	return ce.innerError
}

// IsCustomErr returns `true` if the `err` is a Custom Error.
func IsCustomErr(err error) bool {
	_, ok := err.(*customError)
	return ok
}

// HasTag returns `true` if the `err` is a custom error
// or the error that is wrapped
func HasTag(err error, tag Tag) bool {
	tags := Tags(err)
	for _, t := range tags {
		if t == tag {
			return true
		}
	}
	return false
}

// Tags returns all tags attached to the `err` or
// attached to the errors recursively wrapped by the `err`.
func Tags(err error) []Tag {
	var ce *customError
	ok := errors.As(err, &ce)
	var result []Tag

	for ok {
		if result == nil {
			result = make([]Tag, len((*ce).tags))
		}

		for _, t := range (*ce).tags {
			result = append(result, t)
		}
		ok = errors.As(ce.innerError, &ce)
	}
	return result
}
