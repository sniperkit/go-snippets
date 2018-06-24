use strict;

sub numzero {
    my @s = @{ shift @_ };
    my ( $c, $z ) = @_;
    return ( ( $z == 0 ) ? 1 : 0 ) if ( $c == 0 );
    return (0) if scalar @s < $c || $z + $c * @s[0] > 0 || $z + $c * @s[-1] < 0;
    my @t  = @s;
    my $t0 = shift @t;
    return ( numzero( \@t, $c - 1, $z + $t0 ) + numzero( \@t, $c, $z ) );
}

open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    chomp($line);
    my @r = sort { $a <=> $b } split( /,/, $line );
    printf( "%d\n", numzero( \@r, 4, 0 ) );
}
close(INFILE);
