use strict;
my @s = (
    "",         "One",      "Two",      "Three",   "Four",    "Five",
    "Six",      "Seven",    "Eight",    "Nine",    "Ten",     "Eleven",
    "Twelve",   "Thirteen", "Fourteen", "Fifteen", "Sixteen", "Seventeen",
    "Eighteen", "Nineteen", "Twenty",   "Thirty",  "Forty",   "Fifty",
    "Sixty",    "Seventy",  "Eighty",   "Ninety",  "Hundred", "Thousand",
    "Million"
);

sub wrd {
    my ( $a, $b, $c ) = @_;
    return (0) if ( $a + $b + $c == 0 );
    print $s[$a] . $s[28] if ( $a > 0 );
    printf( "%s", $b > 1 ? $s[ 18 + $b ] . $s[$c] : $s[ 10 * $b + $c ] );
    return (1);
}
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    chomp $line;
    if ( $line == 0 ) { print "Zero"; }
    else {
        my @l = split( //, $line );
        unshift @l, 0 while ( scalar @l < 9 );
        print $s[30] if wrd( $l[0], $l[1], $l[2] );
        print $s[29] if wrd( $l[3], $l[4], $l[5] );
        wrd( $l[6], $l[7], $l[8] );
    }
    print "Dollars\n";
}
close(INFILE);
