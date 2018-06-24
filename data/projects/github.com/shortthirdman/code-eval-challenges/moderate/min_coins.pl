use strict;
use integer;
my @b = ( 0, 1, 2, 1, 2 );
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    next if ( $line =~ m/^\s$/ );
    printf( "%s\n", $line / 5 + @b[ $line % 5 ] );
}
close(INFILE);
