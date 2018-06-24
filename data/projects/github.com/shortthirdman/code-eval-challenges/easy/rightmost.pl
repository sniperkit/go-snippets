use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    next if ( $line =~ m/^\s$/ );
    my ( $n, $m ) = ( $line =~ /([^,]+),([^\n])/ );
    printf( "%d\n", rindex( $n, $m ) );
}
close(INFILE);
