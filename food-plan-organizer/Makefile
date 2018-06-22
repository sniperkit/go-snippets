download:
	curl -L https://www.ars.usda.gov/SP2UserFiles/Place/12354500/Data/SR27/dnload/sr27asc.zip -o sr27asc.zip
	unzip sr27asc.zip -d data
	rm sr27asc.zip

run_database:
	go run cmd/parse/bindata.go cmd/parse/main.go -data.dir data -database sr27.db

build_database:
	cd cmd/parse && go build

run_server:
	go run cmd/serve/bindata.go cmd/serve/main.go -sr27.database sr27.db -user.database user.db

build_server:
	cd cmd/serve && go build

atom:
	/Applications/Atom-Shell.app/Contents/MacOS/Atom .