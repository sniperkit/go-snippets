use strict;

sub balanced {
    my ( $s, $c ) = @_;
    while ( $c >= 0 ) {
        return ( $c == 0 ) if ( !( $s =~ /[:()]/ ) );
        return ( balanced( substr( $s, 2 ), $c )
              || balanced( substr( $s, 2 ), $c + 1 ) )
          if ( substr( $s, 0, 2 ) eq ':(' );
        return ( balanced( substr( $s, 2 ), $c )
              || balanced( substr( $s, 2 ), $c - 1 ) )
          if ( substr( $s, 0, 2 ) eq ':)' );
        $c++ if ( substr( $s, 0, 1 ) eq '(' );
        $c-- if ( substr( $s, 0, 1 ) eq ')' );
        $s =~ s/^.[^:()]*//;
    }
    return 0;
}
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    printf( "%s\n", balanced( $line, 0 ) ? "YES" : "NO" );
}
close(INFILE);
