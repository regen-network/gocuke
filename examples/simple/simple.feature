Feature: simple

  Scenario Outline: eat cukes
    Given I have <x> cukes
    When I eat <y>
    Then I have <z> left

    Examples:
    | x | y | z |
    | 5 | 3 | 2 |
    | 10 | 2 | 8 |
