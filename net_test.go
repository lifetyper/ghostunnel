/*-
 * Copyright 2015 Square Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeAddress(t *testing.T) {
	addr4, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8080")
	assert.Nil(t, err, "must resolve valid address")

	addr6, err := net.ResolveTCPAddr("tcp", "[::1]:8080")
	assert.Nil(t, err, "must resolve valid address")

	net, _ := decodeAddress(addr4)
	if net != "tcp4" {
		t.Errorf("decoding 127.0.0.1 network gave %s, should be tcp4", net)
	}

	net, _ = decodeAddress(addr6)
	if net != "tcp6" {
		t.Errorf("decoding 127.0.0.1 network gave %s, should be tcp6", net)
	}
}

func TestParseUnixOrTcpAddress(t *testing.T) {
	network, address, host, _ := parseUnixOrTCPAddress("unix:/tmp/foo")
	if network != "unix" {
		t.Errorf("unexpected network: %s", network)
	}
	if address != "/tmp/foo" {
		t.Errorf("unexpected address: %s", address)
	}
	if host != "" {
		t.Errorf("unexpected host: %s", host)
	}

	network, address, host, _ = parseUnixOrTCPAddress("localhost:8080")
	// note: ipv6 test is probably fragile, we don't expand ::1.
	if network != "tcp4" && network != "tcp6" {
		t.Errorf("unexpected network: %s", network)
	}
	if address != "127.0.0.1:8080" && address != "[::1]:8080" {
		t.Errorf("unexpected address: %s", address)
	}
	if host != "localhost" {
		t.Errorf("unexpected host: %s", host)
	}

	_, _, _, err := parseUnixOrTCPAddress("localhost")
	assert.NotNil(t, err, "was able to parse invalid host/port")

	_, _, _, err = parseUnixOrTCPAddress("256.256.256.256:99999")
	assert.NotNil(t, err, "was able to parse invalid host/port")
}
