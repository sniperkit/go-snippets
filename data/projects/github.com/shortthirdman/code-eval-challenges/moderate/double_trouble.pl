use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    chomp($line);
    my ( $o, $r ) = ( ( length $line ) / 2, 1 );
    for ( my $i = 0 ; $i < $o ; $i++ ) {
        if (
            (
                   substr( $line, $i, 1 ) eq "A"
                && substr( $line, $i + $o, 1 ) eq "B"
            )
            || (   substr( $line, $i, 1 ) eq "B"
                && substr( $line, $i + $o, 1 ) eq "A" )
          )
        {
            $r = 0;
            last;
        }
        if (   substr( $line, $i, 1 ) eq "*"
            && substr( $line, $i + $o, 1 ) eq "*" )
        {
            $r <<= 1;
        }
    }
    print "$r\n";
}
close(INFILE);
