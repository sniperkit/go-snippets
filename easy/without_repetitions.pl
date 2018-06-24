use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    chomp($line);
    $line =~ s/(.)\1+/\1/g;
    printf "$line\n";
}
close(INFILE);
