use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    chomp($line);
    my ( $l, $d, $n, @s ) = split( / /, $line );
    my ( $count, $t, $tx, $i ) = ( 0, 0, 6 - $d, 6 );
    while ( $i <= $l - 6 ) {
        if ( $i > $tx - $d ) {
            $i = $tx;
            if ( $t == $n ) {
                $tx = $l - 6 + $d;
            }
            else {
                $tx = $s[$t];
                $t++;
            }
        }
        else { $count++; }
        $i += $d;
    }
    print "$count\n";
}
close(INFILE);
