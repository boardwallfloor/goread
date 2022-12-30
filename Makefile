rerun :
	go build
	./goread
reprof :
	go build
	./goread -cpuprofile=perf.prof
profile : 
	go tool pprof -http=localhost:9000 goread perf.prof
