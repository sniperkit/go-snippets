use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    my $r = 0;
    my ( $ls, $rs ) = split( / \| /, $line );
    map { $r += hex $_ } split( / /, $ls );
    map { $r -= oct( "0b" . $_ ) } split( / /, $rs );
    print( $r > 0 ? "False\n" : "True\n" );
}
close(INFILE);
