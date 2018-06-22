yps - Youtube Playlist Synchronizer
===

yps helps you get your Youtube playlist saved as mp3s so you can listen your favorite video lectures while on the move.

A good example for this would be: [Computer Science 61B - Fall 2006](https://www.youtube.com/playlist?list=PL4BBB74C7D2A1049C) or
[Google NYC Tech Talks](https://www.youtube.com/playlist?list=PLAD8A7B6D66DDD297)

[![GoDoc](https://godoc.org/github.com/gophergala/yps?status.svg)](https://godoc.org/github.com/gophergala/yps)

You can test it out [here](http://zeta-rush-835.appspot.com/)

TODO
---
In priority order:

- [x] Make a pretty interface
- [x] Implement an interface to receive a single video ID
- [x] Implement an interface to receive the playlist ID
- [ ] Process playlist request and divide items into tasks (using a MQ) to process files independently
- [ ] Implement a way to get the video in mp4 format from YouTube
    - [ ] If not possible, check what other services are doing
- [ ] Convert the mp4 into a mp3 file
- [ ] Save mp3 file into user destination
    - [ ] Goole Play Music (no apparent API, offer option to save mp3 file on local computer?)
    - [ ] Google Drive (check if they have an API)
    - [ ] Dropbox (maybe? check if they have an API)
    - [ ] local drive
- [ ] Implement a way to get the private videos / playlists of users

Technologies used
---
***TODO*** Check if this is what we need or we need more

- [x] AppEngine
    - [x] Task Queue
    - [ ] Managed VM (for the conversion part. maybe we can run a CLI to do it on the user machine?)
- [x] Go
    - [x] GorrilaWebToolkit - [http://www.gorillatoolkit.org/](http://www.gorillatoolkit.org/)
- Google APIs
    - [ ] Drive API (?)
    - [ ] Drive SDK (?, check difference between the two)
    - [ ] YouTube Data API v3 (to fetcht the playlists and compute the videos needed)
- FFmpeg (conversion stuff)

Resources
---
This project uses resurces from [initializr.com](http://www.initializr.com/)

Legal stuff
---
This should be a way to save material for later listening, while on the move. Please respect the copyright of the authors of the respective videos.

License
---
This project is created under MIT license.

For more information, please see the [LICENSE](https://github.com/gophergala/yps/blob/master/LICENSE.md) file.
