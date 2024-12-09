/*
Ory Identities API

This is the API specification for Ory Identities with features such as registration, login, recovery, account verification, profile settings, password reset, identity management, session management, email and sms delivery, and more. 

API version: v1.2.1
Contact: office@ory.sh
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package client

import (
	"encoding/json"
	"fmt"
)

// checks if the UpdateLoginFlowWithLookupSecretMethod type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &UpdateLoginFlowWithLookupSecretMethod{}

// UpdateLoginFlowWithLookupSecretMethod Update Login Flow with Lookup Secret Method
type UpdateLoginFlowWithLookupSecretMethod struct {
	// Sending the anti-csrf token is only required for browser login flows.
	CsrfToken *string `json:"csrf_token,omitempty"`
	// The lookup secret.
	LookupSecret string `json:"lookup_secret"`
	// Method should be set to \"lookup_secret\" when logging in using the lookup_secret strategy.
	Method string `json:"method"`
	AdditionalProperties map[string]interface{}
}

type _UpdateLoginFlowWithLookupSecretMethod UpdateLoginFlowWithLookupSecretMethod

// NewUpdateLoginFlowWithLookupSecretMethod instantiates a new UpdateLoginFlowWithLookupSecretMethod object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewUpdateLoginFlowWithLookupSecretMethod(lookupSecret string, method string) *UpdateLoginFlowWithLookupSecretMethod {
	this := UpdateLoginFlowWithLookupSecretMethod{}
	this.LookupSecret = lookupSecret
	this.Method = method
	return &this
}

// NewUpdateLoginFlowWithLookupSecretMethodWithDefaults instantiates a new UpdateLoginFlowWithLookupSecretMethod object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewUpdateLoginFlowWithLookupSecretMethodWithDefaults() *UpdateLoginFlowWithLookupSecretMethod {
	this := UpdateLoginFlowWithLookupSecretMethod{}
	return &this
}

// GetCsrfToken returns the CsrfToken field value if set, zero value otherwise.
func (o *UpdateLoginFlowWithLookupSecretMethod) GetCsrfToken() string {
	if o == nil || IsNil(o.CsrfToken) {
		var ret string
		return ret
	}
	return *o.CsrfToken
}

// GetCsrfTokenOk returns a tuple with the CsrfToken field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *UpdateLoginFlowWithLookupSecretMethod) GetCsrfTokenOk() (*string, bool) {
	if o == nil || IsNil(o.CsrfToken) {
		return nil, false
	}
	return o.CsrfToken, true
}

// HasCsrfToken returns a boolean if a field has been set.
func (o *UpdateLoginFlowWithLookupSecretMethod) HasCsrfToken() bool {
	if o != nil && !IsNil(o.CsrfToken) {
		return true
	}

	return false
}

// SetCsrfToken gets a reference to the given string and assigns it to the CsrfToken field.
func (o *UpdateLoginFlowWithLookupSecretMethod) SetCsrfToken(v string) {
	o.CsrfToken = &v
}

// GetLookupSecret returns the LookupSecret field value
func (o *UpdateLoginFlowWithLookupSecretMethod) GetLookupSecret() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.LookupSecret
}

// GetLookupSecretOk returns a tuple with the LookupSecret field value
// and a boolean to check if the value has been set.
func (o *UpdateLoginFlowWithLookupSecretMethod) GetLookupSecretOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.LookupSecret, true
}

// SetLookupSecret sets field value
func (o *UpdateLoginFlowWithLookupSecretMethod) SetLookupSecret(v string) {
	o.LookupSecret = v
}

// GetMethod returns the Method field value
func (o *UpdateLoginFlowWithLookupSecretMethod) GetMethod() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Method
}

// GetMethodOk returns a tuple with the Method field value
// and a boolean to check if the value has been set.
func (o *UpdateLoginFlowWithLookupSecretMethod) GetMethodOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Method, true
}

// SetMethod sets field value
func (o *UpdateLoginFlowWithLookupSecretMethod) SetMethod(v string) {
	o.Method = v
}

func (o UpdateLoginFlowWithLookupSecretMethod) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o UpdateLoginFlowWithLookupSecretMethod) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.CsrfToken) {
		toSerialize["csrf_token"] = o.CsrfToken
	}
	toSerialize["lookup_secret"] = o.LookupSecret
	toSerialize["method"] = o.Method

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *UpdateLoginFlowWithLookupSecretMethod) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"lookup_secret",
		"method",
	}

	allProperties := make(map[string]interface{})

	err = json.Unmarshal(data, &allProperties)

	if err != nil {
		return err;
	}

	for _, requiredProperty := range(requiredProperties) {
		if _, exists := allProperties[requiredProperty]; !exists {
			return fmt.Errorf("no value given for required property %v", requiredProperty)
		}
	}

	varUpdateLoginFlowWithLookupSecretMethod := _UpdateLoginFlowWithLookupSecretMethod{}

	err = json.Unmarshal(data, &varUpdateLoginFlowWithLookupSecretMethod)

	if err != nil {
		return err
	}

	*o = UpdateLoginFlowWithLookupSecretMethod(varUpdateLoginFlowWithLookupSecretMethod)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "csrf_token")
		delete(additionalProperties, "lookup_secret")
		delete(additionalProperties, "method")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableUpdateLoginFlowWithLookupSecretMethod struct {
	value *UpdateLoginFlowWithLookupSecretMethod
	isSet bool
}

func (v NullableUpdateLoginFlowWithLookupSecretMethod) Get() *UpdateLoginFlowWithLookupSecretMethod {
	return v.value
}

func (v *NullableUpdateLoginFlowWithLookupSecretMethod) Set(val *UpdateLoginFlowWithLookupSecretMethod) {
	v.value = val
	v.isSet = true
}

func (v NullableUpdateLoginFlowWithLookupSecretMethod) IsSet() bool {
	return v.isSet
}

func (v *NullableUpdateLoginFlowWithLookupSecretMethod) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableUpdateLoginFlowWithLookupSecretMethod(val *UpdateLoginFlowWithLookupSecretMethod) *NullableUpdateLoginFlowWithLookupSecretMethod {
	return &NullableUpdateLoginFlowWithLookupSecretMethod{value: val, isSet: true}
}

func (v NullableUpdateLoginFlowWithLookupSecretMethod) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableUpdateLoginFlowWithLookupSecretMethod) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

