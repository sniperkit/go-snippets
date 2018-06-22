package main

type Action interface {
	GetButtonTitle() string
	GetLayoutName() string
	Run()
}
