use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    print "1";
    for ( my $i = 1 ; $i < $line ; $i++ ) {
        print " 1";
        my $r = 1;
        for ( my $j = 0 ; $j < $i ; $j++ ) {
            $r = ( $r * ( $i - $j ) ) / ( $j + 1 );
            print " $r";
        }
    }
    print "\n";
}
close(INFILE);
