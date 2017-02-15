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

package services

import (
	"encoding/json"
	"fmt"

	"github.com/gravitational/trace"
)

type ClusterAuthPreference interface {
	GetType() string
	SetType(string)
	GetSecondFactor() string
	SetSecondFactor(string)
	GetOIDC() string
	SetOIDC(string)
	SetU2F() string
	GetU2F(string)
	GetLDAP() string
	SetLDAP(string)
	Check() error
}

func NewClusterAuthPreference(spec ClusterAuthPreferenceV2) (ClusterAuthPreference, error) {
	clusterAuthPreference := ClusterAuthPreferenceV2{
		Kind:    KindClusterAuthPreference,
		Version: V2,
		Metadata: Metadata{
			Name:      MetaNameClusterAuthPreference,
			Namespace: defaults.Namespace,
		},
		Spec: spec,
	}

	err := clusterAuthPreference.Check()
	if err != nil {
		return nil, trace.Wrap(err)
	}

	return &clusterAuthPreference, nil
}

type ClusterAuthPreferenceV2 struct {
	Kind     string                      `json:"kind"`
	Version  string                      `json:"version"`
	Metadata Metadata                    `json:"metadata"`
	Spec     ClusterAuthPreferenceSpecV2 `json:"spec"`
}

type ClusterAuthPreferenceSpecV2 struct {
	Type         string `json:"type"`
	SecondFactor string `json:"second_factor,omitempty"`
	OIDC         string `json:"oidc,omitempty"`
	U2F          string `json:"u2f,omitempty"`
	LDAP         string `json:"ldap,omitempty"`
}

func (c *ClusterAuthPreferenceV2) GetType() string {
	return c.Spec.Type
}

func (c *ClusterAuthPreferenceV2) SetType(t string) {
	c.Spec.Type = t
}

func (c *ClusterAuthPreferenceV2) GetSecondFactor() string {
	return c.Spec.SecondFactor
}

func (c *ClusterAuthPreferenceV2) SetSecondFactor(s string) {
	c.Spec.SecondFactor = s
}

func (c *ClusterAuthPreferenceV2) GetOIDC() string {
	return c.Spec.OIDC
}

func (c *ClusterAuthPreferenceV2) SetOIDC(o string) {
	c.Spec.OIDC = o
}

func (c *ClusterAuthPreferenceV2) GetU2F() string {
	return c.Spec.U2F
}

func (c *ClusterAuthPreferenceV2) SetU2F(u string) {
	c.Spec.U2F = u
}

func (c *ClusterAuthPreferenceV2) GetLDAP() string {
	return c.Spec.LDAP
}

func (c *ClusterAuthPreferenceV2) SetLDAP(l string) {
	c.Spec.LDAP = l
}

func (c *ClusterAuthPreferenceV2) Check() error {
	return nil
}

const ClusterAuthPreferenceSpecSchemaTemplate = `{
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "type": {"type": "string"},
    "second_factor": {"type": "string"},
    "oidc": {"type": "string"},
    "u2f": {"type": "string"},
    "ldap": {"type": "string"}%v
  }
}`

func GetClusterAuthPreferenceSchema(extensionSchema string) string {
	var clusterAuthPreferenceSchema string
	if clusterAuthPreferenceSchema == "" {
		clusterAuthPreferenceSchema = fmt.Sprintf(ClusterAuthPreferenceSpecSchemaTemplate, "")
	} else {
		clusterAuthPreferenceSchema = fmt.Sprintf(ClusterAuthPreferenceSpecSchemaTemplate, ","+extensionSchema)
	}
	return fmt.Sprintf(V2SchemaTemplate, MetadataSchema, clusterAuthPreferenceSchema)
}

type ClusterAuthPreferenceMarshaler interface {
	UnmarshalClusterAuthPreference(bytes []byte) (ClusterAuthPreference, error)
	MarshalClusterAuthPreference(c ClusterAuthPreference, opts ...MarshalOption) ([]byte, error)
}

var clusterAuthPreferenceMarshaler ClusterAuthPreferenceMarshaler = &TeleportClusterAuthPreferenceMarshaler{}

type TeleportClusterAuthPreferenceMarshaler struct{}

func (t *TeleportClusterAuthPreferenceMarshaler) UnmarshalClusterAuthPreference(bytes []byte) (ClusterAuthPreference, error) {
	var clusterAuthPreference ClusterAuthPreferenceV2

	if len(data) == 0 {
		return nil, trace.BadParameter("missing resource data")
	}

	err := utils.UnmarshalWithSchema(GetClusterAuthPreferenceSchema(""), &clusterAuthPreference, data)
	if err != nil {
		return nil, trace.BadParameter(err.Error())
	}

	return &clusterAuthPreference, nil
}

func (t *TeleportClusterAuthPreferenceMarshaler) MarshalClusterAuthPreference(c ClusterAuthPreference, opts ...MarshalOption) ([]byte, error) {
	return json.Marshal(c)
}
