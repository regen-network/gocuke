Feature: cleanup

  Scenario: regular test
    When I open a resource
    Then it is open
    # Then after all tests is done it is closed

  Scenario: rapid tests
    When I open any resources
    Then it is open
    # Then after all tests is done they are closed
