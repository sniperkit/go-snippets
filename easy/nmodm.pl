use strict;
use integer;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    next if ( $line =~ m/^\s$/ );
    my ( $n, $m ) = ( $line =~ /(\d+),(\d+)/ );
    my $r = $n - ( $n / $m ) * $m;
    print "$r\n";
}
close(INFILE);
