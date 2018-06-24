use strict;

my @p = ( 2, 3 );

sub prime {
    my ( $a, $i ) = ( $_[0], 1 );
    while ( @p[$i] * @p[$i] <= $a ) {
        return (0) if ( $a % @p[$i] == 0 );
        $i++;
    }
    return (1);
}

my $c = 5;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    while ( 2**$c - 1 < $line ) {
        $c += 2 while ( !prime($c) );
        push @p, $c;
        $c += 2;
    }
    my @r = grep { $_ < $line } map { 2**$_ - 1 } @p;
    printf( "%s\n", join( ', ', @r ) );
}
close(INFILE);
