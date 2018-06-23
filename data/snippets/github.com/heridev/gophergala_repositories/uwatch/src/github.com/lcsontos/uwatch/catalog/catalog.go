package catalog

import (
	"fmt"
	"time"
)

type VideoRecord struct {
	Id          int64
	Description string
	PublishedAt time.Time
	VideoId     string
	Title       string
}

type NoSuchVideoError struct {
	VideoId string
}

type VideoCatalog interface {
	SearchByID(videoId string) (*VideoRecord, error)
	SearchByTitle(title string, maxResults int64) ([]VideoRecord, error)
}

func (err *NoSuchVideoError) Error() string {
	return fmt.Sprintf("Video with Id %s does not exist", err.VideoId)
}

func NewVideoRecord(videoId, title, description string, publishedAt time.Time) *VideoRecord {
	return &VideoRecord{Description: description, PublishedAt: publishedAt, VideoId: videoId, Title: title}
}

func (videoRecord *VideoRecord) String() string {
	return fmt.Sprintf("[%v] %v: %v", videoRecord.VideoId, videoRecord.Title, videoRecord.Description)
}
