Feature: tags

  Scenario Outline: tag expression 1
    Given the tag expression
    """
    @def and (@xyz or @1!:) and not (@qzy or @abc)
    """
    When I match "<tags>"
    Then the result is "<result>"
    Examples:
      | tags      | result |
      | @abc      | false  |
      | @def      | false  |
      | @def @xyz | true   |
      | @def @1!: | true   |
      | @qzy      | false  |


  @foo @bar[baz=bam]
  Scenario: some tags
    Given I eat some cukes

  @123abc
  Scenario: other tags
    Given I eat other cukes


