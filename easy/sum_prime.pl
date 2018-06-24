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

my ( $c, $sum ) = ( 5, 5 );
for ( my $i = 2 ; $i < 1000 ; $i++ ) {
    $c += 2 while ( !prime($c) );
    push @p, $c;
    ( $sum, $c ) = ( $sum + $c, $c + 2 );
}
print "$sum\n";
