export DEBIAN_FRONTEND=noninteractive

# Updating system.
sudo apt-get update
# sudo apt-get upgrade -y

# Installing system tools.
sudo apt-get install -y vim tmux

# Installing mongoDB.
sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv 7F0CEB10
echo 'deb http://downloads-distro.mongodb.org/repo/ubuntu-upstart dist 10gen' | sudo tee /etc/apt/sources.list.d/mongodb.list
sudo apt-get update
sudo apt-get install -y mongodb-org
