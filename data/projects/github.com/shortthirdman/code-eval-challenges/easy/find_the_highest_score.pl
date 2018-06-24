use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    chomp($line);
    my @r;
    foreach my $i ( split( / \| /, $line ) ) {
        if ( @r == 0 ) {
            @r = split( / /, $i );
        }
        else {
            my $jx = 0;
            foreach my $j ( split( / /, $i ) ) {
                @r[$jx] = $j if ( $j > @r[$jx] );
                $jx++;
            }
        }
    }
    printf( "%s\n", join( ' ', @r ) );
}
close(INFILE);
