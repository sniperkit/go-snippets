use strict;
use POSIX;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $a = <INFILE> ) {
    next if ( $a =~ m/^\s$/ );
    my $b = floor($a);
    print "$b.";
    $a = ( $a - $b ) * 60;
    $b = floor($a);
    printf( "%02d'", $b );
    $a = ( $a - $b ) * 60;
    $b = floor($a);
    printf( "%02d\"\n", $b );
}
