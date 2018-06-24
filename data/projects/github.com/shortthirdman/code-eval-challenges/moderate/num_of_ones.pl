use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    next if ( $line =~ m/^\s$/ );
    my $r = 0;
    while ( $line > 0 ) {
        $line &= $line - 1;
        $r++;
    }
    print "$r\n";
}
close(INFILE);
