use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    chomp $line;
    my @s = split( /;/, $line );
    my $n = pop @s;
    my @t = split( /,/, $s[0] );
    for ( my $c = $n ; $c <= scalar @t ; $c += $n ) {
        for ( my $i = 1 ; 2 * $i <= $n ; $i++ ) {
            ( $t[ $c - $i ], $t[ $c - $n + $i - 1 ] ) =
              ( $t[ $c - $n + $i - 1 ], $t[ $c - $i ] );
        }
    }
    printf( "%s\n", join( ',', @t ) );
}
close(INFILE);
