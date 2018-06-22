#!/bin/bash

echo -n "Please provide the e-mail of your github.com account [ENTER]:"
read github_mail

echo -n "Please provide your name [Enter]:"
read github_name

echo "Using github account: $github_name<$github_mail>"

git config --global user.email "$github_mail"
git config --global user.name "$github_name"
git config --global push.default simple
git config --global core.editor "gedit"

echo "generating a new SSH key that can be added to your github account"
echo "" | ssh-keygen -f $HOME/.ssh/id_rsa -t rsa -b 4096 -C "$github_mail"

# start the ssh-agent in the background
eval $(ssh-agent -s)

ssh-add $HOME/.ssh/id_rsa

echo "Add the new SSH key to your github account."
echo -n "1. Go to https://github.com/settings/keys . Press [ENTER] when done."
read step_done
echo -n "2. Click 'New SSH key' and set a title to your key (e.g. workshop-pc). Press Press [ENTER] when done."
read step_done
echo "3. Paste the following key"
echo "----------------------------------"
cat $HOME/.ssh/id_rsa.pub
echo "----------------------------------"
echo -n "Press [ENTER] when done."
read step_done
