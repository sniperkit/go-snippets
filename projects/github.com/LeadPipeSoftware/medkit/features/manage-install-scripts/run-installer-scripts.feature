Feature: Run installer scripts
  As a user
  In order to have all of my favorite software installed
  I want to run the installer scripts I have written

  Scenario: Installer script has not been run
    Given the installer script has not been run before
     When I run the command to run my installer scripts
     Then the installer scripts are executed

