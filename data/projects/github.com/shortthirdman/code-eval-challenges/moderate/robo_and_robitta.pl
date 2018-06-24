use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    $line =~ /(\d+)x(\d+) \| (\d+) (\d+)/;
    my ( $m, $n, $s, $t ) = ( $1, $2, $3, $4 );
    my $r = 0;
    while ( $n > $t ) {
        $r += $m;
        ( $m, $n, $s, $t ) = ( $n - 1, $m, $n - $t, $s );
    }
    printf( "%d\n", $r + $s );
}
close(INFILE);
