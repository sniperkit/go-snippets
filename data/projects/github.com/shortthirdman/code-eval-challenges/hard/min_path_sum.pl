use strict;
use List::Util qw[min];
my ( $n, $i ) = ( 0, 0 );
my @b;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    chomp $line;
    if ( $n == 0 ) {
        ( $n, $i ) = ( $line, $line );
        next;
    }
    my @s = split( /,/, $line );
    if ( $i == $n ) {
        @b = ();
        $b[0] = $s[0];
        for ( my $j = 1 ; $j < $n ; $j++ ) {
            $b[$j] = $b[ $j - 1 ] + $s[$j];
        }
    }
    else {
        $b[0] += $s[0];
        for ( my $j = 1 ; $j < $n ; $j++ ) {
            $b[$j] = min( $b[ $j - 1 ], $b[$j] ) + $s[$j];
        }
    }
    $i--;
    if ( $i == 0 ) {
        print "$b[-1]\n";
        $n = 0;
    }
}
close(INFILE);
