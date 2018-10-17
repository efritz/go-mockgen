# go-mockgen

[![GoDoc](https://godoc.org/github.com/efritz/go-mockgen?status.svg)](https://godoc.org/github.com/efritz/go-mockgen)
[![Build Status](https://secure.travis-ci.org/efritz/go-mockgen.png)](http://travis-ci.org/efritz/go-mockgen)
[![Maintainability](https://api.codeclimate.com/v1/badges/8546037d609e215de82d/maintainability)](https://codeclimate.com/github/efritz/go-mockgen/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/8546037d609e215de82d/test_coverage)](https://codeclimate.com/github/efritz/go-mockgen/test_coverage)

A mock interface code generator.

## Installation

Simply run `go get -u github.com/efritz/go-mockgen/...`.

## Binary Usage

As an example, we generate a mock for the `Client` interface from the library
[reception](https://github.com/efritz/reception). If the reception library can
be found in the GOPATH, then the following command will generate a file called
`client_mock.go` with the following content. This assumes that the current
working directory (also in the GOPATH) is called *playground*.

```bash
$ go-mockgen github.com/efritz/reception -i Client
```

```go
// Code generated by github.com/efritz/go-mockgen; DO NOT EDIT.
// This file was generated by robots at
// 2018-10-16T09:07:54-05:00
// using the command
// $ go-mockgen github.com/efritz/reception -i Client

package playground

import (
	reception "github.com/efritz/reception"
	"sync"
)

type MockClient struct {
	ListServicesFunc            func(string) ([]*reception.Service, error)
	ListServicesFuncCallHistory []ClientListServicesParamSet
	NewWatcherFunc              func(string) reception.Watcher
	NewWatcherFuncCallHistory   []ClientNewWatcherParamSet
	RegisterFunc                func(*reception.Service, func(error)) error
	RegisterFuncCallHistory     []ClientRegisterParamSet
	mutex                       sync.RWMutex
}

func NewMockClient() *MockClient {
	return &MockClient{ListServicesFunc: func(string) ([]*reception.Service, error) {
		return nil, nil
	}, NewWatcherFunc: func(string) reception.Watcher {
		return nil
	}, RegisterFunc: func(*reception.Service, func(error)) error {
		return nil
	}}
}

type ClientListServicesParamSet struct {
	Arg0 string
}

func (m *MockClient) ListServices(v0 string) ([]*reception.Service, error) {
	m.mutex.RLock()
	m.ListServicesFuncCallHistory = append(m.ListServicesFuncCallHistory, ClientListServicesParamSet{v0})
	m.mutex.RUnlock()
	r0, r1 := m.ListServicesFunc(v0)
	return r0, r1
}

func (m *MockClient) ListServicesFuncCallCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.ListServicesFuncCallHistory)
}

func (m *MockClient) ListServicesFuncCallParams() []ClientListServicesParamSet {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.ListServicesFuncCallHistory
}

type ClientNewWatcherParamSet struct {
	Arg0 string
}

func (m *MockClient) NewWatcher(v0 string) reception.Watcher {
	m.mutex.RLock()
	m.NewWatcherFuncCallHistory = append(m.NewWatcherFuncCallHistory, ClientNewWatcherParamSet{v0})
	m.mutex.RUnlock()
	r0 := m.NewWatcherFunc(v0)
	return r0
}

func (m *MockClient) NewWatcherFuncCallCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.NewWatcherFuncCallHistory)
}

func (m *MockClient) NewWatcherFuncCallParams() []ClientNewWatcherParamSet {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.NewWatcherFuncCallHistory
}

type ClientRegisterParamSet struct {
	Arg0 *reception.Service
	Arg1 func(error)
}

func (m *MockClient) Register(v0 *reception.Service, v1 func(error)) error {
	m.mutex.RLock()
	m.RegisterFuncCallHistory = append(m.RegisterFuncCallHistory, ClientRegisterParamSet{v0, v1})
	m.mutex.RUnlock()
	r0 := m.RegisterFunc(v0, v1)
	return r0
}

func (m *MockClient) RegisterFuncCallCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.RegisterFuncCallHistory)
}

func (m *MockClient) RegisterFuncCallParams() []ClientRegisterParamSet {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.RegisterFuncCallHistory
}
```

Multiple import paths can be supplied, and a mock will be generated for each
exported interface found in each package (but not subpackages). A whitelist
of packages can be supplied in order to prevent generating unnecessary code
(discussed in the flags section below).

The suggested way to generate mocks for a project is to use go-generate. Then,
when the interface is updated, running `go generate` on the package will re-generate
the mocks.

```go
package foo

//go:generate go-mockgen -f github.com/efritz/watchdog -i Retry
//go:generate go-mockgen -f github.com/efritz/overcurrent -i Breaker
```

Mocks can be generated for a dependency next to test code where it is needed, or it
can be provided as a sibling to the interface definition in the library or application.
The latter option is suggested when possible in order to keep the mocks up to date with
changing dependencies (and to reduce repeated, generated code).

### Flags

The following flags are defined by the binary.

| Name       | Short Flag | Description  |
| ---------- | ---------- | ------------ |
| package    | p          | The name of the generated package. Is the name of target directory if dirname or filename is supplied by default. |
| prefix     |            | A prefix used in the name of each mock struct. Should be TitleCase by convention. |
| interfaces | i          | A whitelist of interfaces to generate given the import paths. |
| filename   | o          | The target output file. All mocks are writen to this file. |
| dirname    | d          | The target output directory. Each mock will be written to a unique file. |
| force      | f          | Do not abort if a write to disk would overwrite an existing file. |

If neither dirname nor filename are supplied, then the generated code is printed to standard out.

## Mock Usage

Each mock can be initialized via the no-argument constructor. This is a valid
implementation of the mocked interface that returns zero values on every function
call. For testing, it may be beneficial to force a return value or side effect when
a particular method of the interface is called. This is supported by re-assigning
the function value in the mock struct to a function defined within your test. This
also allows functions to be monkeypatched in-line, capturing values from the test
method such as communication channels, call counters, and maps in which function
call arguments can be stored.

The following (stripped) example from [reception](https://github.com/efritz/reception)
uses this pattern to mock a connection to Zookeeper, returning an error when attempting
to create an ephemeral znode.

```go
func (s *ZkSuite) TestRegisterError(t sweet.T) {
    conn := NewMockZkConn()
    conn.CreateEphemeralFunc = func(path string, data []byte) error {
        return zk.ErrUnknown
    }

    client := newZkClient(conn)
    err := client.Register(&Service{Name: "s", ID: "1234"}, nil)
    Expect(err).To(Equal(zk.ErrUnknown))
    Expect(conn.CreateEphemeralFuncCallCount).To(Equal(1))
    Expect(conn.CreateEphemeralFuncCallParams[0].Arg0).To(Equal("s-1234"))
}
```

## License

Copyright (c) 2018 Eric Fritz

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
