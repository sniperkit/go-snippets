use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    chomp $line;
    my $swap = 0;
    my @r;
    foreach ( split( / /, $line ) ) {
        if ($swap) {
            $_ =~ /(\d+)-(\d+)/;
            ( $r[$1], $r[$2] ) = ( $r[$2], $r[$1] );
        }
        elsif ( $_ eq ":" ) {
            $swap = 1;
        }
        else {
            push @r, $_;
        }
    }
    printf( "%s\n", join( ' ', @r ) );
}
close(INFILE);
