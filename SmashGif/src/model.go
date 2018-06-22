package hello

type Content struct {
	Comments  string
	Upvotes   int
	Subreddit string
}

// type Content interface {
// 	Prepare()
// }

// Gfycat content
type Gif struct {
	Content   Content
	GameTitle string
	GifTitle  string
	GifId     string
}

// Youtube content
type Youtube struct {
	Content
}
