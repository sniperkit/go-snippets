use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    chomp($line);
    my ( $s, $m ) = split( / \| /, $line );
    my @t = split( / /, $s );
    while ( scalar @t > 1 ) {
        my $n = $m % ( scalar @t ) - 1;
        $n >= 0 ? splice( @t, $n, 1 ) : pop @t;
    }
    print "@t[0]\n";
}
close(INFILE);
