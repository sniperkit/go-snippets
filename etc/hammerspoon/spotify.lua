require 'libzen'
local log = hs.logger.new('riptide.spotify', 'debug')

-- Icons
local playIcon = getIcon('media-play-outline', 16)
local pauseIcon = getIcon('media-pause', 16)

-- Mute spotify if there's an Ad. Yeah...
local _spotifyWasMuted
local function muteSpotifyOnAd()
  if hs.spotify.isPlaying() then
    if hs.spotify.getCurrentTrack():lower():has('spotify') then
      _spotifyWasMuted = true
      hs.audiodevice.defaultOutputDevice():setMuted(true)
    elseif _spotifyWasMuted then
      _spotifyWasMuted = false
      hs.audiodevice.defaultOutputDevice():setMuted(false)
    end
  end
end
spotify = hs.timer.new(5, muteSpotifyOnAd, true)
spotify:start()

local function playerControls()
  if hs.spotify.isPlaying() then
    hs.spotify.pause()
    spotifyMenuIcon:setIcon(playIcon)
  else
    hs.spotify.play()
    spotifyMenuIcon:setIcon(pauseIcon)
  end
end

-- spotifyMenuIcon = hs.menubar.new()
--   :setIcon(playIcon)
--   :setClickCallback(playerControls)
