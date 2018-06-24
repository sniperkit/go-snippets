use strict;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    next if ( $line =~ m/^\s$/ );
    my @s = split( / /, $line );
    my $l = $s[-1];
    if ( ( $l > 0 ) && ( $l < scalar @s ) ) {
        my $r = $s[ -$l - 1 ];
        print "$r\n";
    }
}
close(INFILE);
