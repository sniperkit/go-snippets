use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    chomp $line;
    my $p = 1;
    for ( my $i = 1 ; $i < length $line ; $i++ ) {
        $p = $i + 1
          if ( substr( $line, $i, 1 ) ne substr( $line, $i - $p, 1 ) );
    }
    print "$p\n";
}
close(INFILE);
