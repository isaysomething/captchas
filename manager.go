// Copyright 2020 CleverGo. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.

package captchas

import (
	"errors"
	"strings"
)

// Option is a function that receives a pointer of manager.
type Option func(*Manager)

// CaseSensitive is an option that enable or disable case sensitive.
func CaseSensitive(v bool) Option {
	return func(m *Manager) {
		m.caseSensitive = v
	}
}

// Manager is a captchas manager.
type Manager struct {
	store         Store
	driver        Driver
	caseSensitive bool
}

// New returns a manager instance with the given store and driver.
func New(store Store, driver Driver, opts ...Option) *Manager {
	m := &Manager{
		store:         store,
		driver:        driver,
		caseSensitive: true,
	}

	for _, f := range opts {
		f(m)
	}

	return m
}

// Generate generates a new captcha and save it to store, returns an error if failed.
func (m *Manager) Generate() (Captcha, error) {
	captcha, err := m.driver.Generate()
	if err != nil {
		return nil, err
	}
	if err = m.store.Set(captcha.ID(), captcha.Answer()); err != nil {
		return nil, err
	}

	return captcha, nil
}

// Get is a shortcut of Store.Get.
func (m *Manager) Get(id string, clear bool) (string, error) {
	return m.store.Get(id, clear)
}

// Errors
var (
	ErrIncorrectCaptcha = errors.New("incorrect captcha")
	ErrExpiredCaptcha   = errors.New("expired captcha")
)

// Verify verifies whether the given actual value is equal to the
// answer of captcha, returns an error if failed.
func (m *Manager) Verify(id, actual string, clear bool) error {
	answer, err := m.store.Get(id, clear)
	if err != nil {
		return err
	}

	if m.isEqual(actual, answer) {
		return nil
	}

	return ErrIncorrectCaptcha
}

func (m *Manager) isEqual(actual, answer string) bool {
	if answer == "" || actual == "" {
		return false
	}

	if !m.caseSensitive {
		return strings.EqualFold(actual, answer)
	}

	return actual == answer
}
