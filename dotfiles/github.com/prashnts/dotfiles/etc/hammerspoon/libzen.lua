-- === "standard" library ===
local log = hs.logger.new('riptide.zen', 'debug')

-- Check if string contains some other string
-- Usage: "it's all in vain":has("vain") == true
function string:has(pattern)
  -- Don't laugh. Loaded lua might not have a toboolean method.
  if self:match(pattern) then return true else return false end
end

-- Load an icon from ./icons dir.
function getIcon(name, size)
  img = hs.image.imageFromPath('./icons/' .. name .. '.png')
  if size == nil then
    return img
  else
    log:d('size')
    return img:setSize({h = size, w = size})
  end
end
