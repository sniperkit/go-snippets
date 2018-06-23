<!DOCTYPE html>
<html lang="pt-br">
<head>
  <title>About</title>
  <meta charset="utf-8" />
  <link type="text/css" rel="stylesheet" href="/static/css/layout.css"  media="screen,projection"/>
  <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1.0, user-scalable=no"/>
</script>
</head>
<body class="index">

  <div class="navbar-fixed">
    <nav class="transparent" role="navigation">
      <div class="container">
        <div class="nav-wrapper">
          <a id="logo-container" href="/" class="brand-logo thin-text">
            About
          </a>

          <ul class="side-nav">
            <li><a href="#"  data-activates="slide-out" class="button-menu"><i class="mdi-navigation-more-vert"></i></a></li>
          </ul>

         <ul id="slide-out" class="side-nav full">
              <li><a href="/about">ABOUT</a></li>
              <li><a href="/#index-banner">SEARCH</a></li>
          </ul>

          <a href="/">
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
            Special places need to be shown
          </h1>
        </div>
      </div>
    </div>
    <div class="parallax"><img src="/static/assets/background-serra.jpg"></div>
  </div>


  <div class="container">

    <h2>About us</h2>
    <p>Serra da Capivara National Park (Portuguese: Parque Nacional Serra da Capivara) is a national park in the Northeastern region of Brazil. It has many prehistoric paintings. The park was created to protect the prehistoric artifacts and paintings found there. It became a World Heritage Site in 1991. Its head archaeologist is Niède Guidon. Its best known archaeological site is Pedra Furada. <br><br>

      It is located in northeast state of Piauí, between latitudes 8° 26' 50" and 8° 54' 23" south and longitudes 42° 19' 47" and 42° 45' 51" west. It falls within the municipal areas of São Raimundo Nonato, São João do Piauí, Coronel José Dias and Canto do Buriti. It has an area of 1291.4 square kilometres (319,000 acres). The area has the largest concentration of prehistoric sites in the Americas. Scientific studies confirm that the Capivara mountain range was densely populated in the pre-Columbian Era.</p>

    </div>


    <div class="primaryColor">
      <div class="container">
        <h2>The Team</h2>
        <div class="row">
          <div class="col l4">
            <h5 class="center">
              <a href="https://github.com/dannluciano"><img src="http://www.gravatar.com/avatar/9486298aceeee0f05f029c760d89a248.jpg?r=g&s=120" class="responsive-img"><br>
              @dannluciano</a>
            </h5>
          </div>
          <div class="col l4">
            <h5 class="center">
              <a href="http://github.com/dimiro1"><img src="http://www.gravatar.com/avatar/bcf55979bb17c9e73dddd9d9eb81c5f3.jpg?r=g&s=120" class="responsive-img"><br>
              @dimiro1</a>
            </h5>
          </div>
          <div class="col l4">
            <h5 class="center">
              <a href="http://sitma.com.br/"><img src="http://www.gravatar.com/avatar/0d8de0d32a891b4f282f8116354612cb.jpg?r=g&s=120" class="responsive-img"><br>
              @italobrasiil</a>
            </h5>
          </div>
        </div>
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
