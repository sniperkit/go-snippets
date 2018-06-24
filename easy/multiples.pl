use strict;
use POSIX;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    next if ( $line =~ m/^\s$/ );
    my ( $n, $m ) = ( $line =~ /(\d+),(\d+)/ );
    my $c =
      $n - floor( $n * 2**( -log($m) / log(2) ) ) * 2**( log($m) / log(2) );
    printf( "%d\n", ( $c > 0 ) ? $n - $c + $m : $n );
}
close(INFILE);
