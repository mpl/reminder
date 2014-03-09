reminder
========

Silly and convoluted, albeit funny, reminder.
Basically replaces crontab + mail with go http server + webkit notification.
It watches an sshfs mount because that's what I needed it for, but it could be easily generalized to check for user defined commands or whatnot.

Install:
--------

 1. Install Go: http://golang.org/doc/install
 2. go get github.com/mpl/reminder

