use strict;
my %replace = (
    'zero'  => 0,
    'one'   => 1,
    'two'   => 2,
    'three' => 3,
    'four'  => 4,
    'five'  => 5,
    'six'   => 6,
    'seven' => 7,
    'eight' => 8,
    'nine'  => 9,
    ';'     => ''
);
my $words = qr/@{[join '|', map { quotemeta($_) } keys %replace]}/;
open( INFILE, "<", $ARGV[0] )
  or die("Cannot open file $ARGV[0] for reading: $!");
while ( my $line = <INFILE> ) {
    next if ( $line =~ m/^\s$/ );
    chomp($line);
    $line =~ s/($words)/$replace{$1}/g;
    print "$line\n";
}
close(INFILE);
