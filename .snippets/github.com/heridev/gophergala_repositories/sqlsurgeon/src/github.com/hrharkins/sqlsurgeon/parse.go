package sqlsurgeon;
import "strings";

type text string
type output chan AnyStatement
type parser func(output, text) (bool, text, AnyParseError)

//////////////////////////////////////////////////////////////////////////////

func (t text) eatSpace() text {
    return text(strings.TrimSpace(string(t)))
}

func (t text) atKeyword(k string) (bool, text) {
    if strings.HasPrefix(string(t), k) {
        return true, text(strings.TrimPrefix(string(t), k))
    } else {
        return false, t
    }
}

//////////////////////////////////////////////////////////////////////////////

func parsSelect(o output, t text) (bool, text, AnyParseError) {
    t = t.eatSpace()
    if is_kw, _ := t.atKeyword("select"); is_kw {
    }

    return false, t, nil
}

