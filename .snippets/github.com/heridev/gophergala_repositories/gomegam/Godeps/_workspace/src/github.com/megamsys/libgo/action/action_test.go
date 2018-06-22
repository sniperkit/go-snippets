/*
** Copyright [2012-2014] [Megam Systems]
**
** Licensed under the Apache License, Version 2.0 (the "License");
** you may not use this file except in compliance with the License.
** You may obtain a copy of the License at
**
** http://www.apache.org/licenses/LICENSE-2.0
**
** Unless required by applicable law or agreed to in writing, software
** distributed under the License is distributed on an "AS IS" BASIS,
** WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
** See the License for the specific language governing permissions and
** limitations under the License.
*/

package action

import (
	"gopkg.in/check.v1"
	"testing"
)

func Test(t *testing.T) {
	check.TestingT(t)
}

type S struct{}

var _ = check.Suite(&S{})

func (s *S) TestSucessAndParameters(c *check.C) {
	actions := []*Action{
		{
			Forward: func(ctx FWContext) (Result, error) {
				c.Assert(ctx.Params, check.DeepEquals, []interface{}{"hello"})
				return "ok", nil
			},
		},
	}
	pipeline := NewPipeline(actions...)
	err := pipeline.Execute("hello")
	c.Assert(err, check.IsNil)
}

func (s *S) TestRollback(c *check.C) {
	actions := []*Action{
		{
			Forward: func(ctx FWContext) (Result, error) {
				return "ok", nil
			},
			Backward: func(ctx BWContext) {
				c.Assert(ctx.Params, check.DeepEquals, []interface{}{"hello", "world"})
				c.Assert(ctx.FWResult, check.DeepEquals, "ok")
			},
		},
		&errorAction,
	}
	pipeline := NewPipeline(actions...)
	err := pipeline.Execute("hello", "world")
	c.Assert(err, check.NotNil)
	c.Assert(err.Error(), check.Equals, "Failed to execute.")
}

func (s *S) TestRollbackUnrollbackableAction(c *check.C) {
	actions := []*Action{
		&helloAction,
		&unrollbackableAction,
		&errorAction,
	}
	pipeline := NewPipeline(actions...)
	err := pipeline.Execute("hello")
	c.Assert(err, check.NotNil)
	c.Assert(err.Error(), check.Equals, "Failed to execute.")
}

func (s *S) TestExecuteNoActions(c *check.C) {
	pipeline := NewPipeline()
	err := pipeline.Execute()
	c.Assert(err, check.NotNil)
	c.Assert(err.Error(), check.Equals, "No actions to execute.")
}

func (s *S) TestExecuteActionWithNilForward(c *check.C) {
	var executed bool
	actions := []*Action{
		{
			Forward: func(ctx FWContext) (Result, error) {
				return "ok", nil
			},
			Backward: func(ctx BWContext) {
				executed = true
			},
		},
		{
			Forward:  nil,
			Backward: nil,
		},
	}
	pipeline := NewPipeline(actions...)
	err := pipeline.Execute()
	c.Assert(err, check.NotNil)
	c.Assert(err.Error(), check.Equals, "All actions must define the forward function.")
	c.Assert(executed, check.Equals, true)
}

func (s *S) TestExecuteMinParams(c *check.C) {
	var executed bool
	actions := []*Action{
		{
			Forward: func(ctx FWContext) (Result, error) {
				return "ok", nil
			},
			Backward: func(ctx BWContext) {
				executed = true
			},
			MinParams: 0,
		},
		{
			Forward: func(ctx FWContext) (Result, error) {
				return "ok", nil
			},
			MinParams: 1,
		},
	}
	pipeline := NewPipeline(actions...)
	err := pipeline.Execute()
	c.Assert(err, check.NotNil)
	c.Assert(err.Error(), check.Equals, "Not enough parameters to call Action.Forward.")
	c.Assert(executed, check.Equals, true)
}

func (s *S) TestResult(c *check.C) {
	actions := []*Action{
		{
			Forward: func(ctx FWContext) (Result, error) {
				return "ok", nil
			},
			Backward: func(ctx BWContext) {
			},
		},
	}
	pipeline := NewPipeline(actions...)
	err := pipeline.Execute()
	c.Assert(err, check.IsNil)
	r := pipeline.Result()
	c.Assert(r, check.Equals, "ok")
}
