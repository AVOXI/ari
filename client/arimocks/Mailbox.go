// Code generated by mockery v1.0.0. DO NOT EDIT.

package arimocks

import ari "github.com/CyCoreSystems/ari"
import mock "github.com/stretchr/testify/mock"

// Mailbox is an autogenerated mock type for the Mailbox type
type Mailbox struct {
	mock.Mock
}

// Data provides a mock function with given fields: key
func (_m *Mailbox) Data(key *ari.Key) (*ari.MailboxData, error) {
	ret := _m.Called(key)

	var r0 *ari.MailboxData
	if rf, ok := ret.Get(0).(func(*ari.Key) *ari.MailboxData); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ari.MailboxData)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*ari.Key) error); ok {
		r1 = rf(key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: key
func (_m *Mailbox) Delete(key *ari.Key) error {
	ret := _m.Called(key)

	var r0 error
	if rf, ok := ret.Get(0).(func(*ari.Key) error); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: key
func (_m *Mailbox) Get(key *ari.Key) *ari.MailboxHandle {
	ret := _m.Called(key)

	var r0 *ari.MailboxHandle
	if rf, ok := ret.Get(0).(func(*ari.Key) *ari.MailboxHandle); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ari.MailboxHandle)
		}
	}

	return r0
}

// List provides a mock function with given fields: filter
func (_m *Mailbox) List(filter *ari.Key) ([]*ari.Key, error) {
	ret := _m.Called(filter)

	var r0 []*ari.Key
	if rf, ok := ret.Get(0).(func(*ari.Key) []*ari.Key); ok {
		r0 = rf(filter)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*ari.Key)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*ari.Key) error); ok {
		r1 = rf(filter)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: key, oldMessages, newMessages
func (_m *Mailbox) Update(key *ari.Key, oldMessages int, newMessages int) error {
	ret := _m.Called(key, oldMessages, newMessages)

	var r0 error
	if rf, ok := ret.Get(0).(func(*ari.Key, int, int) error); ok {
		r0 = rf(key, oldMessages, newMessages)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
