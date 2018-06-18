#!/usr/bin/env python

def getDecision(name, game, answer):
  print name, game, answer


my_game = 'my game'
my_answer = 'my answer'

getDecision('company 0', my_game, my_answer)


def make_decision(game, answer):
  def callDecision (name):
    getDecision(name, game, answer)
  return callDecision

caller = make_decision(my_game, my_answer)
caller('company 1')
caller('company 2')
