/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package api

import "github.com/hyperledger/aries-framework-go/cmd/aries-agent-mobile/pkg/wrappers/models"

// VDRController defines methods for the VDR controller.
type VDRController interface {

	// ResolveDID resolve did.
	ResolveDID(request []byte) *models.ResponseEnvelope

	// SaveDID saves the did doc to the store.
	SaveDID(request []byte) *models.ResponseEnvelope

	// CreateDID create the did doc.
	CreateDID(request []byte) *models.ResponseEnvelope

	// GetDID retrieves the did from the store.
	GetDID(request []byte) *models.ResponseEnvelope

	// GetDIDRecords retrieves the did doc containing name and didID.
	GetDIDRecords(request []byte) *models.ResponseEnvelope
}
