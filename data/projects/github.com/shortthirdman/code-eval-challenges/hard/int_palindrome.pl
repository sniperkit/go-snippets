use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    $line =~ /(\d+) (\d+)/;
    my ( $n, $l, $r ) = ( 0, $1, $2 );
    for ( my $i = $l ; $i <= $r ; $i++ ) {
        my $prev = -1;
        for ( my $j = $i ; $j <= $r ; $j++ ) {
            my $p;
            if ( $prev > 0 ) {
                $p = $prev;
                $p++ if ( $j eq reverse($j) );
            }
            else {
                $p = 0;
                for ( my $k = $i ; $k <= $j ; $k++ ) {
                    $p++ if ( $k eq reverse($k) );
                }
            }
            $n++ if ( ( $p & 1 ) == 0 );
            $prev = $p;
        }
    }
    print "$n\n";
}
close(INFILE);
