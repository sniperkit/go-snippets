use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    chomp($line);
    my ( $b, $c ) = ( -1, 1 );
    foreach my $i ( split( / /, $line ) ) {
        if ( $i == $b ) {
            $c++;
        }
        else {
            if ( $b >= 0 ) {
                print "$c $b ";
                $c = 1;
            }
            $b = $i;
        }
    }
    print "$c $b\n";
}
