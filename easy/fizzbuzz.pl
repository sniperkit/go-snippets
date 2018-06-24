use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    chomp($line);
    my ( $f, $b, $n ) = split( / /, $line );
    my @r;
    for ( my $i = 1 ; $i <= $n ; $i++ ) {
        @r[ $i - 1 ] =
          $i % $f == 0 ? $i % $b == 0 ? "FB" : "F" : $i % $b == 0 ? "B" : $i;
    }
    printf( "%s\n", join( ' ', @r ) );
}
close(INFILE);
