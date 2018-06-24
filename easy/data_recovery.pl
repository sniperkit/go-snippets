use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    chomp $line;
    my ( $a, $b ) = split( /;/, $line );
    my @l = split( / /, $a );
    my @r = split( / /, $b );
    for ( my $i = 1 ; $i <= scalar @l ; $i++ ) {
        push( @r, $i ) if ( !grep { $_ == $i } @r );
    }
    my @s;
    for ( my $i = 0 ; $i < scalar @l ; $i++ ) {
        $s[ @r[$i] - 1 ] = @l[$i];
    }
    printf( "%s\n", join( ' ', @s ) );
}
close(INFILE);
