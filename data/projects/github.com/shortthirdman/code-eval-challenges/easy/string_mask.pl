use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    next if ( $line =~ m/^\s$/ );
    chomp($line);
    my ( $l, $r ) = split( / /, $line );
    for ( my $i = 0 ; $i < length $l ; $i++ ) {
        if ( substr( $r, $i, 1 ) == 1 ) {
            print uc substr( $l, $i, 1 );
        }
        else {
            print lc substr( $l, $i, 1 );
        }
    }
    print "\n";
}
close(INFILE);
