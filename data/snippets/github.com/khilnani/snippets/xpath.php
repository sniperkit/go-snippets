#!/usr/bin/php -q
<?php

//-------------------------------------------------------

$inputFile = "./xpath.xml";

//-------------------------------------------------------

function processFile( $from )
{
	
  echo "--------------------------------------------------------------------------\n";
  echo "Processing: $from\n";
  echo "--------------------------------------------------------------------------\n";

  $in = file_get_contents($from);
  //var_dump( $in );

  $xml = new SimpleXMLElement($in);
  $xml->registerXPathNamespace('ms', 'http://microsoft.com/catalog.xsd');

  $books = $xml->xpath('//ms:book');
  //var_dump($books);

  foreach ($books as $book) {
    //var_dump($book);

    $author = $book->xpath('.//ms:author');
    $title = $book->xpath('.//ms:title');
    $id = $book->xpath('.//@id');
    echo "Book: $author[0] - $title[0] ($id[0])\n";

    $publish_date = $book->xpath('.//ms:publish_date');
    $publish_date[0][0] = 'TODAY'; // Change the publish date
  }

  $outputXml = $xml->asXML();

  echo "--------------------------------------------------------------------------\n";
  echo "Updated XML (publish_date = 'TODAY')\n";
  echo "--------------------------------------------------------------------------\n";
  echo "$outputXml\n";
  echo "--------------------------------------------------------------------------\n";

}

//-------------------------------------------------------

function start() {
  global $inputFile;
  processFile($inputFile);
}

//-------------------------------------------------------

start();
// var_dump( $hash );

//-------------------------------------------------------

?>
