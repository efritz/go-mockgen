package mockassert

import (
	"fmt"

	"github.com/stretchr/testify/assert"
)

// TODO - document
type ArgAssertionFunc func(assert.TestingT, interface{}) bool

// TODO - document
func Called(t assert.TestingT, mockFn interface{}, msgAndArgs ...interface{}) bool {
	callCount, ok := callCount(t, mockFn, msgAndArgs...)
	if !ok {
		return false
	}
	if callCount == 0 {
		return assert.Fail(t, fmt.Sprintf("Expected %T to be called at least once", mockFn), msgAndArgs...)
	}

	return true
}

// TODO - document
func CalledOnce(t assert.TestingT, mockFn interface{}, msgAndArgs ...interface{}) bool {
	return CalledN(t, mockFn, 1, msgAndArgs...)
}

// TODO - document
func CalledN(t assert.TestingT, mockFn interface{}, n int, msgAndArgs ...interface{}) bool {
	callCount, ok := callCount(t, mockFn, msgAndArgs...)
	if !ok {
		return false
	}
	if callCount != n {
		return assert.Fail(t, fmt.Sprintf("Expected %T to be called exactly %d times", n, mockFn), msgAndArgs...)
	}

	return true
}

// TODO - document
func CalledMatching(t assert.TestingT, mockFn interface{}, assertion ArgAssertionFunc, msgAndArgs ...interface{}) bool {
	matchingCallCount, ok := matchingCallCount(t, mockFn, assertion, msgAndArgs...)
	if !ok {
		return false
	}
	if matchingCallCount == 0 {
		return assert.Fail(t, fmt.Sprintf("Expected %T to be called at least once", mockFn), msgAndArgs...)
	}
	return true
}

// TODO - document
func CalledOnceMatching(t assert.TestingT, mockFn interface{}, assertion ArgAssertionFunc, msgAndArgs ...interface{}) bool {
	return CalledNMatching(t, mockFn, 1, assertion, msgAndArgs...)
}

// TODO - document
func CalledNMatching(t assert.TestingT, mockFn interface{}, n int, assertion ArgAssertionFunc, msgAndArgs ...interface{}) bool {
	matchingCallCount, ok := matchingCallCount(t, mockFn, assertion, msgAndArgs...)
	if !ok {
		return false
	}
	if matchingCallCount != 0 {
		return assert.Fail(t, fmt.Sprintf("Expected %T to be called exactly %d times", n, mockFn), msgAndArgs...)
	}
	return true
}

// TODO - document
func callCount(t assert.TestingT, mockFn interface{}, msgAndArgs ...interface{}) (int, bool) {
	history, ok := getCallHistory(mockFn)
	if !ok {
		return 0, assert.Fail(t, fmt.Sprintf("Parameters must be a mock function description, got %T", mockFn), msgAndArgs...)
	}

	return len(history), true
}

// TODO - document
func matchingCallCount(t assert.TestingT, mockFn interface{}, assertion ArgAssertionFunc, msgAndArgs ...interface{}) (int, bool) {
	history, ok := getCallHistory(mockFn)
	if !ok {
		return 0, assert.Fail(t, fmt.Sprintf("Parameters must be a mock function description, got %T", mockFn), msgAndArgs...)
	}

	n := 0
	for _, call := range history {
		if assertion(t, call) {
			n++
		}
	}

	return n, true
}
