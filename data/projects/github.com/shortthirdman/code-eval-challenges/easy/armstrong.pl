use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    next if ( $line =~ m/^\s$/ );
    chomp($line);
    my $r;
    map { $r += $_**length $line } split( //, $line );
    print( $r == $line ? "True\n" : "False\n" );
}
close(INFILE);
