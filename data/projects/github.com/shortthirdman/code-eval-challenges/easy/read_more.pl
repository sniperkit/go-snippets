use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    chomp $line;
    if ( length $line > 55 ) {
        my $p = 40;
        for ( my $i = 39 ; $i >= 0 ; $i-- ) {
            if ( substr( $line, $i, 1 ) eq ' ' ) {
                $p = $i;
                last;
            }
        }
        $line = substr( $line, 0, $p ) . "... <Read More>";
    }
    print "$line\n";
}
close(INFILE);
