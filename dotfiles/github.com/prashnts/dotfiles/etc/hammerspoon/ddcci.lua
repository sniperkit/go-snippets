require 'consts'

local log = hs.logger.new('riptide.ddcci', 'debug')

local ddcctl = '/Users/prashantsinha/opt/bin/ddcctl'
local CMDFMT = '%s -d 1 -%s %s'

-- Use volume keys when external monitor is connected with ddc
if hs.settings.get(HS_VOLUME) == nil then
  hs.settings.set(HS_VOLUME, 8)
end

function syncVolume()
  if hs.settings.get(HS_EXT_DISP) == true then
    local volume = hs.settings.get(HS_VOLUME)
    local cmd = CMDFMT:format(ddcctl, 'v', volume)
    log:d('sync volume with', cmd)
    hs.execute(cmd)
  else
    log:d('no external display, wont sync volume')
  end
end

function volumeKeyUp()
  local basevol = hs.settings.get(HS_VOLUME)
  if basevol < 254 then
    hs.settings.set(HS_VOLUME, basevol + 1)
  end
  syncVolume()
end

function volumeKeyDown()
  local basevol = hs.settings.get(HS_VOLUME)
  if basevol > 1 then
    hs.settings.set(HS_VOLUME, basevol - 1)
  end
  syncVolume()
end


-- Brightness settings
function setMonitorBrightness(val)
  local cmd = CMDFMT:format(ddcctl, 'b', val)
  log:d('sync brightness with', cmd)
  hs.execute(cmd)
end

-- Sync main display brightness with external monitor
function syncMonitorBrightness()
  local mainDisplayBrightness = hs.brightness.get()
  setMonitorBrightness(mainDisplayBrightness)
end

function setMonitorBrightnessLow()
  setMonitorBrightness(0)
end

function setMonitorBrightnessMedium()
  setMonitorBrightness(50)
end

function setMonitorBrightnessHigh()
  setMonitorBrightness(100)
end
