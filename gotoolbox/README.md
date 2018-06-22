# Go Toolbox

The Go Toolbox is a little web project where developers of go libraries can submit their open source work and categorize it.
People who are kickstarting projects can see for every category which projects are out there and what are the options.

## Tools used

https://github.com/pilu/fresh - Reload the server everytime a go or ace file changes

https://github.com/codegangsta/negroni - Express style middleware for go

https://github.com/gorilla/mux - Router

https://github.com/markbates/goth/ - Authentication via OAuth

https://github.com/jinzhu/gorm - Dealing with the database

http://www.gorillatoolkit.org/pkg/schema - Parsing web forms

https://github.com/yosssi/ace/ - For templates

## Deployment via Docker

### Building the image

    docker build -t 9elements/gotoolbox .

    docker run -p 8080:8080 -h gotoolbox -d 9elements/gotoolbox

    docker run -p 8080:8080 -it --entrypoint=/bin/bash 9elements/gotoolbox -i

## License

Go Toolbox is released under the [MIT License](http://www.opensource.org/licenses/MIT).
