package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	idstring = "http://golang.org/pkg/http/#ListenAndServe"
)

var (
	help     = flag.Bool("h", false, "show this help.")
	host     = flag.String("host", "localhost:8080", "listening port and hostname.")
	interval = flag.Int("interval", 10, "interval in minutes between checks")
	testFile = flag.String("testfile", "~/mnt/serenity/.bashrc", "file to stat to determine is mount effective")
)

func usage() {
	fmt.Fprintf(os.Stderr, "reminder\n")
	flag.PrintDefaults()
	os.Exit(2)
}

var tpl *template.Template

func main() {
	flag.Usage = usage
	flag.Parse()
	if *help {
		usage()
	}

	nargs := flag.NArg()
	if nargs > 0 {
		usage()
	}

	cleanTestFile()

	tpl = template.Must(template.New("main").Parse(mainHTML()))
	http.HandleFunc("/", mainHandler)
	go func() {
		if err := http.ListenAndServe(*host, nil); err != nil {
			log.Fatalf("Could not start http server: %v", err)
		}
	}()
	time.Sleep(2 * time.Second)

	url := "http://" + *host
	for {
		if !isMounted(*testFile) {
			if err := exec.Command("xdg-open", url).Run(); err != nil {
				log.Fatalf("Could not open url %v in browser: %v", url, err)
			}
		}
		time.Sleep(time.Duration(*interval) * time.Minute)
	}
}

func isMounted(testFile string) bool {
	if _, err := os.Stat(testFile); err != nil {
		return false
	}
	return true
}

func cleanTestFile() {
	if *testFile == "" || len(*testFile) < 2 {
		log.Fatal("Invalid test file")
	}
	if (*testFile)[0] == '~' {
		e := os.Getenv("HOME")
		if e == "" {
			log.Fatal("no HOME env var")
		}
		*testFile = filepath.Join(e, (*testFile)[1:])
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Server", idstring)
	if err := tpl.Execute(w, nil); err != nil {
		log.Printf("Could not execute template: %v", err)
	}
}

func mainHTML() string {
	return `<!DOCTYPE HTML>
<html>
	<head>
		<title>Reminder</title>
	</head>

	<body>
	<script>
setTimeout(window.close, 10000);
window.onload=function(){notify()};

function enableNotify() {
	if (!(window.webkitNotifications)) {
		alert("Notifications not supported on this browser.");
		return;
	}
	var havePermission = window.webkitNotifications.checkPermission();
	if (havePermission == 0) {
		alert("Notifications already allowed.");
		return;
	}
	window.webkitNotifications.requestPermission();
}

function notify() {
	if (!(window.webkitNotifications)) {
		console.log("Notifications not supported");
		return;
	}
	var havePermission = window.webkitNotifications.checkPermission();
	if (havePermission != 0) {
		console.log("Notifications not allowed.");
		return;
	}
	var notification = window.webkitNotifications.createNotification(
		'',
		'Reminder notification',
		'do the sshfs dance'
	);

	// NOTE: the tab/window needs to be still open for the cancellation
	// of the notification to work.
	notification.onclick = function () {
		this.cancel();
	};

	notification.ondisplay = function(event) {
		setTimeout(function() {
			event.currentTarget.cancel();
		}, 7000);
	};

	notification.show();
} 

	</script>

	<a id="notifyLink" href="#" onclick="enableNotify();return false;">Enable notifications?</a>

	<h2> Bazinga </h2>
	</body>
</html>
`
}
