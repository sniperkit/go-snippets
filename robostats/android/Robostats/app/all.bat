@echo off

setlocal

echo # building robostats
call make.bat

echo # installing bin/Robostats.apk
adb install -r bin/Robostats.apk >nul

echo # starting com.robostats.remote.MainActivity
adb shell am start -a android.intent.action.MAIN -n com.robostats.remote/com.robostats.remote.MainActivity >nul
