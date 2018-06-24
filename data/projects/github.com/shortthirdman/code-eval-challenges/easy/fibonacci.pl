use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    next if ( $line =~ m/^\s$/ );
    chomp($line);
    if ( $line < 2 ) {
        print "$line\n";
        next;
    }
    my ( $b, $c ) = ( 0, 1 );
    for ( ; $line > 1 ; $line-- ) {
        ( $b, $c ) = ( $c, $b + $c );
    }
    print "$c\n";
}
close(INFILE);
