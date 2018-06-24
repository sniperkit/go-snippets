use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    my ( $a, $b ) = split( / \| /, $line );
    my $r = 0;
    for ( my $i = 0 ; $i < length $a ; $i++ ) {
        $r++ if substr( $a, $i, 1 ) ne substr( $b, $i, 1 );
    }
    if ( $r == 0 ) {
        print "Done\n";
    }
    elsif ( $r < 3 ) {
        print "Low\n";
    }
    elsif ( $r < 5 ) {
        print "Medium\n";
    }
    elsif ( $r < 7 ) {
        print "High\n";
    }
    else {
        print "Critical\n";
    }
}
close(INFILE);
