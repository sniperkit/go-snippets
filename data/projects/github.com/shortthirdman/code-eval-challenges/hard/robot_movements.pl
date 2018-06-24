use strict;

sub robot {
    my ( $f, $x, $y ) = @_;
    return (1) if ( $x == 3 && $y == 3 );
    my $r = 0;
    $r += robot( $f | ( 1 << ( $x - 1 + 4 * $y ) ), $x - 1, $y )
      if ( $x > 0 && ( $f & ( 1 << ( $x - 1 + 4 * $y ) ) ) == 0 );
    $r += robot( $f | ( 1 << ( $x + 4 * ( $y - 1 ) ) ), $x, $y - 1 )
      if ( $y > 0 && ( $f & ( 1 << ( $x + 4 * ( $y - 1 ) ) ) ) == 0 );
    $r += robot( $f | ( 1 << ( $x + 1 + 4 * $y ) ), $x + 1, $y )
      if ( $x < 3 && ( $f & ( 1 << ( $x + 1 + 4 * $y ) ) ) == 0 );
    $r += robot( $f | ( 1 << ( $x + 4 * ( $y + 1 ) ) ), $x, $y + 1 )
      if ( $y < 3 && ( $f & ( 1 << ( $x + 4 * ( $y + 1 ) ) ) ) == 0 );
    return ($r);
}

printf( "%d\n", robot( 1, 0, 0 ) );
