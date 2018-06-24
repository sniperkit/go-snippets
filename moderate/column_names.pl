use strict;
use integer;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    next if ( $line =~ m/^\s$/ );
    $line--;
    my $r = chr( $line % 26 + ord 'A' );
    if ( $line >= 26 ) {
        $line = $line / 26 - 1;
        $r    = chr( $line % 26 + ord 'A' ) . $r;
        if ( $line >= 26 ) {
            $line = $line / 26 - 1;
            $r    = chr( $line % 26 + ord 'A' ) . $r;
        }
    }
    print "$r\n";
}
close(INFILE);
