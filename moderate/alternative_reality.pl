use strict;
my @m = ( 50, 25, 10, 5, 1 );

sub alter {
    my ( $n, $p ) = @_;
    return (1) if ( $n == 0 );
    $p += 1 while ( $m[$p] > $n );
    return (1) if ( $m[$p] == 1 );
    return ( alter( $n - $m[$p], $p ) + alter( $n, $p + 1 ) );
}
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    printf( "%d\n", alter( $line, 0 ) );
}
close(INFILE);
