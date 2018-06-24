use strict;
my $morse = "ETIANMSURWDKGOHVF L PJBXCYZQ  54 3   2       16       7   8 90";
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    chomp $line;
    my $m = 1;
    foreach my $i ( split( //, $line ) ) {
        if ( $i eq '.' ) { $m *= 2 }
        elsif ( $i eq '-' ) { $m = 2 * $m + 1 }
        else {
            if ( $m == 1 ) { print " " }
            else {
                print substr( $morse, $m - 2, 1 ) if ( $m < 64 );
                $m = 1;
            }
        }
    }
    print substr( $morse, $m - 2, 1 ) if ( $m > 1 && $m < 64 );
    print "\n";
}
close(INFILE);
