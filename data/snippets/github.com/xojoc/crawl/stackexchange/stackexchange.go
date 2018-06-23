/*  Copyright (C) 2018 Alexandru Cojocaru

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>. */

package stackexchange

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"time"

	"xojoc.pw/must"
	"xojoc.pw/u"
)

type Post struct {
	ID int // `json:"-"`
	// 1 - question
	// 2 - answer
	PostTypeID int
	ParentID   int
	// 0 if no accepted anwser
	AcceptedAnswerID int
	Score            int
	ViewCount        int
	Body             []byte
	Title            []byte
	Tags             [][]byte
	AnswerCount      int
	CommentCount     int
	FavoriteCount    int

	CreationDate time.Time
	LastEditDate time.Time
}

var (
	rowStart = []byte(`  <row`)

	idLabel               = []byte(` Id=`)
	postTypeIDLabel       = []byte(` PostTypeId=`)
	parentIDLabel         = []byte(` ParentId=`)
	acceptedAnswerIDLabel = []byte(` AcceptedAnswerId=`)
	scoreLabel            = []byte(` Score=`)
	viewCountLabel        = []byte(` ViewCount=`)
	bodyLabel             = []byte(` Body=`)
	titleLabel            = []byte(` Title=`)
	tagsLabel             = []byte(` Tags=`)
	answerCountLabel      = []byte(` AnswerCount=`)
	commentCountLabel     = []byte(` CommentCount=`)
	favoriteCountLabel    = []byte(` FavoriteCount=`)

	creationDateLabel = []byte(` CreationDate=`)
	lastEditDateLabel = []byte(` LastEditDate=`)
)

func LoadDB(r io.Reader) map[int]*Post {
	dec := json.NewDecoder(r)
	var posts map[int]*Post
	must.OK(dec.Decode(&posts))
	return posts
}

func StoreDB(w io.Writer, posts interface{}) error {
	// func StoreDB(w io.Writer, posts map[int]*Post) error {
	enc := json.NewEncoder(w)
	return enc.Encode(posts)
}

var (
	intSyntaxError = errors.New("syntax error")
)

func toInt(bs []byte) (int, error) {
	if len(bs) == 0 {
		return 0, intSyntaxError
	}
	bs0 := bs
	if bs[0] == '-' {
		bs = bs[1:]
	}

	n := 0
	for _, b := range bs {
		b -= '0'
		if b > 9 {
			return 0, intSyntaxError
		}
		n = n*10 + int(b)
	}
	if bs0[0] == '-' {
		n = -n
	}
	return n, nil
}

