# nginx rtmp

nginx 直播

    docker run -d -p 80:80 -p 1935:1935 --name nginx ncc0706/nginx-rtmp
    
    # ffplay 播放
    ffplay rtmp://52.79.113.198/playback/01.mp4
    
    ## port
    tcp 1935
    tcp 80