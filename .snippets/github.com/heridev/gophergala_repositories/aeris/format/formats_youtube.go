package format

type YoutubeFormat struct {
	Format
	Itag int
}

func YoutubeFormats() map[string]*YoutubeFormat {
	return youtubeFormats
}

// https://en.wikipedia.org/wiki/YouTube#Quality_and_formats
var youtubeFormats = map[string]*YoutubeFormat{
	"5": &YoutubeFormat{
		Itag: 5,
		Format: Format{
			Container: "FLV",
			Video: VideoFormat{
				Resolution: "240p",
				Encoding:   "Sorenson H.263",
			},
			Audio: AudioFormat{
				Encoding: "MP3",
				Bitrate:  64,
			},
		},
	},
	"17": &YoutubeFormat{
		Itag: 17,
		Format: Format{
			Container: "3GP",
			Video: VideoFormat{
				Resolution: "144p",
				Encoding:   "MPEG-4 Visual",
			},
			Audio: AudioFormat{
				Encoding: "AAC",
				Bitrate:  24,
			},
		},
	},
	"18": &YoutubeFormat{
		Itag: 18,
		Format: Format{
			Container: "MP4",
			Video: VideoFormat{
				Resolution: "360p",
				Encoding:   "H.264",
			},
			Audio: AudioFormat{
				Encoding: "AAC",
				Bitrate:  96,
			},
		},
	},
	"22": &YoutubeFormat{
		Itag: 22,
		Format: Format{
			Container: "MP4",
			Video: VideoFormat{
				Resolution: "720p",
				Encoding:   "H.264",
			},
			Audio: AudioFormat{
				Encoding: "AAC",
				Bitrate:  192,
			},
		},
	},
	"36": &YoutubeFormat{
		Itag: 36,
		Format: Format{
			Container: "3GP",
			Video: VideoFormat{
				Resolution: "240p",
				Encoding:   "MPEG-4 Visual",
			},
			Audio: AudioFormat{
				Encoding: "AAC",
				Bitrate:  36,
			},
		},
	},
	"43": &YoutubeFormat{
		Itag: 43,
		Format: Format{
			Container: "WebM",
			Video: VideoFormat{
				Resolution: "360p",
				Encoding:   "VP8",
			},
			Audio: AudioFormat{
				Encoding: "Vorbis",
				Bitrate:  128,
			},
		},
	},
	"82": &YoutubeFormat{
		Itag: 82,
		Format: Format{
			Container: "MP4",
			Video: VideoFormat{
				Resolution: "360p",
				Encoding:   "H.264",
			},
			Audio: AudioFormat{
				Encoding: "AAC",
				Bitrate:  96,
			},
		},
	},
	"83": &YoutubeFormat{
		Itag: 83,
		Format: Format{
			Container: "MP4",
			Video: VideoFormat{
				Resolution: "240p",
				Encoding:   "H.264",
			},
			Audio: AudioFormat{
				Encoding: "AAC",
				Bitrate:  96,
			},
		},
	},
	"84": &YoutubeFormat{
		Itag: 84,
		Format: Format{
			Container: "MP4",
			Video: VideoFormat{
				Resolution: "720p",
				Encoding:   "H.264",
			},
			Audio: AudioFormat{
				Encoding: "AAC",
				Bitrate:  192,
			},
		},
	},
	"85": &YoutubeFormat{
		Itag: 85,
		Format: Format{
			Container: "MP4",
			Video: VideoFormat{
				Resolution: "1080p",
				Encoding:   "H.264",
			},
			Audio: AudioFormat{
				Encoding: "AAC",
				Bitrate:  192,
			},
		},
	},
	"100": &YoutubeFormat{
		Itag: 100,
		Format: Format{
			Container: "WebM",
			Video: VideoFormat{
				Resolution: "360p",
				Encoding:   "VP8",
			},
			Audio: AudioFormat{
				Encoding: "Vorbis",
				Bitrate:  128,
			},
		},
	},
}
