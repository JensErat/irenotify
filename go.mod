module github.com/JensErat/irenotify

go 1.19

require (
	github.com/go-logr/stdr v1.2.2
	github.com/gokyle/fswatch v0.0.0-20121217010029-1dbdf8320a69
	github.com/patrickmn/go-cache v2.1.0+incompatible
)

require github.com/go-logr/logr v1.2.2 // indirect

// https://github.com/gokyle/fswatch/pull/1
replace github.com/gokyle/fswatch => github.com/JensErat/fswatch v0.0.0-20230213220220-0486f1e76cef
