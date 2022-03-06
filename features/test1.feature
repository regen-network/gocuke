Feature: DataTables
  Rule: test 1
    Scenario: simple
      Given I have a data table
      | 1 |
      | 3 |
      | 7 |
      When some doc string:
      """markdown
      # Hello
      """
      And add 5
      ```go
      x := 5
      ```
      Then pass

    Scenario Outline: eating
      Given there are <start> cucumbers
      When I eat <eat> cucumbers
      Then I should have <left> cucumbers

      Examples:
        | start | eat | left |
        |    12 |   5 |    7 |
        |    20 |   5 |   15 |

