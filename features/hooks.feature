Feature: hooks

  @short
  Scenario: regular test
    When I open a resource
    Then it is open
    And expect scenario name "regular test"
    And expect scenario tag "@short"
    # And step level resources are cleaned up
    # And after all tests are done resources are closed

  @long
  Scenario: rapid tests
    When I open any resources
    Then it is open
    And expect scenario name "rapid tests"
    And expect scenario tag "@long"
    # And step level resources are cleaned up
    # And after all tests are done resources are closed
