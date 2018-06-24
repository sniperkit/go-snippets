<!DOCTYPE html>
<html lang="pt-br">
    <head>
        <title>Serra da Capivara</title>
        <meta charset="utf-8" />
        <link type="text/css" rel="stylesheet" href="/static/css/layout.css"  media="screen,projection"/>
        <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1.0, user-scalable=no"/>
    </head>
    <body class="index">
        <div class="navbar-fixed">
            <nav class="transparent" role="navigation">
                <div class="container">
                    <div class="nav-wrapper">
                        <a id="logo-container" href="/" class="brand-logo thin-text">Serra da Capivara</a>
                        <ul id="slide-out" class="side-nav full">
                            <li><a href="/about">ABOUT</a></li>
                            <li><a href="/#index-banner">SEARCH</a></li>
                            <li><a href="/map">SITES</a></li>
                        </ul>
                        <a href="#" data-activates="slide-out" class="button-menu">
                        <i class="mdi-navigation-menu"></i>
                        </a>
                    </div>
                </div>
            </nav>
        </div>
        <div class="info-page">
            <div class="container">
                <div class="row row-container">
                    <div class="col s12 m6">
                        <img src="/static/assets/brand-home.png" class="responsive-img">
                    </div>
                    <div class="col s12 m6">
                        <h1 class="white-text">Serra da Capivara</h1>
                        <p class="white-text opaque">
                            Serra da Capivara National Park is a national park in the Northeastern region of Brazil. It has many prehistoric paintings. The park was created to protect the prehistoric artifacts and paintings found there. It became a World Heritage Site in 1991.
                        </p>
                        <a href="/about" class="waves-effect waves-dark btn white red-text">I want to know more</a>
                    </div>
                </div>
                <a class="btn-floating btn-large waves-effect waves-light white">
                <i class="mdi-navigation-expand-more grey-text"></i>
                </a>
            </div>
        </div>
        <div class="parallax-container valign-wrapper">
            <div class="section no-pad-bot">
                <div class="container">
                    <div class="row center">
                        <h1 class="header col s12 white-text">
                            Special places need to be shown
                        </h1>
                        <a href="/map" id="download-button" class="btn-large waves-effect waves-light red">See all sites</a>
                    </div>
                </div>
            </div>
            <div class="parallax">
                <img src="/static/assets/background-serra.jpg">
            </div>
        </div>
        <div class="section no-pad-bot" id="index-banner">
            <div class="container">
                <br><br>
                <div class="row center">
                    <h5 class="header col s12 light">Search sites</h5>
                </div>
                <div class="row center">
                    <div class="input-field col s12">
                        <form method="GET" action="/search">
                            <input id="search_input" type="text" class="validate" name="q" />
                            <label for="search_input">Search for Title.</label>
                            <input type="submit" value="SEARCH NOW" class="btn-large waves-effect waves-light red"/>
                        </form>
                    </div>
                </div>
                <br><br>
            </div>
        </div>
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
