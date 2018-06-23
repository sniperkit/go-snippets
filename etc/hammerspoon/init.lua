require 'hs.ipc'
require 'libzen'
require 'ddcci'
require 'consts'

local log = hs.logger.new('riptide', 'debug')

-- Bring all Finder windows forward when one gets activated
function finderActivationHandler(appName, eventType, appObject)
  if (eventType == hs.application.watcher.activated) then
    if (appName == "Finder") then
      appObject:selectMenuItem({"Window", "Bring All to Front"})
    end
  end
end
appWatcher = hs.application.watcher.new(finderActivationHandler)
appWatcher:start()

-- Misc actions
function startScreenSaver()
  hs.caffeinate.startScreensaver()
end

function lockScreen()
  hs.caffeinate.lockScreen()
end

-- Add menubar controls
local actionIcon = getIcon('call_to_action', 16)
bar = hs.menubar.new()
  :setIcon(actionIcon)
  :setMenu({
    { title = '-' },
    { title = 'Secondary Monitor Brightness', disabled = true },
    { title = 'Sync', fn = syncMonitorBrightness, image = getIcon('broken_image', 16) },
    { title = 'Turn Down', fn = setMonitorBrightnessLow, image = getIcon('brightness_low', 16) },
    { title = 'Okay-Okay', fn = setMonitorBrightnessMedium, image = getIcon('brightness_medium', 16) },
    { title = 'Crank Up', fn = setMonitorBrightnessHigh, image = getIcon('brightness_high', 16) },
    { title = '-' },
    { title = 'Start Screensaver', fn = startScreenSaver, image = getIcon('lightbulb_outline', 16) },
    { title = 'Lock', fn = lockScreen, image = getIcon('lock', 16) },
  })

function logev(event)
  local keyProps = event:systemKey()
  if keyProps.down then
    if keyProps.key == 'SOUND_UP' then
      volumeKeyUp()
    elseif keyProps.key == 'SOUND_DOWN' then
      volumeKeyDown()
    end
  end
end

tap = hs.eventtap.new({ hs.eventtap.event.types.NSSystemDefined }, logev)
  :start()

-- Watch screen for changes
function onScreenLayoutChange()
  if hs.screen.allScreens()[2] ~= nil then
    hs.settings.set(HS_EXT_DISP, true)
  else
    hs.settings.set(HS_EXT_DISP, false)
  end
end

screentap = hs.screen.watcher.new(onScreenLayoutChange)
  :start()
onScreenLayoutChange()
