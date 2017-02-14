/*
Copyright 2017 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package config

import (
	"encoding/base64"
	"fmt"

	"gopkg.in/check.v1"
)

type FileTestSuite struct {
}

var _ = check.Suite(&FileTestSuite{})
var _ = fmt.Printf

func (s *FileTestSuite) SetUpSuite(c *check.C) {
}

func (s *FileTestSuite) TearDownSuite(c *check.C) {
}

func (s *FileTestSuite) SetUpTest(c *check.C) {
}

// TestLegacySection ensures we continue to parse and correctly load deprecated
// OIDC connector and U2F authentication configuration.
func (s *FileTestSuite) TestLegacySection(c *check.C) {
	encodedLegacyAuthenticationSection := base64.StdEncoding.EncodeToString([]byte(LegacyAuthenticationSection))

	// read config into struct
	fc, err := ReadFromString(encodedLegacyAuthenticationSection)
	c.Assert(err, check.IsNil)

	// validate oidc connector
	c.Assert(fc.Auth.OIDCConnectors, check.HasLen, 1)
	c.Assert(fc.Auth.OIDCConnectors[0].ID, check.Equals, "google")
	c.Assert(fc.Auth.OIDCConnectors[0].RedirectURL, check.Equals, "https://localhost:3080/v1/webapi/oidc/callback")
	c.Assert(fc.Auth.OIDCConnectors[0].ClientID, check.Equals, "id-from-google.apps.googleusercontent.com")
	c.Assert(fc.Auth.OIDCConnectors[0].ClientSecret, check.Equals, "secret-key-from-google")
	c.Assert(fc.Auth.OIDCConnectors[0].IssuerURL, check.Equals, "https://accounts.google.com")

	// validate u2f
	c.Assert(fc.Auth.U2F.EnabledFlag, check.Equals, "yes")
	c.Assert(fc.Auth.U2F.AppID, check.Equals, "https://graviton:3080")
	c.Assert(fc.Auth.U2F.Facets, check.HasLen, 1)
	c.Assert(fc.Auth.U2F.Facets[0], check.Equals, "https://graviton:3080")
}
