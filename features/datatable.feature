Feature: data tables

  Scenario: header tables
    Given this data table
    | x | y | z | x + y + z |
    | 1 | 2 | 3 | 6  |
    | 2 | 4 | 6 | 12 |
    | 3 | 5 | 1 | 9  |
    Then it has 4 rows and 4 columns
    And 3 rows as a header table
    Then the values add up when accessed as a header table
    And the total sum of the x + y + z column is 27