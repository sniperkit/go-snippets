-- === iTerm2 Bindings ===

local log = hs.logger.new('iterm2', 'debug')

-- Alternative iTerm2 hotkey
local function iTermHotkeyHandler()
  local iterm = hs.application.get('iTerm2')
  if (iterm) then
    local window = iterm:allWindows()
    if iterm:isFrontmost() then
      log:d('Gonna hide')
      iterm:hide()
    elseif #window > 0 then
      log:d('Mmm, iter windows')
      for _, w in pairs(window) do
        w:focus()
      end
    else
      log:d('Gonna create a window')
      iterm:activate()
      iterm:selectMenuItem('New Window')
    end
  end
end
-- hs.hotkey.bind('cmd', 'space', iTermHotkeyHandler)
