Feature: DataTables

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