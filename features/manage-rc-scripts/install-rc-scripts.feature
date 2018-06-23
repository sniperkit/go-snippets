Feature: Install rc scripts
  As a user
  In order to use the rc scripts in my repository
  I want to install the rc scripts on my machine

  Scenario: rc script does not exist
    Given no rc script exists
    When I run the command to install my rc scripts
    Then my rc scripts are created in their respective locations
