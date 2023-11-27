Feature: zero

  Scenario: eat cukes
    Given I have 10 cukes
    When I eat 0
    Then I have 10 left
