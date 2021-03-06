Feature: Uninstall package

  Scenario: uninstall
    Given package ledger is uninstalled
    Then  systemctl does not contain following active units
      | name        | type    |
      | ledger-rest | service |
      | ledger      | service |
      | ledger      | path    |
