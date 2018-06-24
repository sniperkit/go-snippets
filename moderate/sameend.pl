use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    next if ( $line =~ m/^\s$/ );
    chomp($line);
    my @s = split( /,/, $line );
    print( ( substr( @s[0], -length @s[1], length @s[1] ) eq @s[1] )
        ? "1\n"
        : "0\n" );
}
close(INFILE);
