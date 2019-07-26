// Copyright (c) 2019 Faye Amacker. All rights reserved.
// Use of this source code is governed by Apache License 2.0 found in the LICENSE file.

package webauthn

import (
	"errors"
	"net/url"
	"strconv"
)

// Config represents Relying Party settings used to create attestation and assertion options.
// Zero value Config is not valid.
type Config struct {
	ChallengeLength         int
	Timeout                 uint64
	RPID                    string
	RPName                  string
	RPIcon                  string
	AuthenticatorAttachment AuthenticatorAttachment
	ResidentKey             ResidentKeyRequirement
	UserVerification        UserVerificationRequirement
	Attestation             AttestationConveyancePreference
	CredentialAlgs          []int
}

const (
	challengeMinLength = 16
	challengeMaxLength = 64
)

// Valid checks Config settings and returns error if it is invalid.
func (c *Config) Valid() error {
	if c.RPName == "" {
		return errors.New("rp name is required")
	}
	if c.RPID == "" {
		return errors.New("rp id is required")
	}
	if _, err := url.Parse(c.RPID); err != nil {
		return errors.New("rp id " + c.RPID + " is not a valid URI: " + err.Error())
	}
	if c.Timeout <= 0 {
		return errors.New("timeout must be a positive number")
	}
	if c.ChallengeLength < challengeMinLength {
		return errors.New("challenge must be at least " + strconv.Itoa(challengeMinLength) + " bytes long")
	}
	if c.ChallengeLength > challengeMaxLength {
		return errors.New("challenge must be no more than" + strconv.Itoa(challengeMaxLength) + " bytes long")
	}
	if c.AuthenticatorAttachment != "" &&
		c.AuthenticatorAttachment != AuthenticatorPlatform &&
		c.AuthenticatorAttachment != AuthenticatorCrossPlatform {
		return errors.New("authenticator attachment must be \"\", \"platform\", or \"cross-platform\"")
	}
	if c.ResidentKey != ResidentKeyRequired &&
		c.ResidentKey != ResidentKeyPreferred &&
		c.ResidentKey != ResidentKeyDiscouraged {
		return errors.New("resident key must be \"required\", \"preferred\", or \"discouraged\"")
	}
	if c.UserVerification != UserVerificationRequired &&
		c.UserVerification != UserVerificationPreferred &&
		c.UserVerification != UserVerificationDiscouraged {
		return errors.New("user verification must be \"required\", \"preferred\", or \"discouraged\"")
	}
	if c.Attestation != AttestationNone &&
		c.Attestation != AttestationIndirect &&
		c.Attestation != AttestationDirect {
		return errors.New("attestation must be \"none\", \"indirect\", or \"direct\"")
	}
	if len(c.CredentialAlgs) == 0 {
		return errors.New("there must be at least one credential algorithm")
	}
	for _, alg := range c.CredentialAlgs {
		if !signatureAlgorithmRegistered(alg) {
			return errors.New("credential algorithm " + strconv.Itoa(alg) + " is not registered")
		}
	}

	return nil
}
