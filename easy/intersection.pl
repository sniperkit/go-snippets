use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    chomp($line);
    my ( $s, $t ) = split( /;/, $line );
    my @v = split( /,/, $t );
    my @r;
    foreach my $i ( split( /,/, $s ) ) {
        shift @v while ( @v && @v[0] < $i );
        if ( @v && @v[0] == $i ) {
            push @r, shift @v;
        }
    }
    printf( "%s\n", join( ',', @r ) );
}
close(INFILE);
