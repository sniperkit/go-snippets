Feature: Install dotfiles
  As a user
  In order to use the dotfiles in my repository
  I want to install the dotfiles on my machine

  Scenario: Dotfile does not exist
    Given no dotfile exists
     When I run the command to install my dotfiles
     Then my dotfiles are created in their respective locations

