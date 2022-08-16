/*
 *
 * Copyright SecureKey Technologies Inc. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 * /
 *
 */

package basicv2

import "time"

// Message is message model for basic message version 2 protocol
// Reference: https://didcomm.org/basicmessage/2.0/
type Message struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"`
	Lang        string    `json:"lang"`
	CreatedTime time.Time `json:"created_time"`
	Body        struct {
		Content string `json:"content"`
	} `json:"body"`
}
