use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    my $l = scalar grep { $_ =~ /^[[:lower:]]/ } split( //, $line );
    my $u = scalar grep { $_ =~ /^[[:upper:]]/ } split( //, $line );
    ( $l, $u ) = ( 100 * $l / ( $l + $u ), 100 * $u / ( $l + $u ) )
      if ( $l + $u > 0 );
    printf( "lowercase: %.2f uppercase: %.2f\n", $l, $u );
}
close(INFILE);
