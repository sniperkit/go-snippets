use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    chomp $line;
    my ( $l, $r ) = split( /;/, $line );
    if ( $r eq '' ) {
        print "\n";
        next;
    }
    my @s = split( /,/, $r );
    for ( my $i = 0 ; $i < scalar @s ; $i += 2 ) {
        if ( $s[ $i + 1 ] eq '' ) {
            $s[ $i + 1 ] = "x";
        }
        else {
            $s[ $i + 1 ] =~ tr/01/ab/;
        }
        $l =~ s/$s[$i]/$s[$i + 1]/g;
    }
    $l =~ tr/abx/01/d;
    printf "$l\n";
}
close(INFILE);
