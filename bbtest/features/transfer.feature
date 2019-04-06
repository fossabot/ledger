Feature: High level Transfer workflow

  Scenario: Transfer forward
    Given vault is empty
    And   tenant FORWARD is onbdoarded
    And   ledger is restarted
    And   pasive account FORWARD/OriginCredit with currency EUR exist
    And   pasive account FORWARD/OriginDebit with currency EUR exist
    And   pasive account FORWARD/Target with currency EUR exist

    When following transaction is created from tenant FORWARD
    """
      {
        "id": "forward_id",
        "transfers": [
          {
            "id": "transfer_1",
            "credit": {
              "tenant": "FORWARD",
              "name": "OriginCredit"
            },
            "debit": {
              "tenant": "FORWARD",
              "name": "OriginDebit"
            },
            "amount": "1",
            "currency": "EUR"
          },
          {
            "id": "transfer_2",
            "credit": {
              "tenant": "FORWARD",
              "name": "OriginCredit"
            },
            "debit": {
              "tenant": "FORWARD",
              "name": "OriginDebit"
            },
            "amount": "2",
            "currency": "EUR"
          }
        ]
      }
    """
    Then  request should succeed
    Then  FORWARD/OriginDebit balance should be -3 EUR
    And   FORWARD/OriginCredit balance should be 3 EUR
    And   FORWARD/Target balance should be 0 EUR

    When  forward_id transfer_1 credit side is forwarded to FORWARD/Target from tenant FORWARD
    Then  request should succeed
    Then  FORWARD/OriginDebit balance should be -3 EUR
    And   FORWARD/OriginCredit balance should be 2 EUR
    And   FORWARD/Target balance should be 1 EUR

    When  forward_id transfer_2 credit side is forwarded to FORWARD/Target from tenant FORWARD
    Then  request should succeed
    Then  FORWARD/OriginDebit balance should be -3 EUR
    And   FORWARD/OriginCredit balance should be 0 EUR
    And   FORWARD/Target balance should be 3 EUR

    When  forward_id transfer_1 debit side is forwarded to FORWARD/Target from tenant FORWARD
    Then  request should succeed
    Then  FORWARD/OriginDebit balance should be -2 EUR
    And   FORWARD/OriginCredit balance should be 0 EUR
    And   FORWARD/Target balance should be 2 EUR

    When  forward_id transfer_2 debit side is forwarded to FORWARD/Target from tenant FORWARD
    Then  request should succeed
    Then  FORWARD/OriginDebit balance should be 0 EUR
    And   FORWARD/OriginCredit balance should be 0 EUR
    And   FORWARD/Target balance should be 0 EUR
