use strict;

sub g2d {
    my $a = $_[0];
    $a ^= $a >> 4;
    $a ^= $a >> 2;
    $a ^= $a >> 1;
    return ($a);
}
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    chomp($line);
    my @r = map { g2d( oct( "0b" . $_ ) ) } split( / \| /, $line );
    printf( "%s\n", join( " | ", @r ) );
}
close(INFILE);
