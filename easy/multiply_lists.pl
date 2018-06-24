use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    chomp($line);
    my ( $s, $t ) = split( / \| /, $line );
    my @s = split( / /, $s );
    my @t = split( / /, $t );
    for ( my $i = 0 ; $i < scalar @s ; $i++ ) {
        @s[$i] *= @t[$i];
    }
    printf( "%s\n", join( ' ', @s ) );
}
close(INFILE);
