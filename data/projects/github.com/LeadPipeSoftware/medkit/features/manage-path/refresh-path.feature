Feature: Refresh PATH
  As a user
  In order to make sure my PATH is up-to-date
  I want to refresh my PATH using my configuration

  Scenario: PATH files exist
    Given at least one PATH file exists
     When run the command to refresh my PATH variable
     Then the PATH variable is updated

