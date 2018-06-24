use strict;
use bigint;

open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    my $r = 0;
    for ( my $i = 0 ; $i < $line / 2 ; $i++ ) {
        $r +=
          (-1)**$i *
          ( 0 + $line )->bnok($i) *
          ( 11 * $line / 2 - 1 - 10 * $i )->bnok( $line - 1 );
    }
    print "$r\n";
}
close(INFILE);
