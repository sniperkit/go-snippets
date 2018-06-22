package main

import (
	"encoding/json"
	"testing"
)

const invalidLua = `
func Foo() (int, error) {
	return 0, nil
}
`

// taken from http://lua-users.org/files/wiki_insecure/users/chill/table.binsearch-0.3.lua
const validLua = `
--[[
   table.bininsert( table, value [, comp] )
   
   Inserts a given value through BinaryInsert into the table sorted by [, comp].
   
   If 'comp' is given, then it must be a function that receives
   two table elements, and returns true when the first is less
   than the second, e.g. comp = function(a, b) return a > b end,
   will give a sorted table, with the biggest value on position 1.
   [, comp] behaves as in table.sort(table, value [, comp])
   returns the index where 'value' was inserted
]]--
do
   -- Avoid heap allocs for performance
   local fcomp_default = function( a,b ) return a < b end
   function table.bininsert(t, value, fcomp)
      -- Initialise compare function
      local fcomp = fcomp or fcomp_default
      --  Initialise numbers
      local iStart,iEnd,iMid,iState = 1,#t,1,0
      -- Get insert position
      while iStart <= iEnd do
         -- calculate middle
         iMid = math.floor( (iStart+iEnd)/2 )
         -- compare
         if fcomp( value,t[iMid] ) then
            iEnd,iState = iMid - 1,0
         else
            iStart,iState = iMid + 1,1
         end
      end
      table.insert( t,(iMid+iState),value )
      return (iMid+iState)
   end
end
-- CHILLCODE™
`

func TestProcessRequest(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		numErrors int
		status    Status
		astEmpty  bool
		marshal   bool
	}{
		{
			"valid input",
			validLua,
			0,
			Ok,
			false,
			true,
		},
		{
			"invalid input",
			invalidLua,
			1,
			Error,
			true,
			true,
		},
		{
			"not json",
			"fooo bar baz",
			1,
			Fatal,
			true,
			false,
		},
	}

	for _, c := range cases {
		var req []byte
		var err error
		if c.marshal {
			req, err = json.Marshal(&Request{
				Content: c.input,
			})
			if err != nil {
				t.Fatalf("%s: unexpected error: %s", c.name, err)
			}
		} else {
			req = []byte(c.input)
		}

		resp := processRequest(req)
		if resp.Status != c.status {
			t.Errorf("%s: expecting status %s, got %s", c.name, c.status, resp.Status)
		}

		if resp.AST == nil && !c.astEmpty {
			t.Errorf("%s: unexpected empty AST", c.name)
		}

		if resp.AST != nil && c.astEmpty {
			t.Errorf("%s: expected empty AST", c.name)
		}

		if len(resp.Errors) != c.numErrors {
			t.Errorf("%s: expecting %d errors, got %d", c.name, c.numErrors, len(resp.Errors))
		}
	}
}
