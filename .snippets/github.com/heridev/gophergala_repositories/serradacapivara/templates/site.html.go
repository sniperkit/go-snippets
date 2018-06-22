<!DOCTYPE html>
<html lang="pt-br">
    <head>
        <title>{{.Name}}</title>
        <meta charset="utf-8" />
        <link type="text/css" rel="stylesheet" href="/static/css/layout.css"  media="screen,projection"/>
        <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1.0, user-scalable=no"/>
        <script src="https://maps.googleapis.com/maps/api/js?v=3.exp&signed_in=true"></script>
        <script>
            var map;
            function initialize() {
            	var mapOptions = {
            		scrollwheel: false,
            		zoom: 14,
            		center: new google.maps.LatLng({{.Latitude}}, {{.Longitude}})
            	};
            
            	map = new google.maps.Map(document.getElementById('map-canvas'), mapOptions);
            
            	marker = new google.maps.Marker({
                    position: new google.maps.LatLng({{.Latitude}}, {{.Longitude}}),
                    map: map,
                    title:{{.Name}}
                  });
            }
            
            google.maps.event.addDomListener(window, 'load', initialize);
            
        </script>
    </head>
    <body class="index">
        <div class="navbar-fixed">
            <nav class="transparent" role="navigation">
                <div class="container">
                    <div class="nav-wrapper">
                        <a id="logo-container" href="javascript:history.back();" class="brand-logo thin-text">
                        {{.Name}}
                        </a>
                        <ul class="side-nav">
                            <li><a href="/#index-banner"><i class="mdi-action-search"></i></a></li>
                            <li><a href="#" data-activates="slide-out" class="button-menu"><i class="mdi-navigation-more-vert"></i></a></li>
                        </ul>
                        <ul id="slide-out" class="side-nav full">
                            <li><a href="/about">ABOUT US</a></li>
                            <li><a href="/">SEARCH</a></li>
                            <li><a href="/map">SITES</a></li>
                        </ul>
                        <a href="javascript:history.back();">
                        <i class="mdi-navigation-arrow-back"></i>
                        </a>
                    </div>
                </div>
            </nav>
        </div>
        <div class="parallax-container valign-wrapper">
            <div class="section no-pad-bot">
                <div class="container">
                    <div class="row center">
                        <h1 class="header col s12 white-text">
                            {{.Name}}
                        </h1>
                        <a href="#info" id="download-button" class="btn-large waves-effect waves-light deep-orange darken-4">VIEW GALLERY</a>
                        <p class="col s12 white-text">
                        <div class="col s4">
                            <p class="white-text center"><b>National Park:</b><br>
                                {{.NationalPark}}
                            </p>
                        </div>
                        <div class="col s4">
                            <p class="white-text center"><b>City:</b><br>
                                {{.City}}
                            </p>
                        </div>
                        <div class="col s4">
                            <p class="white-text center"><b>Year of discovery:</b><br>
                                {{.YearOfDiscovery}}
                            </p>
                        </div>
                        </p>

                        <p class="col s12 white-text">
                        <div class="col s4">
                            <p class="white-text center"><b>Circuit:</b><br>
                                {{.Circuit}}
                            </p>
                        </div>
                        <div class="col s4">
                            <p class="white-text center"><b>Location:</b><br>
                                {{.Location}}
                            </p>
                        </div>
                        <div class="col s4">
                            <p class="white-text center"><b>Geo:</b><br>
                                ({{.Latitude}}, {{.Longitude}})
                            </p>
                        </div>
                        </p>
                    </div>
                </div>
            </div>
            <div class="parallax"><img src="/static/assets/background-serra.jpg"></div>
        </div>
        <div class="container">
            <div class="section">
                <div class="row" id="info">
                    <div class="col s12 center">
                        <h3><i class="mdi-action-help deep-orange-text darken-4"></i></h3>
                        <h4>Info</h4>
                        <p class="left-align light">
                            {{.Description}}
                        </p>
                    </div>
                </div>
            </div>
        </div>
        <div class="container margin-bottom">
            <section id="gallery">
                <div class="row">
                    {{ range $i, $url := .Pictures }}
                    <img class="materialboxed col s6 m4 l2 responsive-img" src="{{$url}}" >
                    {{ end }}
                </div>
            </section>
        </div>
        <div id="map-canvas"></div>
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
            	$('.materialboxed').materialbox();
            });
        </script>
    </body>
</html>
