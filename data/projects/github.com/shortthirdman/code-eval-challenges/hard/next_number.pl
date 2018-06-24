use strict;
use integer;

sub dsig {
    my ( $a, $d ) = ( $_[0], 0 );
    while ( $a > 0 ) {
        my $r = $a % 10;
        if ( $r > 0 ) {
            $d += 1 << ( 3 * $r );
        }
        $a /= 10;
    }
    return ($d);
}

open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    next if ( $line =~ m/^\s$/ );
    my $b = $line + 9;
    my $s = dsig($line);
    while ( $s != dsig($b) ) {
        $b += 9;
    }
    printf "$b\n";
}
close(INFILE);
