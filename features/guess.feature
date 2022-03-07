Feature: guess the method names, parameter types and regexes of step definitions

  Scenario: no parameters
    Given the step
    """
    I have a cucumber
    """
    When we guess the step definition
    Then we get the method signature
    """
    IHaveACucumber()
    """

  Scenario: int64 parameter
    Given the step
    """
    I have 5 cucumbers
    """
    When we guess the step definition
    Then we get the method signature
    """
    IHaveCucumbers(a int64)
    """
    When we match the step
    Then we get the values
    | 5 |

  Scenario: decimal parameter with a doc string
    Given the step
    """
    I have 5.0 kilos of cucumbers
    """
    And with a doc string
    When we guess the step definition
    Then we get the method signature
    """
    IHaveKilosOfCucumbers(a *apd.Decimal, b gocuke.DocString)
    """
    When we match the step
    Then we get the values
    | 5.0 |

  Scenario: string parameter with a data table
    Given the step
    """
    I have a "red" cucumber
    """
    And with a data table
    When we guess the step definition
    Then we get the method signature
    """
    IHaveACucumber(a string, b gocuke.DataTable)
    """
    When we match the step
    Then we get the values
      | red |


  Scenario: many parameters
    Given the step
    """
    I have 10 "green" cucumbers which weigh 1.3 kilos.
    """
    When we guess the step definition
    Then we get the method signature
    """
    IHaveCucumbersWhichWeighKilos(a int64, b string, c *apd.Decimal)
    """
    When we match the step
    Then we get the values
      | 10 | green | 1.3 |
