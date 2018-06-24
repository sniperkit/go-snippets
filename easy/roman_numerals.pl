use strict;

my @ronum = ( 1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1 );
my @rostr =
  ( "M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I" );
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    chomp $line;
    my $i;
    while ( $line > 0 ) {
        if ( $line >= @ronum[$i] ) {
            print "@rostr[$i]";
            $line -= @ronum[$i];
        }
        else {
            $i++;
        }
    }
    print "\n";
}
close(INFILE);
