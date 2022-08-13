/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package api

import "github.com/hyperledger/aries-framework-go/cmd/aries-agent-mobile/pkg/wrappers/models"

// OutOfBandController defines methods for the out-of-band protocol controller.
type OutOfBandController interface {
	// CreateInvitation creates and saves an out-of-band invitation.
	CreateInvitation(request []byte) *models.ResponseEnvelope

	// AcceptInvitation from another agent and return the ID of the new connection records.
	AcceptInvitation(request []byte) *models.ResponseEnvelope

	// Actions returns pending actions that have not yet to be executed or canceled.
	Actions(request []byte) *models.ResponseEnvelope

	// ActionContinue allows continuing with the protocol after an action event was triggered.
	ActionContinue(request []byte) *models.ResponseEnvelope

	// ActionStop stops the protocol after an action event was triggered.
	ActionStop(request []byte) *models.ResponseEnvelope
}
