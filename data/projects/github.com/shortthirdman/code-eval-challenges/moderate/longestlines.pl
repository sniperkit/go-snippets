use strict;
my $n = 0;
my @l;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    chomp($line);
    if ( $n == 0 ) {
        $n = $line;
        next;
    }
    for ( my $i = 0 ; $i < $n ; $i++ ) {
        if ( $i == scalar @l ) {
            push @l, $line;
            last;
        }
        ( @l[$i], $line ) = ( $line, @l[$i] )
          if ( length $line > length @l[$i] );
    }
}
printf( "%s\n", join( "\n", @l ) );
close(INFILE);
