<!DOCTYPE html>
<html lang="pt-br">
    <head>
        <title>Serra da Capivara</title>
        <meta charset="utf-8" />
        <link type="text/css" rel="stylesheet" href="/static/css/layout.css"  media="screen,projection"/>
        <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1.0, user-scalable=no"/>
        <script src="https://maps.googleapis.com/maps/api/js?v=3.exp&signed_in=true"></script>
        <script>
            var map;
            function initialize() {
                var mapOptions = {
                    scrollwheel: true,
                    zoom: 12,
                    center: new google.maps.LatLng(-8.836630, -42.555310)
                };

                map = new google.maps.Map(document.getElementById('map-canvas'),
                    mapOptions);

                {{ range $i, $site := . }}
                    var marker = new google.maps.Marker({
                        position: new google.maps.LatLng({{$site.Latitude}}, {{$site.Longitude}}),
                        map: map,
                        title:{{$site.Name}}
                    });

                    google.maps.event.addListener(marker, 'click', function() {
                        document.location.href = "/site/{{$site.Id}}"
                    });
                {{ end }}
            }
            
            google.maps.event.addDomListener(window, 'load', initialize);
            
        </script>
    </head>
    <body class="index">
        <div class="navbar-fixed">
            <nav class="transparent" role="navigation">
                <div class="container">
                    <div class="nav-wrapper">
                        <a id="logo-container" href="#" class="brand-logo thin-text">
                        All Sites
                        </a>
                        <ul class="side-nav">
                            <li><a href="/#index-banner"><i class="mdi-action-search"></i></a></li>
                            <li><a href="#"  data-activates="slide-out" class="button-menu"><i class="mdi-navigation-more-vert"></i></a></li>
                        </ul>
                        <ul id="slide-out" class="side-nav full">
                            <li><a href="/about">ABOUT US</a></li>
                            <li><a href="/">SEARCH</a></li>
                            <li><a href="/map">SITES</a></li>
                        </ul>
                        <a href="#">
                        <i class="mdi-navigation-arrow-back"></i>
                        </a>
                    </div>
                </div>
            </nav>
        </div>
        <div id="map-canvas" class="full"></div>
        <footer class="page-footer">
            <div class="footer-copyright">
                <div class="container">
                    Made with Go, Love, Doritos and <a class="deep-orange-text darken-4" href="http://materializecss.com">Materialize</a>
                </div>
            </div>
        </footer>
        <!--Import jQuery before materialize.js-->
        <script type="text/javascript" src="/static/js/jquery.js"></script>
        <script type="text/javascript" src="/static/js/materialize.min.js"></script>
        <script type="text/javascript">
            $(document).ready(function(){
                $('.parallax').parallax();
                $('.button-menu').sideNav();
            });
        </script>
    </body>
</html>
