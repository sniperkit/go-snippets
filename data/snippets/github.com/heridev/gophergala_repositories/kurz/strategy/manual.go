package strategy

import (
	"errors"
	"fmt"
	"github.com/FGM/kurz/storage"
	"github.com/FGM/kurz/url"
	"time"
)

type ManualStrategy struct {
	base baseStrategy
}

func (y ManualStrategy) Name() string {
	return "manual"
}

/*
Alias() implements AliasingStrategy.Alias().

  - options is expected to be a non empty single string
*/
func (y ManualStrategy) Alias(long *url.LongUrl, s storage.Storage, options ...interface{}) (url.ShortUrl, error) {
	var ret url.ShortUrl
	var err error
	if len(options) != 1 {
		err = errors.New("ManualString.Alias() takes a single optional parameter, which is the requested short URL.")
		return ret, err
	} else {
		requestedAlias, ok := options[0].(string)
		if !ok {
			err = errors.New(fmt.Sprintf("ManualString.Alias() optional parameter must be a string, got: %+v", requestedAlias))
			return ret, err
		}

		err = nil
		/** TODO
		 * - validate alias is available
		 */
		ret = url.ShortUrl{
			Value:       requestedAlias,
			ShortFor:    *long,
			Domain:      long.Domain(),
			Strategy:    y.Name(),
			SubmittedBy: storage.CurrentUser(),
			SubmittedOn: time.Now().UTC().Unix(),
			IsEnabled:   true,
		}
	}
	return ret, err
}

func (y ManualStrategy) UseCount(s storage.Storage) int {
	return y.base.UseCount(s)
}