func clone(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

type Filter struct {
	PostTypeID        int
	HasAcceptedAnswer bool
	MinScore          int
	MinAnwserCount    int
}

func NewFilter() Filter {
	return Filter{MinScore: int(^uint(0) >> 1)}
}

func Posts(r io.Reader, f Filter) (chan *Post, chan error) {
	posts := make(chan *Post, 1000)
	cerr := make(chan error, 1)

	go func() {
		buf := bufio.NewReaderSize(r, 100*1024*1024)
		buf.ReadSlice('\n')
		buf.ReadSlice('\n')
		for {
		next:
			var err error
			var id, postTypeID, parentID, acceptedAnswerID, score, viewCount int
			var body, title []byte
			var tags [][]byte
			var answerCount, commentCount, favoriteCount int
			var creationDate, lastEditDate time.Time

			_ = parentID
			_ = commentCount
			_ = creationDate
			_ = lastEditDate
			_ = favoriteCount

			l, err := buf.ReadSlice('\n')
			if err == io.EOF {
				break
			}
			must.OK(err)

			l = l[len(rowStart):]

			for {
				i := bytes.IndexByte(l, '"')
				if i == -1 {
					break
				}
				s1 := l[:i]
				l = l[i+1:]
				i = bytes.IndexByte(l, '"')
				s2 := l[:i]
				l = l[i+1:]

				switch {
				case bytes.Equal(s1, idLabel):
					id, err = toInt(s2)
					if err != nil {
						cerr <- err
						goto stop
					}
				case bytes.Equal(s1, postTypeIDLabel):
					postTypeID, err = toInt(s2)
					if err != nil {
						cerr <- err
						goto stop
					}
					if f.PostTypeID != 0 {
						if postTypeID != f.PostTypeID {
							goto next
						}
					}
					/*
						case bytes.Equal(s1, parentIDLabel):
							id, err := strconv.Atoi(string(s2))
							must.OK(err)
							p.ParentID = id
					*/
				case bytes.Equal(s1, acceptedAnswerIDLabel):
					acceptedAnswerID, err = toInt(s2)
					if err != nil {
						cerr <- err
						goto stop
					}
				case bytes.Equal(s1, scoreLabel):
					score, err = toInt(s2)
					if err != nil {
						cerr <- err
						goto stop
					}
					if score < f.MinScore {
						goto next
					}
				case bytes.Equal(s1, viewCountLabel):
					viewCount, err = toInt(s2)

					if err != nil {
						cerr <- err
						goto stop
					}
				case bytes.Equal(s1, tagsLabel):
					ts := bytes.Split(s2, []byte("&gt;&lt;"))
					ts[0] = ts[0][4:]
					tl := ts[len(ts)-1]
					ts[len(ts)-1] = tl[:len(tl)-4]
					for _, t := range ts {
						tags = append(tags, u.Clone(t))
					}
				case bytes.Equal(s1, answerCountLabel):
					answerCount, err = toInt(s2)

					if err != nil {
						cerr <- err
						goto stop
					}
				case bytes.Equal(s1, bodyLabel):
					//		p.Body = html.UnescapeString(string(s2))
					body = clone(s2)
				case bytes.Equal(s1, titleLabel):
					//				p.Title = html.UnescapeString(string(s2))
					title = clone(s2)
				}
			}

			if f.HasAcceptedAnswer {
				if acceptedAnswerID == 0 {
					goto next
				}
			}

			p := Post{}
			p.ID = id
			p.PostTypeID = postTypeID
			p.AcceptedAnswerID = acceptedAnswerID
			p.Score = score
			p.ViewCount = viewCount
			p.Body = body
			p.Title = title
			p.Tags = tags
			p.AnswerCount = answerCount

			posts <- &p
		}
	stop:
		close(posts)
		close(cerr)
	}()

	return posts, cerr
}

//
// func LoadPosts(r io.Reader) map[int]*Post {
// 	buf := bufio.NewReaderSize(r, 500*1024*1024)
// 	buf.ReadSlice('\n')
// 	buf.ReadSlice('\n')
//
// 	posts := map[int]*Post{}
//
// 	for {
// 	next:
// 		var err error
// 		var id, postTypeID, parentID, acceptedAnswerID, score, viewCount int
// 		var body, title []byte
// 		var tags [][]byte
// 		var answerCount, commentCount, favoriteCount int
// 		var creationDate, lastEditDate time.Time
//
// 		l, err := buf.ReadSlice('\n')
// 		if err == io.EOF {
// 			break
// 		}
// 		must.OK(err)
//
// 		l = l[len(rowStart):]
//
// 		for {
// 			i := bytes.IndexByte(l, '"')
// 			if i == -1 {
// 				break
// 			}
// 			s1 := l[:i]
// 			l = l[i+1:]
// 			i = bytes.IndexByte(l, '"')
// 			s2 := l[:i]
// 			l = l[i+1:]
//
// 			switch {
// 			case bytes.Equal(s1, idLabel):
// 				id, err = toInt(s2)
// 				must.OK(err)
// 			case bytes.Equal(s1, postTypeIDLabel):
// 				postTypeID, err = toInt(s2)
// 				must.OK(err)
// 				if postTypeID != 1 {
// 					goto next
// 				}
// 				/*
// 					case bytes.Equal(s1, parentIDLabel):
// 						id, err := strconv.Atoi(string(s2))
// 						must.OK(err)
// 						p.ParentID = id
// 				*/
// case bytes.Equal(s1, acceptedAnswerIDLabel):
// 				acceptedAnswerID, err = toInt(s2)
// 				must.OK(err)
// 			case bytes.Equal(s1, scoreLabel):
// 				score, err = toInt(s2)
// 				must.OK(err)
// 				/*
// 					if score < 9 {
// 						goto next
// 					}
// 				*/
// 			case bytes.Equal(s1, viewCountLabel):
// 				viewCount, err = toInt(s2)
// 				must.OK(err)
// 			case bytes.Equal(s1, tagsLabel):
// 				ts := bytes.Split(s2, []byte("&gt;&lt;"))
// 				ts[0] = ts[0][4:]
// 				tl := ts[len(ts)-1]
// 				ts[len(ts)-1] = tl[:len(tl)-4]
// 				for _, t := range ts {
// 					tags = append(tags, clone(t))
// 				}
// 			case bytes.Equal(s1, answerCountLabel):
// 				answerCount, err = toInt(s2)
// 				must.OK(err)
// 			case bytes.Equal(s1, bodyLabel):
// 				//		p.Body = html.UnescapeString(string(s2))
// 				body = clone(s2)
// 			case bytes.Equal(s1, titleLabel):
// 				//				p.Title = html.UnescapeString(string(s2))
// 				title = clone(s2)
// 			}
// 		}
//
// 		/*
// 			if acceptedAnswerID == 0 {
// 				goto next
// 			}
// 		*/
//
// 		p := Post{}
// 		p.ID = id
// 		p.PostTypeID = postTypeID
// 		p.AcceptedAnswerID = acceptedAnswerID
// 		p.Score = score
// 		p.ViewCount = viewCount
// 		p.Body = body
// 		p.Title = title
// 		p.Tags = tags
// 		p.AnswerCount = answerCount
//
// 		posts[p.ID] = &p
//
// 		if len(posts)%(10*1000) == 0 {
// 			runtime.GC()
// 		}
// 	}
//
// 	return posts
// }
//
