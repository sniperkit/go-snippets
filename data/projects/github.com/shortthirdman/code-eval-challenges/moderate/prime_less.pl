use strict;

my @p = ( 2, 3 );
my $c = 5;

sub prime {
    my ( $a, $i ) = ( $_[0], 1 );
    while ( @p[$i] * @p[$i] <= $a ) {
        return if ( $a % @p[$i] == 0 );
        $i++;
    }
    push @p, $c;
    return;
}

open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    while ( $c < $line ) {
        prime($c);
        $c += 2;
    }
    printf( "%s\n", join( ',', grep { $_ < $line } @p ) );
}
close(INFILE);
