package main

import (
	"database/sql"
	"fmt"
	"github.com/datatogether/core"
	"github.com/ipfs/go-datastore"
	"time"
)

// ValidArchivingUrl checks to see if this url pattern-matches the list of subprimers
// TODO - there are many ways to spoof this, replace with actual URL matching.
func ValidArchivingUrl(db *sql.DB, url string) error {
	var exists bool
	err := db.QueryRow("select exists(select 1 from subprimers where $1 ilike concat('%', url ,'%'))", url).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("Oops! Only urls contained in subprimers can be archived. cannot archive %s", url)
	}
	return nil
}

func (c *Client) ArchiveUrl(db *sql.DB, reqId, url string) {
	if err := ValidArchivingUrl(db, url); err != nil {
		log.Info(err.Error())
		c.SendResponse(&ClientResponse{
			Type:      "URL_ARCHIVE_ERROR",
			RequestId: reqId,
			Error:     err.Error(),
		})
		return
	}

	// TODO - plumb userId into this
	_, err := db.Exec("insert into archive_requests (created,url,user_id) values ($1, $2, $3)", time.Now().Round(time.Second).In(time.UTC), url, "")
	if err != nil {
		log.Info(err.Error())
		c.SendResponse(&ClientResponse{
			Type:      "URL_ARCHIVE_ERROR",
			RequestId: reqId,
			Error:     err.Error(),
		})
		return
	}

	log.Info("archiving %s", url)
	u := &core.Url{Url: url}
	if _, err := u.ParsedUrl(); err != nil {
		log.Info(err.Error())
		c.SendResponse(&ClientResponse{
			Type:      "URL_ARCHIVE_ERROR",
			RequestId: reqId,
			Error:     fmt.Sprintf("url parse error: %s", err.Error()),
		})
		return
	}

	if err := u.Read(store); err != nil {
		if err == core.ErrNotFound {
			if err := u.Save(store); err != nil {
				log.Info(err.Error())
				c.SendResponse(&ClientResponse{
					Type:      "URL_ARCHIVE_ERROR",
					RequestId: reqId,
					Error:     fmt.Sprintf("internal server error"),
				})
				return
			}
		} else {
			log.Info(err.Error())
			c.SendResponse(&ClientResponse{
				Type:      "URL_ARCHIVE_ERROR",
				RequestId: reqId,
				Error:     fmt.Sprintf("internal server error"),
			})
			return
		}
	}

	// Initial get succeeded, let the client know
	c.SendResponse(&ClientResponse{
		Type:      "URL_ARCHIVE_SUCCESS",
		RequestId: reqId,
		Schema:    "URL",
		Data:      u,
	})

	// Perform base GET request
	_, links, err := u.Get(store)
	if err != nil {
		log.Info(err.Error())
		c.SendResponse(&ClientResponse{
			Type:      "URL_ARCHIVE_ERROR",
			RequestId: reqId,
			Error:     err.Error(),
		})
		return
	}

	// push our new links to client
	c.SendResponse(&ClientResponse{
		Type:      FetchOutboundLinksAct{}.SuccessType(),
		RequestId: "server",
		Schema:    "LINK_ARRAY",
		Data:      links,
	})

	go func(db *sql.DB, links []*core.Link) {
		// GET each destination link from this page in parallel
		for _, l := range links {
			// need a sleep here to avoid bombing server with requests
			// tooooo hard, also we sleep first b/c the websocket trips up if
			// we jam the messages to hard.
			time.Sleep(time.Second * 3)

			c.SendResponse(&ClientResponse{
				Type:      "URL_SET_LOADING",
				RequestId: "server",
				Data: map[string]interface{}{
					"url":     l.Dst.Url,
					"loading": true,
				},
			})

			if _, _, err := l.Dst.Get(store); err != nil {
				log.Info(err.Error())
				c.SendResponse(&ClientResponse{
					Type:      "URL_SET_ERROR",
					RequestId: "server",
					Data: map[string]interface{}{
						"url":   l.Dst.Url,
						"error": err.Error(),
					},
				})
			}

			c.SendResponse(&ClientResponse{
				Type:      "URL_SET_SUCCESS",
				RequestId: "server",
				Data: map[string]interface{}{
					"url":     l.Dst.Url,
					"success": true,
				},
			})
		}
	}(db, links)
}

// ArchiveUrl GET's a url and if it's an HTML page, any links it directly references
func ArchiveUrl(db *sql.DB, url string, done func(err error)) (*core.Url, []*core.Link, error) {
	u := &core.Url{Url: url}
	if _, err := u.ParsedUrl(); err != nil {
		done(err)
		return nil, nil, err
	}

	if err := u.Read(store); err != nil {
		if err == core.ErrNotFound {
			if err := u.Save(store); err != nil {
				done(err)
				return nil, nil, err
			}
		} else {
			done(err)
			return nil, nil, err
		}
	}

	// Perform GET request
	_, links, err := u.Get(store)
	if err != nil {
		done(err)
		return u, links, err
	}

	tasks := len(links)
	errs := make(chan error, tasks)

	go func(store datastore.Datastore, links []*core.Link) {
		// GET each destination link from this page in parallel
		for _, l := range links {
			if _, _, err := l.Dst.Get(store); err != nil {
				log.Info(err.Error())
			}
			errs <- nil

			// need a sleep here to avoid bombing server with requests
			// tooooo hard
			time.Sleep(time.Second * 3)
		}
	}(store, links)

	go func() {
		for i := 0; i < tasks; i++ {
			err := <-errs
			if err != nil {
				done(err)
				return
			}
		}
		done(nil)
	}()

	return u, links, err
}

func ArchiveUrlSync(db *sql.DB, url string) (*core.Url, error) {
	done := make(chan error)
	u, _, err := ArchiveUrl(db, url, func(err error) {
		done <- err
	})
	if err != nil {
		return u, err
	}

	err = <-done
	return u, err
}
