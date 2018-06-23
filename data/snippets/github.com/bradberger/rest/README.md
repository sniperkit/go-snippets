A simple package for quick bootstrapping of REST API creation in the 
Google App Engine Go environment. It's currently used in production on 
several sites.

Documentation forthcoming.

### Testing handlers is easy ###

There are built-in test functions to make testing easy. For example:

```
import (
    "github.com/bradberger/context"
    "github.com/bradberger/rest"
)

type Foo struct {
    Bar string `json:"bar"`
}

func myHandler(ctx context.Context) error {
    return rest.JSON(ctx, http.StatusOK, &Foo{Bar: "foobar"})
}

func TestSomeHandler() {
    var f Foo
    r := TestPostJSON("/foo/bar", &f)
    
    assert.NoError(t, r.Do(myHandler))
    assert.NoError(t, r.Decode(&f))
    assert.Equal(t, http.StatusOK, r.Writer.Code)
    assert.Equal(t, "foobar", f.Bar)
}
```

If you need more control over the request, like custom headers, etc., then
you can update it after creation and before calling `Do()`:

```
func TestSomeHandler() {
    var f Foo

    r := TestPostJSON("/foo/bar", &f)
    r.Request.Header.Set("Authorization", "Bearer abc123")
    
    assert.NoError(t, r.Do(myHandler))
    assert.NoError(t, r.Decode(&f))
    assert.Equal(t, http.StatusOK, r.Writer.Code)
    assert.Equal(t, "foobar", f.Bar)
}
```