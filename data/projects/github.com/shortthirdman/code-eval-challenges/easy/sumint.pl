use strict;
my $sum = 0;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    next if ( $line =~ m/^\s$/ );
    $sum += $line;
}
print "$sum\n";
close(INFILE);
