use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    next if ( $line =~ m/^\s$/ );
    my $r = 0;
    map { $r += $_ } $line =~ /"id": (\d+),/g;
    print "$r\n";
}
close(INFILE);
