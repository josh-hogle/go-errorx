package errorx

import (
	"fmt"
	"runtime"
	"time"
)

// Error represents an extension to the standard error interface by adding the ability to include an error code,
// nested errors, and any attributes associated with the error.
type Error interface {
	// Attrs should return any additional attributes which may be used when logging the error.
	Attrs() map[string]any

	// Code should return the corresponding error code.
	Code() int

	// Error should return the string version of the error.
	Error() string

	// File should return the name of the file where the error occurred if caller information is included.
	File() string

	// InternalError should return the internal standard error object if there is one.
	InternalError() error

	// Line should return the line number where the error occurred if caller information is included.
	Line() int

	// Method should return the name of the function where the error occurred if caller information is included.
	Method() string

	// NestedErrors should return the list of nested errors if there are any.
	NestedErrors() []Error
}

// BaseError is a base error object that mostly implements the Error interface and can be used to more easily compose
// more type-specific errors. Composed objects must implement the Error() function.
//
// Do not create this object directly. Use NewBaseError() to construct a new object so its values are initialized
// properly.
type BaseError struct {
	errAttrs   map[string]any
	errCode    int
	err        error
	file       string
	method     string
	line       int
	nestedErrs []Error
	pc         uintptr
}

// NewBaseError returns a new BaseError object without caller information included.
func NewBaseError(code int, err error) *BaseError {
	if err == nil {
		err = fmt.Errorf("an unknown error occurred (code=%d)", code)
	}
	return &BaseError{
		errAttrs:   map[string]any{},
		errCode:    code,
		err:        err,
		line:       0,
		nestedErrs: []Error{},
		pc:         0,
	}
}

// NewBaseError returns a new BaseError object with caller information included.
func NewBaseErrorWithCaller(code int, err error, skip int) *BaseError {
	if err == nil {
		err = fmt.Errorf("an unknown error occurred (code=%d)", code)
	}
	// skip Caller + this function + caller (which will be the actual NewError()... function)
	pc, file, line, ok := runtime.Caller(skip + 2)
	if ok {
		f := runtime.FuncForPC(pc)
		return &BaseError{
			errAttrs:   map[string]any{},
			errCode:    code,
			err:        err,
			file:       file,
			line:       line,
			method:     f.Name(),
			nestedErrs: []Error{},
			pc:         pc,
		}
	}
	return &BaseError{
		errAttrs:   map[string]any{},
		errCode:    code,
		err:        err,
		file:       "????",
		line:       0,
		method:     "????",
		nestedErrs: []Error{},
		pc:         0,
	}
}

// Append appends one or more non-nil errors to the end of the list of nested errors associated with this error.
//
// Any nil errors passed are ignored.
func (b *BaseError) Append(errs ...Error) {
	for _, err := range errs {
		if err != nil {
			b.nestedErrs = append(b.nestedErrs, err)
		}
	}
}

// Attr returns the value of the attribute with the given key if it exists.
func (b *BaseError) Attr(key string) (any, error) {
	v, ok := b.errAttrs[key]
	if !ok {
		return "", fmt.Errorf("%s: attribute not found", key)
	}
	return v, nil
}

// AttrDuration returns the value of the attribute with the given key (if it exists) as a time duration.
func (b *BaseError) AttrDuration(key string) (time.Duration, error) {
	v, err := b.Attr(key)
	if err != nil {
		return time.Duration(0), err
	}
	t, ok := v.(time.Duration)
	if !ok {
		return time.Duration(0), fmt.Errorf("%s: cannot convert attribute value to time.Duration", key)
	}
	return t, nil
}

// AttrInt returns the value of the attribute with the given key (if it exists) as an integer if it exists.
func (b *BaseError) AttrInt(key string) (int, error) {
	v, err := b.Attr(key)
	if err != nil {
		return 0, err
	}
	i, ok := v.(int)
	if !ok {
		return 0, fmt.Errorf("%s: cannot convert attribute value to an int", key)
	}
	return i, nil
}

// AttrInt64 returns the value of the attribute with the given key (if it exists) as a 64-bit integer.
func (b *BaseError) AttrInt64(key string) (int64, error) {
	v, err := b.Attr(key)
	if err != nil {
		return 0, err
	}
	i, ok := v.(int64)
	if !ok {
		return 0, fmt.Errorf("%s: cannot convert attribute value to an int64", key)
	}
	return i, nil
}

// Attrs returns all of the attributes associated with the error.
//
// If this object has no attributes defined, an empty map is returned.
func (b *BaseError) Attrs() map[string]any {
	return b.errAttrs
}

// AttrString returns the value of the attribute with the given key (if it exists) as a string.
func (b *BaseError) AttrString(key string) (string, error) {
	v, err := b.Attr(key)
	if err != nil {
		return "", err
	}
	str, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("%s: cannot convert attribute value to a string", key)
	}
	return str, nil
}

// AttrTime returns the value of the attribute with the given key (if it exists) as a time object.
func (b *BaseError) AttrTime(key string) (time.Time, error) {
	v, err := b.Attr(key)
	if err != nil {
		return time.Time{}, err
	}
	t, ok := v.(time.Time)
	if !ok {
		return time.Time{}, fmt.Errorf("%s: cannot convert attribute value to time.Time", key)
	}
	return t, nil
}

// AttrUint returns the value of the attribute with the given key (if it exists) as an unsigned integer.
func (b *BaseError) AttrUint(key string) (uint, error) {
	v, err := b.Attr(key)
	if err != nil {
		return 0, err
	}
	i, ok := v.(uint)
	if !ok {
		return 0, fmt.Errorf("%s: cannot convert attribute value to an uint", key)
	}
	return i, nil
}

// AttrUint64 returns the value of the attribute with the given key (if it exists) as an unsigned 64-bit integer.
func (b *BaseError) AttrUint64(key string) (uint64, error) {
	v, err := b.Attr(key)
	if err != nil {
		return 0, err
	}
	i, ok := v.(uint64)
	if !ok {
		return 0, fmt.Errorf("%s: cannot convert attribute value to an uint64", key)
	}
	return i, nil
}

// Code returns the error code.
func (b *BaseError) Code() int {
	return b.errCode
}

// Error returns the corresponding error message.
func (b *BaseError) Error() string {
	return fmt.Sprintf("error: %s", b.InternalError().Error())
}

// File returns the name of the file where the error occurred if caller information is included.
func (b *BaseError) File() string {
	return b.file
}

// InternalError returns the standard error associated with the object.
//
// The error returned by this function is guaranteed to never be `nil`.
func (b *BaseError) InternalError() error {
	return b.err
}

// Line returns the line number where the error occurred if caller information is included.
func (b *BaseError) Line() int {
	return b.line
}

// Method returns the name of the function where the error occurred if caller information is included.
func (b *BaseError) Method() string {
	return b.method
}

// NestedErrors returns the list of errors that were generated by a call to another function.
//
// If there are no nested errors, an empty list is returned.
func (b *BaseError) NestedErrors() []Error {
	return b.nestedErrs
}

// WithAttr adds the given key/value pair to the list of attributes associated with this error and
// returns itself
func (b *BaseError) WithAttr(attrKey string, attrValue any) *BaseError {
	b.errAttrs[attrKey] = attrValue
	return b
}

// WithAttrs adds each of the given key/value pairs to the list of attributes associated with this error
// and returns itself.
func (b *BaseError) WithAttrs(attrs map[string]any) *BaseError {
	for k, v := range attrs {
		b.errAttrs[k] = v
	}
	return b
}
