use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    next if ( $line =~ m/^\s$/ );
    chomp($line);
    my @s = split( / /, $line );
    my $n = shift @s;
    my $j = 1;
    if ( $n > 1 ) {
        my %u;
        for ( my $i = 1 ; $i < $n ; $i++ ) {
            my $x = abs( @s[$i] - @s[ $i - 1 ] );
            if ( $x >= $n || $x == 0 || $u{ $x - 1 } ) {
                $j = 0;
                last;
            }
            $u{ $x - 1 } = 1;
        }
    }
    printf( "%s\n", $j ? "Jolly" : "Not jolly" );
}
close(INFILE);
