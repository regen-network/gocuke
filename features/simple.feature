Feature: simple

  Scenario: eat cukes
    Given I have 5 cukes
    When I eat 3
    Then I have 2 left
