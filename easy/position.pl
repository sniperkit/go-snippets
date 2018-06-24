use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    my ( $x, $p, $q ) = ( $line =~ /(\d+),(\d+),(\d+)/ );
    print( ( $x & 1 << $p - 1 ) > 0 == ( $x & 1 << $q - 1 ) > 0
        ? "true\n"
        : "false\n" );
}
close(INFILE);
