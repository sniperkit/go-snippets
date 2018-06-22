# alpine-php7

    docker build -t ncc0706/alpine-php7 .
    
    docker run --rm -p 80:80 -v $(pwd):/application ncc0706/alpine-php7