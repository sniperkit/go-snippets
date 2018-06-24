use strict;
use integer;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    chomp $line;
    my @s = split( / /, $line );
    my @r;
    my $n = sqrt( scalar @s );
    for ( my $i = 0 ; $i < $n ; $i++ ) {
        for ( my $j = $n - 1 ; $j >= 0 ; $j-- ) {
            push @r, $s[ $j * $n + $i ];
        }
    }
    printf( "%s\n", join( ' ', @r ) );
}
close(INFILE);
