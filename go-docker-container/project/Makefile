run:
	go run *.go
gomod-exp:
	export GO111MODULE=on
gobuild:
	GOOS=linux GOARCH=amd64 go build -o k8golangtest
dockerbuild:
	docker build -t k8golangtest .
dockerbuildandpush:
	docker build -t k8golangtest .
	docker tag k8golangtest americanwonton/k8golangtest
	docker push americanwonton/k8golangtest
dockerrun:
	docker run -it -p 8080:8080 k8golangtest
dockerrundetached:
	docker run -d -p 8080:8080 k8golangtest
dockerrunitvolume:
	docker run -it -p 8080:8080 -v photo-images:/static/images k8golangtest
dockerrundetvolume:
	docker run -d -p 8080:8080 -v photo-images:/static/images k8golangtest
dockertagimage:
	docker tag k8golangtest americanwonton/k8golangtest
dockerimagepush:
	docker push americanwonton/k8golangtest
dockerallpush:
	docker tag k8golangtest americanwonton/k8golangtest
	docker push americanwonton/k8golangtest
dockerseeshell:
	docker run -it k8golangtest sh