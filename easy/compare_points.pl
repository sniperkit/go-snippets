use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    chomp($line);
    my ( $x1, $y1, $x2, $y2 ) = split( / /, $line );
    print(
          ( $x1 == $x2 )
        ? ( $y1 == $y2 ? "here\n" : ( $y1 < $y2 ? "N\n" : "S\n" ) )
        : (
            $y1 == $y2
            ? ( $x1 < $x2 ? "E\n" : "W\n" )
            : (
                $x1 < $x2
                ? ( $y1 < $y2 ? "NE\n" : "SE\n" )
                : ( $y1 < $y2 ? "NW\n" : "SW\n" )
            )
        )
    );
}
close(INFILE);
