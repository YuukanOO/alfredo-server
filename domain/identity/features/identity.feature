Feature: Authentification
  In order to access the system
  As a user
  I need to be able to authenticate myself

  Scenario: Registering a new user
    Given I have the email "bob@roger.org"
    And the password "bobpwd"
    When I register a new account
    Then my account should have been created