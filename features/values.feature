Feature: convert parameter strings to step definition values

  Scenario: any int64 value
    Given any int64 string
    When when I convert it to an int64
    Then I get back the original value

  Scenario: any decimal value
    Given any decimal string
    When when I convert it to a decimal
    Then I get back the original value

  Scenario: any big integer value
    Given any big integer string
    When when I convert it to a big integer
    Then I get back the original value


  Scenario Outline: some int64 values
    Given an int64 <x>
    When when I convert it to an int64
    Then I get back the original value

    Examples:
    | x |
    | 0 |
    | -23572034732 |
    | 5702482349215 |


  Scenario Outline: some decimal values
    Given a decimal <x>
    When when I convert it to a decimal
    Then I get back the original value

    Examples:
      | x |
      | 0.0 |
      | -523482300347322357234124.2357123129 |
      | 0.3571353723814251325367 |
      | 957239240672175829402.9867463518239606482 |


  Scenario Outline: some big integer values
    Given a big integer <x>
    When when I convert it to a big integer
    Then I get back the original value

    Examples:
      | x |
      | 0 |
      | -12345678987654321137635923752357823982473278236482372423523 |
      | 10                                                           |
      | 3287634093469384572349823569834674309582349823657348767359823754235 |
