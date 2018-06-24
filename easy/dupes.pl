use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    chomp($line);
    my $l = -1;
    foreach my $i ( split( /,/, $line ) ) {
        if ( $l == -1 ) {
            print "$i";
        }
        elsif ( $i != $l ) {
            print ",$i";
        }
        $l = $i;
    }
    print "\n";
}
close(INFILE);
