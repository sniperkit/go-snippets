-- === Webview Popouts ===
--
require 'libzen'

local wktoolbar = require 'hs.webview.toolbar'
local log = hs.logger.new('riptide.popouts', 'debug')

toolbar = wktoolbar.new('poptools', {
  { id = 'navGroup', label = 'Navigation', groupMembers = { 'navLeft', 'navRight' }},
  { id = 'navLeft', image = getIcon('arrow-left'), allowedAlone = false},
  { id = 'navRight', image = getIcon('arrow-right'), allowedAlone = false},
  { id = 'reload', image = getIcon('reload')},
})
toolbar:displayMode('icon')
toolbar:sizeMode('small')

local rect = hs.geometry.rect(900, 20, 520, 550)
ww = hs.webview.newBrowser(rect)
ww:windowStyle(
  hs.webview.windowMasks['titled'] |
  hs.webview.windowMasks['resizable'] |
  hs.webview.windowMasks['closable'] |
  hs.webview.windowMasks['utility'] |
  hs.webview.windowMasks['HUD']
)
ww:userAgent('Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2725.0 Safari/537.36')
ww:allowNewWindows(false)
-- ww:toolbar(toolbar)
ww:url('https://messenger.com')
ww:level(1000)

local hidden = true

function webViewClickHandler()
  if hidden then
    ww:show(.2)
  else
    ww:hide(.2)
  end
  hidden = not hidden
end

popoutMenuIcon = hs.menubar.new()
  :setIcon(getIcon('facebook', 16))
  :setClickCallback(webViewClickHandler)
