build-restart:
	go build -o kubectl-check_restart -tags=restarted ./run
build-cleaner:
	go build -o kubectl-cleaner -tags=cleaner ./run