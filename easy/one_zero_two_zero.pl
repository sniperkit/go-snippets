use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    my ( $n, $a ) = ( $line =~ /(\d+) (\d+)/ );
    my $r = 0;
    for ( my $i = 1 << $n ; $i <= $a ; $i++ ) {
        my ( $m, $b ) = ( $n, $i );
        while ( $b > 0 && $m >= 0 ) {
            $m-- if ( $b & 1 ) == 0;
            $b >>= 1;
        }
        $r++ if $m == 0;
    }
    print "$r\n";
}
close(INFILE);
