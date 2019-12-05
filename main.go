package main

import _ "site-monitor/init"

import "site-monitor/jobs"

func main() {
	jobs.CheckUrl()
}
