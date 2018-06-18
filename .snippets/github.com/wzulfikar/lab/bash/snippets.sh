# move user's folder to external directory (OSX).
NEW_DIR=/Volumes/data
ln -s $NEW_DIR/Downloads ~
ln -s $NEW_DIR/Pictures ~
ln -s $NEW_DIR/Music ~
ln -s $NEW_DIR/Archive/dump_desktop ~/Desktop
