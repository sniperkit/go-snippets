# Script to build and launch the app on an android device.

set -e

./make.bash

adb install -r bin/Robostats.apk

adb shell am start -a android.intent.action.MAIN \
	-n com.remote.robostats/com.remote.robostats.MainActivity
