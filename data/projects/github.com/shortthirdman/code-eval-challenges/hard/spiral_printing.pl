use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    chomp $line;
    my ( $row, $col, $s ) = split( /;/, $line );
    my @t = split( / /, $s );
    my ( $i, $tn, $tw, $j, $te, $ts ) = ( 0, 0, 0, 1, $col - 1, $row - 1 );
    print "$t[0]";
    while (1) {
        for ( ; $j <= $te ; $j++ ) {
            print " $t[$i * $col + $j]";
        }
        $j--;
        $i++;
        $tn++;
        last if ( $i > $ts );
        for ( ; $i <= $ts ; $i++ ) {
            print " $t[$i * $col + $j]";
        }
        $i--;
        $j--;
        $te--;
        last if ( $j < $tw );
        for ( ; $j >= $tw ; $j-- ) {
            print " $t[$i * $col + $j]";
        }
        $j++;
        $i--;
        $ts--;
        last if ( $i < $tn );
        for ( ; $i >= $tn ; $i-- ) {
            print " $t[$i * $col + $j]";
        }
        $i++;
        $j++;
        $tw++;
        last if ( $j > $te );
    }
    print "\n";
}
close(INFILE);
