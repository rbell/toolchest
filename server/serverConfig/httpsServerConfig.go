/*
 * Copyright (c) 2024 by Randy Bell.  All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Apache Public License, version 2.0. If a copy of the APL was not distributed with this file, you can obtain one at https://www.apache.org/licenses/LICENSE-2.0.txt.
 */

package serverConfig

import (
	"crypto/tls"
)

type HttpsServerConfig struct {
	*HttpServerConfig
	tlsConfig *tls.Config
	certFile  string
	keyFile   string
}

func (c *HttpsServerConfig) SetTlsConfig(tlsConfig *tls.Config) {
	c.tlsConfig = tlsConfig
}

func (c *HttpsServerConfig) GetTlsConfig() *tls.Config {
	return c.tlsConfig
}

func (c *HttpsServerConfig) SetCertFile(certFile string) {
	c.certFile = certFile
}

func (c *HttpsServerConfig) GetCertFile() string {
	return c.certFile
}

func (c *HttpsServerConfig) SetKeyFile(keyFile string) {
	c.keyFile = keyFile
}

func (c *HttpsServerConfig) GetKeyFile() string {
	return c.keyFile
}
