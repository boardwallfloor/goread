rerun :
	go build
	./goread $(flags)
reprof :
	go build
	./goread -cpuprofile=perf.prof
profile : 
	go tool pprof -http=localhost:9000 goread perf.prof
build :
	go mod tidy
	go build
	./goread