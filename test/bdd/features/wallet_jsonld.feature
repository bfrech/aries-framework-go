#
# Copyright SecureKey Technologies Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

@wallet_jsonld
Feature: Support for JSON-LD in universal wallet

  Scenario Outline: Issue credentials, add to wallet, query presentation, verify presentation and credentials
    Given credentials crypto algorithm "<crypto>"
    And  "Berkley" agent is running on "localhost" port "random" with http-binding did resolver url "${SIDETREE_URL}" which accepts did method "sidetree"
    And  "Alice" agent is running on "localhost" port "random" with http-binding did resolver url "${SIDETREE_URL}" which accepts did method "sidetree"
    When "Alice" creates wallet profile
    And "Alice" opens wallet
    And "Berkley" issues credentials at "2022-04-12" regarding "Master Degree" to "Alice"
    Then "Alice" adds credentials to the wallet issued by "Berkley"
    When "Vanna" queries credentials issued by "Berkley" using "QueryByExample" query type
    Then "Alice" resolves query
    And "Alice" adds presentations proof
    And "Alice" closes wallet
    Then "Vanna" receives presentations signed by "Alice" and verifies it
    And "Vanna" receives credentials from presentation signed by "Berkley" and verifies it
    Examples:
      | crypto            |
      | "Ed25519"         |
      | "ECDSA Secp256r1" |
      | "ECDSA Secp384r1" |

  Scenario Outline: Issue credentials using the wallet, add to wallet, query presentation, verify presentation and credentials
    Given credentials crypto algorithm "<crypto>"
    And  "Alice" agent is running on "localhost" port "random" with http-binding did resolver url "${SIDETREE_URL}" which accepts did method "sidetree"
    When "Alice" creates wallet profile
    And "Alice" opens wallet
    And "Alice" creates credentials at "2022-04-12" regarding "Master Degree" without proof
    And "Alice" issues credentials using the wallet
    Then "Alice" adds credentials to the wallet issued by "Alice"
    When "Vanna" queries credentials issued by "Alice" using "PresentationExchange" query type
    And "Alice" resolves query
    And "Alice" adds presentations proof
    And "Alice" closes wallet
    Then "Vanna" receives presentations signed by "Alice" and verifies it
    And "Vanna" receives credentials from presentation signed by "Alice" and verifies it
    Examples:
      | crypto            |
      | "Ed25519"         |
      | "ECDSA Secp256r1" |
      | "ECDSA Secp384r1" |

  Scenario Outline: Issue multiple credentials, add to wallet, query all, verify credentials
    Given credentials crypto algorithm "<crypto>"
    And  "Berkley" agent is running on "localhost" port "random" with http-binding did resolver url "${SIDETREE_URL}" which accepts did method "sidetree"
    And  "MIT" agent is running on "localhost" port "random" with http-binding did resolver url "${SIDETREE_URL}" which accepts did method "sidetree"
    And  "Alice" agent is running on "localhost" port "random" with http-binding did resolver url "${SIDETREE_URL}" which accepts did method "sidetree"
    When "Alice" creates wallet profile
    And "Alice" opens wallet
    And "Berkley" issues credentials at "2022-04-12" regarding "Master Degree" to "Alice"
    And "MIT" issues credentials at "2022-04-12" regarding "Bachelor Degree" to "Alice"
    Then "Alice" adds credentials to the wallet issued by "Berkley"
    And "Alice" adds credentials to the wallet issued by "MIT"
    When "Vanna" queries all credentials from "Alice"
    Then "Vanna" receives "2" credentials
    And "Alice" closes wallet
    And "Vanna" verifies credentials issued by "Berkley"
    And "Vanna" verifies credentials issued by "MIT"
    Examples:
      | crypto            |
      | "Ed25519"         |
      | "ECDSA Secp256r1" |
      | "ECDSA Secp384r1" |
