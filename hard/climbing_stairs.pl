use strict;
use bignum;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    next if ( $line =~ m/^\s$/ );
    chomp($line);
    if ( $line < 4 ) {
        print "$line\n";
        next;
    }
    my ( $b, $c ) = ( 1, 1 );
    for ( ; $line > 1 ; $line-- ) {
        ( $b, $c ) = ( $c, $b + $c );
    }
    print "$c\n";
}
close(INFILE);
