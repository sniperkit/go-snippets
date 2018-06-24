use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    my ( $l, $r ) = ( ord $line, substr( $line, 1, 1 ) );
    my @s;
    push @s, chr( $l - 2 ) . ( $r - 1 ) if $l > ord "b" && $r > 1;
    push @s, chr( $l - 2 ) . ( $r + 1 ) if $l > ord "b" && $r < 8;
    push @s, chr( $l - 1 ) . ( $r - 2 ) if $l > ord "a" && $r > 2;
    push @s, chr( $l - 1 ) . ( $r + 2 ) if $l > ord "a" && $r < 7;
    push @s, chr( $l + 1 ) . ( $r - 2 ) if $l < ord "h" && $r > 2;
    push @s, chr( $l + 1 ) . ( $r + 2 ) if $l < ord "h" && $r < 7;
    push @s, chr( $l + 2 ) . ( $r - 1 ) if $l < ord "g" && $r > 1;
    push @s, chr( $l + 2 ) . ( $r + 1 ) if $l < ord "g" && $r < 8;
    printf( "%s\n", join( ' ', @s ) );
}
close(INFILE);
