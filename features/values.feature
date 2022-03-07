Feature: convert parameter strings to step definition values

  Scenario: int64 values
    Given any int64 string
    When when I convert it to an int64
    Then I get back the original value

  Scenario: decimal values
    Given any decimal string
    When when I convert it to a decimal
    Then I get back the original value

  Scenario: big integer values
    Given any big integer string
    When when I convert it to a big integer
    Then I get back the original value
