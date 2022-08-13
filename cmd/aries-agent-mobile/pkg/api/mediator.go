/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package api

import "github.com/hyperledger/aries-framework-go/cmd/aries-agent-mobile/pkg/wrappers/models"

// MediatorController defines methods for the Mediator controller.
type MediatorController interface {

	// Register registers the agent with the router.
	Register(request []byte) *models.ResponseEnvelope

	// Unregister unregisters the agent with the router.
	Unregister(request []byte) *models.ResponseEnvelope

	// Connections returns router`s connections.
	Connections(request []byte) *models.ResponseEnvelope

	// Reconnect sends noop message to given mediator connection to re-establish network connection
	Reconnect(request []byte) *models.ResponseEnvelope

	// ReconnectAll Reconnect sends noop message to all mediator connection to re-establish network connections.
	ReconnectAll(request []byte) *models.ResponseEnvelope

	// Status returns details about pending messages for given connection.
	Status(request []byte) *models.ResponseEnvelope

	// BatchPickup dispatches pending messages for given connection.
	BatchPickup(request []byte) *models.ResponseEnvelope

	// RegisterKey registers a new key with the mediator
	RegisterKey(request []byte) *models.ResponseEnvelope
}
