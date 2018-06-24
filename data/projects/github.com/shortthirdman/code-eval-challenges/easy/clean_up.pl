use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    chomp($line);
    $line = lc $line;
    $line =~ s/^[^a-z]+|[^a-z]+$//g;
    $line =~ s/[^a-z]+/ /g;
    printf "$line\n";
}
close(INFILE);
