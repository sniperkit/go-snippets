use strict;
use List::Util qw[min];
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    my @l = map { length($_) } $line =~ /X\.*Y/g;
    printf( "%d\n", min(@l) - 2 );
}
