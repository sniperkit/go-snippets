use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    chomp $line;
    my ( $i, $r );
    for ( $i = 0 ; $i < 100 ; $i++ ) {
        $r = reverse($line);
        last if ( $r == $line );
        $line += $r;
    }
    printf( "%s\n", $r == $line ? $i . ' ' . $line : "not found" );
}
close(INFILE);
