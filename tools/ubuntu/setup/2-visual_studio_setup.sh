#!/bin/bash

echo "Downloading Visual Studio Code"

printf "\n*************************************\n"
printf "* 1. Downloading Visual Studio Code *\n"
printf "*************************************\n\n"

wget https://vscode-update.azurewebsites.net/latest/linux-deb-x64/stable -O $HOME/Downloads/visual_studio.deb


printf "\n*************************************\n"
printf "* 2. Installing Visual Studio Code  *\n"
printf "*************************************\n\n"
sudo apt install $HOME/Downloads/visual_studio.deb

printf "\n************************************************\n"
printf "* 3. Deleting the Visual Studio Code installer *\n"
printf "************************************************\n\n"
rm $HOME/Downloads/visual_studio.deb
