use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    chomp($line);
    my @s =
      map { substr( $_, -1 ) . substr( $_, 1, -1 ) . substr( $_, 0, 1 ) }
      split( / /, $line );
    printf( "%s\n", join( ' ', @s ) );
}
close(INFILE);
