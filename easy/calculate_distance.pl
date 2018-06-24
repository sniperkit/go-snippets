use strict;
use integer;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    next if ( $line =~ m/^\s$/ );
    my ( $x, $y, $p, $q ) =
      ( $line =~ /(-?\d+), (-?\d+)\) \((-?\d+), (-?\d+)/ );
    my ( $dx, $dy ) = ( $p - $x, $q - $y );
    printf( "%d\n", sqrt( $dx * $dx + $dy * $dy ) );
}
close(INFILE);
