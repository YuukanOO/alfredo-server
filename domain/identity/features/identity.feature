Feature: Authentification
  In order to access the system
  As a user
  I need to be able to register and authenticate myself

  Scenario: Registering a new user
    Given I have the email "bob@roger.org"
    And the password "bobpwd"
    And the email is not already registered
    When I register a new account
    Then my account should have been created

  Scenario: Registering a user with an email already taken
    Given I have the email "already@taken.org"
    And the email is already registered
    When I register a new account
    Then it should complains about the email being already taken

  Scenario: Authenticating an invalid user
    Given I have the email "anything@taken.org"
    And the email is not already registered
    When I authenticate
    Then it should complains about the email not being registered

  Scenario: Authenticating a valid user with a wrong password
    Given I have the email "already@taken.org"
    And the password "alreadytakenwrongpassword"
    And the valid password is "alreadytak3n"
    When I authenticate
    Then it should complains about the password being incorrect

  Scenario: Authenticating valid credentials
    Given I have the email "already@taken.org"
    And the password "alreadytak3n"
    And the valid password is "alreadytak3n"
    When I authenticate
    Then I should receive an access token