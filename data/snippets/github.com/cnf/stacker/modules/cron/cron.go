package cron

import (
	"github.com/cnf/stacker/engine"
)

// Cron object
type Cron struct {
	// Config *Config
}

// NewCron returns a new cron object
func NewCron() engine.Action {
	return &Cron{}
}

// RegisterAction registers this action with the engine
func RegisterAction() {
	engine.RegisterAction("cron", NewCron)
}

// Setup ...
func (c *Cron) Setup(eng *engine.Engine, cfg engine.Config) error {
	// c.Config = cfg.(*Config)
	return nil
}

// Start the cron runner
func (c *Cron) Start() error {
	return nil
}

// Stop the cron runner
func (c *Cron) Stop() error {
	return nil
}

// Add things
func (c *Cron) Add() error {
	return nil
}

// Remove things
func (c *Cron) Remove() error {
	return nil
}

// RemoveAll things
func (c *Cron) RemoveAll() error {
	return nil
}
