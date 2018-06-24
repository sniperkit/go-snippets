use strict;
use bigint;
for ( my $i = 9 ; $i >= 1 ; $i -= 2 ) {
    for ( my $j = 9 ; $j >= 0 ; $j-- ) {
        my $n = 101 * $i + 10 * $j;
        if ( 2**( $n - 1 ) % $n == 1 ) {
            print "$n\n";
            exit 0;
        }
    }
}
