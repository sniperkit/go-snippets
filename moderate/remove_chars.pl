use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    my ( $l, $r ) = split( /, /, $line );
    $l =~ s/[$r]//g;
    print "$l\n";
}
close(INFILE);
